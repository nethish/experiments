docker exec -it spark-master /opt/spark/bin/spark-submit \
  --master spark://spark-master:7077 \
  --conf "spark.jars.ivy=/tmp/.ivy2" \
  --packages org.apache.spark:spark-sql-kafka-0-10_2.12:3.4.1 \
  /app/main.py
