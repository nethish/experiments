# startup.sh
#!/bin/bash

echo "Starting DataNode and NodeManager..."

# Start DataNode in background
hdfs --daemon start datanode

# Start NodeManager in background
yarn --daemon start nodemanager

# Check if they started (optional, for debugging)
sleep 5 # Give them a moment to start
jps

echo "Startup script finished. Keeping container alive..."
# Keep the container running in the foreground
tail -F /dev/null
