# Malloc experiments
## Custom TCMalloc
* Thread Caching malloc - Each thread will have a pool of memory nodes (with fixed size classes like 4KiB, 8KiB etc), and memory is given from the pool when the thread requests it
* When there is a large memory request, it will go to the page heap managed by the tcmalloc
* The freed memory is released to the thread cache, and to the central cache
