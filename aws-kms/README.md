# AWS KMS
* AWS KMS manages exactly one key and provides encrypt decrypt service
* The Key never leaves the service and is secure.
* One of the use cases is to generate data key
  * KMS returns two values - unencrypted key and encrypted key (encrypted using the master key)
  * You encrypt your data using the unencryped key and throw away
  * You store the data and the encrypted data key
  * When you want to decrypt, you decrypt the data key and decrypt the data
  * Probably use the stdlib crypto and cipher for local encryption and decryption
