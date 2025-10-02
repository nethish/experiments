FROM python:3.12-slim

WORKDIR /locust
COPY locustfile.py .

RUN pip install --no-cache-dir locust

EXPOSE 8089
CMD ["locust", "-f", "locustfile.py", "--host=http://server:8080"]

