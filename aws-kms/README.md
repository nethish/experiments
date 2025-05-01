# AWS KMS
* AWS KMS manages exactly one key and provides encrypt decrypt service
* The Key never leaves the service and is secure.
* One of the use cases is to generate data key
  * KMS returns two values - unencrypted key and encrypted key (encrypted using the master key)
  * You encrypt your data using the unencryped key and throw away
  * You store the data and the encrypted data key
  * When you want to decrypt, you decrypt the data key and decrypt the data
  * Probably use the stdlib crypto and cipher for local encryption and decryption
* There is also **AEAD** - Authenticated Encrytion with Additional Data
  * You provide additional data with the data you want to encrypt.
  * Then the additional data is used to authenticate while decrypting
  * Encrypt(data: "Nethish's Data", additional_data: "User: Nethish")
  * Some ways according to wikipedia
    * Compute HMAC with ad and cipher text
