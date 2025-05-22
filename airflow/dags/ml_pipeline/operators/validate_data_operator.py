from airflow.models import BaseOperator
from airflow.utils.decorators import apply_defaults
from airflow.exceptions import AirflowSkipException
import random


class ValidateDataOperator(BaseOperator):
    @apply_defaults
    def __init__(self, region: str, *args, **kwargs):
        self.region = region
        super().__init__(*args, **kwargs)

    def execute(self, context):
        # Simulate validation
        if random.random() < 0.2:  # 20% chance of failure
            context["ti"].xcom_push(key="validation_failed", value=True)
            raise AirflowSkipException(f"Validation failed for region {self.region}")
        self.log.info(f"Validation passed for {self.region}")
