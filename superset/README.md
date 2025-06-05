# Superset
Lightweight data viz tool that can connect to different data sources.

Supports scheduling email/ slack alerts that can send the dashboard or chart. This needs celery beat and worker setup

On the backend, it uses Python and SQL Alchemy with celery as the asynchronous task scheduler
