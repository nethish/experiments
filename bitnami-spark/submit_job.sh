docker exec spark-master spark-submit \
  --master spark://spark-master:7077 \
  --conf "spark.jars.ivy=/opt/bitnami/spark/.ivy2" \
  --conf "spark.sql.warehouse.dir=/opt/bitnami/spark/warehouse" \
  /opt/bitnami/spark/jobs/sum.py
