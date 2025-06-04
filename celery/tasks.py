# tasks.py
from celery import Celery

# Connect to Redis broker
app = Celery(
    "tasks",
    broker="redis://localhost:6379/0",
    backend="redis://localhost:6379/0",  # Add this line
)

app.conf.beat_schedule = {
    "say-hello-every-10-seconds": {
        "task": "tasks.say_hello",
        "schedule": 10.0,  # seconds
        # Optional args:
        # "args": (),
    },
}


@app.task
def add(x, y):
    return x + y


@app.task
def say_hello():
    print("Hello! This runs periodically.")
