# SQS Experiments
* Receive message from queue using long polling
* Delete the message once processed. If you forget to delete, the message will be redelivered after visibility timeout is over
* After max receive attempts, the message will be moved to DLQ
* Then either use AWS Redrive to move messages from DLQ to their source queue or you do it yourself
