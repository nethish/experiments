# Dockerfile
FROM python:3.10-slim

# Install dependencies
RUN pip install --no-cache-dir hdfs
RUN pip install --no-cache-dir pyarrow

# Copy Python script
WORKDIR /app
COPY main.py .

# Run the script
CMD ["python", "main.py"]

