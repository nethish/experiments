docker exec -it spark-master /opt/spark/bin/spark-submit --master "spark://spark-master:7077" --conf "spark.jars.ivy=/opt/bitnami/spark/.ivy2" /app/app.py
