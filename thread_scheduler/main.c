#include <stdio.h>
#include <stdlib.h>
#include <ucontext.h>
#include <signal.h>
#include <sys/time.h>
#include <unistd.h>
#include <stdbool.h>

#define STACK_SIZE 8192
#define MAX_THREADS 10

// Thread states
typedef enum {
    READY,
    RUNNING,
    SLEEPING,
    DEAD
} ThreadState;

// Thread control block
typedef struct {
    int id;
    ucontext_t context;
    char stack[STACK_SIZE];
    ThreadState state;
    uint64_t sleep_until;  // Microseconds since epoch
} Thread;

// Scheduler state
typedef struct {
    Thread threads[MAX_THREADS];
    int current_thread;
    int thread_count;
    ucontext_t main_context;
} Scheduler;

Scheduler scheduler;

// Get current time in microseconds
uint64_t get_current_time_us() {
    struct timeval tv;
    gettimeofday(&tv, NULL);
    return tv.tv_sec * 1000000ULL + tv.tv_usec;
}

// Initialize scheduler
void scheduler_init() {
    scheduler.current_thread = -1;
    scheduler.thread_count = 0;
}

// Find next runnable thread
int find_next_thread() {
    uint64_t now = get_current_time_us();
    
    // Wake up sleeping threads
    for (int i = 0; i < scheduler.thread_count; i++) {
        if (scheduler.threads[i].state == SLEEPING && 
            now >= scheduler.threads[i].sleep_until) {
            scheduler.threads[i].state = READY;
        }
    }
    
    // Find next ready thread (round-robin)
    for (int i = 0; i < scheduler.thread_count; i++) {
        int idx = (scheduler.current_thread + 1 + i) % scheduler.thread_count;
        if (scheduler.threads[idx].state == READY) {
            return idx;
        }
    }
    
    return -1;
}

// Context switch to next thread
void schedule() {
    int next = find_next_thread();
    
    if (next == -1) {
        // No runnable threads, return to main
        printf("[Scheduler] No more runnable threads\n");
        return;
    }
    
    int prev = scheduler.current_thread;
    scheduler.current_thread = next;
    scheduler.threads[next].state = RUNNING;
    
    if (prev != -1 && scheduler.threads[prev].state == RUNNING) {
        scheduler.threads[prev].state = READY;
    }
    
    printf("[Scheduler] Switching from thread %d to thread %d\n", prev, next);
    
    if (prev == -1) {
        // First switch from main
        swapcontext(&scheduler.main_context, &scheduler.threads[next].context);
    } else {
        swapcontext(&scheduler.threads[prev].context, &scheduler.threads[next].context);
    }
}

// Timer signal handler for preemptive scheduling
void timer_handler(int signum) {
    (void)signum;
    
    if (scheduler.current_thread != -1) {
        printf("[Timer] Preempting thread %d\n", scheduler.current_thread);
        schedule();
    }
}

// Setup preemptive timer (context switches every 100ms)
void setup_timer() {
    struct sigaction sa;
    struct itimerval timer;
    
    sa.sa_handler = timer_handler;
    sigemptyset(&sa.sa_mask);
    sa.sa_flags = 0;
    sigaction(SIGALRM, &sa, NULL);
    
    // Fire every 100ms
    timer.it_value.tv_sec = 0;
    timer.it_value.tv_usec = 100000;
    timer.it_interval.tv_sec = 0;
    timer.it_interval.tv_usec = 100000;
    
    setitimer(ITIMER_REAL, &timer, NULL);
}

// Thread yields CPU voluntarily
void thread_yield() {
    printf("[Thread %d] Yielding\n", scheduler.current_thread);
    schedule();
}

// Thread sleeps for specified microseconds
void thread_sleep(uint64_t microseconds) {
    int current = scheduler.current_thread;
    printf("[Thread %d] Sleeping for %lu us\n", current, microseconds);
    
    scheduler.threads[current].state = SLEEPING;
    scheduler.threads[current].sleep_until = get_current_time_us() + microseconds;
    
    schedule();
}

// Thread exits
void thread_exit() {
    int current = scheduler.current_thread;
    printf("[Thread %d] Exiting\n", current);
    
    scheduler.threads[current].state = DEAD;
    schedule();
}

// Wrapper to call thread function and exit
void thread_wrapper(void (*func)(int), int arg) {
    func(arg);
    thread_exit();
}

// Spawn a new thread
int thread_spawn(void (*func)(int), int arg) {
    if (scheduler.thread_count >= MAX_THREADS) {
        return -1;
    }
    
    int tid = scheduler.thread_count++;
    Thread *t = &scheduler.threads[tid];
    
    t->id = tid;
    t->state = READY;
    
    getcontext(&t->context);
    t->context.uc_stack.ss_sp = t->stack;
    t->context.uc_stack.ss_size = STACK_SIZE;
    t->context.uc_link = &scheduler.main_context;
    
    makecontext(&t->context, (void (*)(void))thread_wrapper, 2, func, arg);
    
    printf("[Scheduler] Spawned thread %d\n", tid);
    return tid;
}

// ============ Example User Threads ============

void worker_thread(int id) {
    for (int i = 0; i < 3; i++) {
        printf("  Worker %d: iteration %d\n", id, i);
        
        if (i == 1) {
            thread_sleep(200000);  // Sleep 200ms
        } else {
            thread_yield();  // Voluntary yield
        }
    }
}

void fast_thread(int id) {
    for (int i = 0; i < 5; i++) {
        printf("  Fast %d: tick %d\n", id, i);
        thread_yield();
    }
}

// ============ Main ============

int main() {
    printf("=== Simple Thread Scheduler Demo ===\n");
    printf("Current user: nethish\n");
    printf("Current time: 2025-12-21 00:30:48\n\n");
    
    scheduler_init();
    setup_timer();  // Enable preemptive scheduling
    
    // Spawn some threads
    thread_spawn(worker_thread, 1);
    thread_spawn(worker_thread, 2);
    thread_spawn(fast_thread, 3);
    
    printf("\n[Main] Starting scheduler...\n\n");
    
    // Start scheduling
    schedule();
    
    // After all threads finish, we return here
    printf("\n[Main] All threads completed!\n");
    
    // Disable timer
    struct itimerval timer = {0};
    setitimer(ITIMER_REAL, &timer, NULL);
    
    return 0;
}