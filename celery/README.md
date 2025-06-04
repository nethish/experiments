# Python Celery
Celery is an asynchronous task queue for Python. It lets you run time-consuming or scheduled tasks in the background, outside your main application process.

The Celery result backend is a storage system where Celery stores the results of tasks so that you can later:

Retrieve the result of a task (result.get())

Track the status of a task (PENDING, STARTED, SUCCESS, FAILURE, etc.)

Handle retries, chaining, and monitoring

## What can it do for you?
* Run Asynchronous tasks
* Run scheduled or periodic jobs
* Run Distributed tasks

## Components
### Broker
Sends and queues tasks in Redis, RabbitMQ etc

### Backend
Stores metadata, and tasks result 

## Commands
```bash
# Run this command in multiple terminals
# Install celery and redis library
celery -A tasks worker --loglevel=info

# Bean scheduler for periodic job
celery -A tasks beat --loglevel=info


python main.py
```
