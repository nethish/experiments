from airflow.models import BaseOperator
from airflow.utils.decorators import apply_defaults

class TrainModelOperator(BaseOperator):

    @apply_defaults
    def __init__(self, model_name: str, *args, **kwargs):
        self.model_name = model_name
        super().__init__(*args, **kwargs)

    def execute(self, context):
        self.log.info(f"Training model: {self.model_name}")
        # Simulate training
        import time
        time.sleep(10)
        self.log.info("Training complete.")
