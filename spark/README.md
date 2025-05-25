# Spark Cluster
* This folder demonstrates a very basic setup with `apache:spark` cluster
* The bitnami-spark throws Kerberos error while reading files from local, and LLMs couldn't resolve the errors for me
* Use this setup to run your experiments.

# Spark
* Ivy is the dependency manager for Spark (LOL how many for Java?)
  * Hence you have to set the ivy2 cache directory in env var
* Lazy Evaluation, Transformations vs Actions, DataFrame and Dataset API, RDD
* Once the cluster is up, run the `submit_job.py` to submit a job
* The jobs are placed inside the spark-apps folder. Each python file is like a job that you can submit.
* This experiment demonstrates a spark cluster running locally in standalone mode.
  * This has a master node and 2 worker nodes
  * The worker node registers itself with the master
  * The jobs are submitted to master, which then splits the work to workers
* Job vs Stage vs Task
  * Job is an user action -  `collect()`, `count()`, `save()`
    * Spark is lazy - Transformations do nothing until an action is called.
    * `df.count()` triggers one job
  * Stage = A Set of Transformations that can run without shuffles
    * Spark breaks the job into stages based on shuffle boundaries.
    > 💡 Example:
    > df.map(...).groupBy(...).agg(...)
    > → First stage: map
    > → Second stage: groupBy + agg (requires a shuffle)
  * Task = smallest unit of work; runs on a single executor thread
    * Each stage is divided into tasks, one per partition.
    * All tasks in the same stage run in parallel across executors.
* Warehouse directory is where spark manages data by default

## Hadoop + Spark 
* See hadoop readme for an example

# FAQ
## What Happens When a Worker Fails?
* Task Fails:
  * Tasks running on that worker crash or stop responding.
  * Spark notices the failure (via heartbeat timeout or task error).
* Task is Rescheduled:
  * The Spark Driver detects the failed task.
  * It reschedules the task on a different active executor (another worker).
  * The data required for the task is fetched again if needed.
* No Data Loss (usually):
  * If the data is in persistent storage (e.g., HDFS, S3, local disk), Spark can re-read it.
  * If you're using cached data that was only in memory on the failed worker, Spark may recompute it using lineage info.

## Other aspects
* RDD Lineage. Spark knows how to construct RDDs from other RDDs (parents). It just reapplies the transformation
* Retries Task execution in another executor if node goes down
* Persistence to storage
* Checkpointing - Store intermediate RDDs reliably in HDFS to avoid full recompute
