# Airflow
* Create your DAG in `dags` folder
* Goto `localhost:8080` and login with `admin` and `admin`
* See your DAG runs

## How to scale airflow?
* CeleryExecutor runs tasks in multiple worker nodes
* KubernetesExecutor runs each task in it's own k8s pod. 
  * I've heard it costs too much
