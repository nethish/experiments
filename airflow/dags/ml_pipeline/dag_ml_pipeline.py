from airflow import DAG
from airflow.operators.python import PythonOperator, BranchPythonOperator
from airflow.operators.dummy import DummyOperator
from datetime import datetime
from ml_pipeline.utils.api import fetch_region_data
from ml_pipeline.operators.validate_data_operator import ValidateDataOperator
from ml_pipeline.operators.train_model_operator import TrainModelOperator

regions = ["us", "eu", "asia"]

default_args = {
    "owner": "data_engineering",
    "start_date": datetime(2024, 1, 1),
    "retries": 1,
}


def decide_branch(**kwargs):
    if kwargs["ti"].xcom_pull(key="validation_failed"):
        return "notify_failure"
    return "train_model"


with DAG(
    "ml_pipeline",
    default_args=default_args,
    schedule_interval="@daily",
    catchup=False,
    tags=["ml", "advanced"],
) as dag:
    start = DummyOperator(task_id="start")

    fetch_tasks = []
    validate_tasks = []

    for region in regions:
        fetch_task = PythonOperator(
            task_id=f"fetch_{region}_data",
            python_callable=fetch_region_data,
            op_kwargs={"region": region},
        )
        validate_task = ValidateDataOperator(
            task_id=f"validate_{region}_data",
            region=region,
        )
        fetch_task >> validate_task
        fetch_tasks.append(fetch_task)
        validate_tasks.append(validate_task)

    validation_check = BranchPythonOperator(
        task_id="validation_check",
        python_callable=decide_branch,
        provide_context=True,
    )

    train_model = TrainModelOperator(
        task_id="train_model",
        model_name="sales_forecast",
    )

    notify_success = DummyOperator(task_id="notify_success")
    notify_failure = DummyOperator(task_id="notify_failure")

    end = DummyOperator(task_id="end")

    start >> fetch_tasks
    validate_tasks >> validation_check
    validation_check >> [train_model, notify_failure]
    train_model >> notify_success
    [notify_success, notify_failure] >> end
