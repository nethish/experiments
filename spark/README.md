# Sparkkkk
* Ivy is the dependency manager for Spark (LOL how many for Java?)
  * Hence you have to set the ivy2 cache directory in env var
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


# TODO
* I could not make spark read the `data.csv` across workers
* After that work with parquet format to see how fast it is?
