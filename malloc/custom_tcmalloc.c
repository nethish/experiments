#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <pthread.h>
#include <sys/mman.h>
#include <unistd.h>

// ============================================
// 1. SIZE CLASSES (like TCMalloc)
// ============================================
#define NUM_SIZE_CLASSES 8
static const size_t SIZE_CLASSES[] = {
    8, 16, 32, 64, 128, 256, 512, 1024
};

// ============================================
// 2. FREE LIST NODE
// ============================================
typedef struct FreeNode {
    struct FreeNode* next;
} FreeNode;

// ============================================
// 3. THREAD CACHE (per-thread, NO LOCKS!)
// ============================================
typedef struct {
    FreeNode* free_lists[NUM_SIZE_CLASSES];  // One list per size class
    size_t num_objects[NUM_SIZE_CLASSES];    // Count of free objects
} ThreadCache;

// Each thread gets its own cache
__thread ThreadCache* thread_cache = NULL;

// ============================================
// 4. CENTRAL CACHE (shared, HAS LOCKS)
// ============================================
typedef struct {
    FreeNode* free_lists[NUM_SIZE_CLASSES];
    pthread_mutex_t locks[NUM_SIZE_CLASSES];  // One lock per size class
    size_t num_objects[NUM_SIZE_CLASSES];
} CentralCache;

static CentralCache central_cache;

// ============================================
// 5. PAGE HEAP (global memory pool)
// ============================================
typedef struct {
    void* heap_start;
    void* heap_current;
    size_t heap_size;
    pthread_mutex_t lock;  // GLOBAL LOCK
} PageHeap;

static PageHeap page_heap;

// ============================================
// 6. HELPER FUNCTIONS
// ============================================

// Find which size class to use
static int get_size_class(size_t size) {
    for (int i = 0; i < NUM_SIZE_CLASSES; i++) {
        if (size <= SIZE_CLASSES[i]) {
            return i;
        }
    }
    return -1;  // Too large
}

// Initialize the page heap (gets memory from OS)
static void init_page_heap() {
    page_heap.heap_size = 10 * 1024 * 1024;  // 10 MB
    page_heap.heap_start = mmap(NULL, page_heap.heap_size,
                                PROT_READ | PROT_WRITE,
                                MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    page_heap.heap_current = page_heap.heap_start;
    pthread_mutex_init(&page_heap.lock, NULL);
    
    printf("[PageHeap] Allocated 10 MB from OS\n");
}

// Allocate memory from page heap (SLOW PATH, has GLOBAL LOCK)
static void* allocate_from_page_heap(size_t size) {
    pthread_mutex_lock(&page_heap.lock);  // ← GLOBAL LOCK
    
    void* ptr = page_heap.heap_current;
    page_heap.heap_current = (char*)page_heap.heap_current + size;
    
    pthread_mutex_unlock(&page_heap.lock);
    
    printf("[PageHeap] Allocated %zu bytes (GLOBAL LOCK)\n", size);
    return ptr;
}

// Initialize central cache
static void init_central_cache() {
    for (int i = 0; i < NUM_SIZE_CLASSES; i++) {
        central_cache.free_lists[i] = NULL;
        central_cache.num_objects[i] = 0;
        pthread_mutex_init(&central_cache.locks[i], NULL);
    }
}

// Refill central cache from page heap
static void refill_central_cache(int size_class) {
    size_t obj_size = SIZE_CLASSES[size_class];
    size_t batch_size = 64;  // Allocate 64 objects at once
    
    // Get memory from page heap
    void* chunk = allocate_from_page_heap(obj_size * batch_size);
    
    // Split into individual objects and add to free list
    for (size_t i = 0; i < batch_size; i++) {
        FreeNode* node = (FreeNode*)((char*)chunk + i * obj_size);
        node->next = central_cache.free_lists[size_class];
        central_cache.free_lists[size_class] = node;
    }
    
    central_cache.num_objects[size_class] += batch_size;
    printf("[CentralCache] Refilled size class %d with %zu objects\n", 
           size_class, batch_size);
}

// Initialize thread cache
static void init_thread_cache() {
    if (thread_cache != NULL) return;
    
    thread_cache = (ThreadCache*)malloc(sizeof(ThreadCache));
    for (int i = 0; i < NUM_SIZE_CLASSES; i++) {
        thread_cache->free_lists[i] = NULL;
        thread_cache->num_objects[i] = 0;
    }
    
    printf("[ThreadCache] Initialized for thread %lu\n", pthread_self());
}

// Refill thread cache from central cache (has PER-SIZE-CLASS LOCK)
static void refill_thread_cache(int size_class) {
    pthread_mutex_lock(&central_cache.locks[size_class]);  // ← PER-SIZE LOCK
    
    // If central cache is empty, refill it first
    if (central_cache.free_lists[size_class] == NULL) {
        pthread_mutex_unlock(&central_cache.locks[size_class]);
        refill_central_cache(size_class);
        pthread_mutex_lock(&central_cache.locks[size_class]);
    }
    
    // Move 10 objects from central to thread cache
    int batch = 10;
    for (int i = 0; i < batch && central_cache.free_lists[size_class]; i++) {
        FreeNode* node = central_cache.free_lists[size_class];
        central_cache.free_lists[size_class] = node->next;
        central_cache.num_objects[size_class]--;
        
        node->next = thread_cache->free_lists[size_class];
        thread_cache->free_lists[size_class] = node;
        thread_cache->num_objects[size_class]++;
    }
    
    pthread_mutex_unlock(&central_cache.locks[size_class]);
    
    printf("[ThreadCache] Refilled size class %d (per-size lock)\n", size_class);
}

// ============================================
// 7. PUBLIC API
// ============================================

void* my_malloc(size_t size) {
    if (size == 0) return NULL;
    
    init_thread_cache();
    
    int size_class = get_size_class(size);
    
    // Handle large allocations (> 1KB) - go directly to page heap
    if (size_class == -1) {
        printf("[my_malloc] Large allocation, going to PageHeap\n");
        return allocate_from_page_heap(size);
    }
    
    // FAST PATH: Check thread cache (NO LOCK!)
    if (thread_cache->free_lists[size_class] != NULL) {
        FreeNode* node = thread_cache->free_lists[size_class];
        thread_cache->free_lists[size_class] = node->next;
        thread_cache->num_objects[size_class]--;
        
        printf("[my_malloc] FAST PATH: %zu bytes from ThreadCache (NO LOCK!)\n", 
               SIZE_CLASSES[size_class]);
        return (void*)node;
    }
    
    // SLOW PATH: Thread cache empty, refill from central cache
    printf("[my_malloc] ThreadCache empty, refilling...\n");
    refill_thread_cache(size_class);
    
    // Now try again
    FreeNode* node = thread_cache->free_lists[size_class];
    thread_cache->free_lists[size_class] = node->next;
    thread_cache->num_objects[size_class]--;
    
    return (void*)node;
}

void my_free(void* ptr, size_t size) {
    if (ptr == NULL) return;
    
    init_thread_cache();
    
    int size_class = get_size_class(size);
    if (size_class == -1) {
        printf("[my_free] Large allocation, ignoring (simplified)\n");
        return;  // Simplified: don't handle large frees
    }
    
    // Add to thread cache (NO LOCK!)
    FreeNode* node = (FreeNode*)ptr;
    node->next = thread_cache->free_lists[size_class];
    thread_cache->free_lists[size_class] = node;
    thread_cache->num_objects[size_class]++;
    
    printf("[my_free] Returned %zu bytes to ThreadCache (NO LOCK!)\n", 
           SIZE_CLASSES[size_class]);
}

// ============================================
// 8. TEST PROGRAM
// ============================================

void* thread_function(void* arg) {
    int thread_id = *(int*)arg;
    
    printf("\n=== Thread %d started ===\n", thread_id);
    
    // Each thread does some allocations
    void* ptrs[5];
    
    for (int i = 0; i < 5; i++) {
        ptrs[i] = my_malloc(64);
        printf("Thread %d: Allocated ptr[%d] = %p\n", thread_id, i, ptrs[i]);
    }
    
    for (int i = 0; i < 5; i++) {
        my_free(ptrs[i], 64);
        printf("Thread %d: Freed ptr[%d]\n", thread_id, i);
    }
    
    // Allocate again (should be FAST from thread cache)
    printf("\nThread %d: Allocating again (should hit cache)...\n", thread_id);
    void* ptr = my_malloc(64);
    printf("Thread %d: Got %p (FAST PATH!)\n", thread_id, ptr);
    
    return NULL;
}

int main() {
    printf("========================================\n");
    printf("Simple TCMalloc-Style Allocator Demo\n");
    printf("========================================\n\n");
    
    // Initialize
    init_page_heap();
    init_central_cache();
    
    // Create multiple threads
    pthread_t threads[3];
    int thread_ids[3] = {1, 2, 3};
    
    for (int i = 0; i < 3; i++) {
        pthread_create(&threads[i], NULL, thread_function, &thread_ids[i]);
    }
    
    for (int i = 0; i < 3; i++) {
        pthread_join(threads[i], NULL);
    }
    
    printf("\n========================================\n");
    printf("Demo completed!\n");
    printf("========================================\n");
    
    return 0;
}
