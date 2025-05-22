from airflow import DAG
from airflow.operators.python import PythonOperator
from datetime import datetime, timedelta


def print_start(**context):
    print(context)
    print(f"DAG started at {datetime.utcnow()}")


def print_end(**context):
    print(context)
    print(f"DAG ended at {datetime.utcnow()}")


default_args = {
    "owner": "airflow",
    "retries": 0,
}

with DAG(
    dag_id="print_start_end_every_2_minutes",
    default_args=default_args,
    description="Prints start and end time every 2 minutes",
    schedule_interval="*/2 * * * *",
    start_date=datetime(2024, 1, 1),
    catchup=False,
    tags=["example"],
) as dag:
    start = PythonOperator(task_id="print_start_time", python_callable=print_start)

    end = PythonOperator(task_id="print_end_time", python_callable=print_end)

    start >> end
