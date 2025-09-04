sudo sh -c 'echo 3 > /proc/sys/vm/drop_caches'
sudo perf trace -e block:block_rq_issue ./main a.txt > /dev/null
