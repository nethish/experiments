# AWS Secret Manager
* Used to store secrets - Passwords, API Tokens etc.
* Encrypted using KMS Key.
  * To decrypt the secret, the associated KMS key permissions is also needed
* Key rotation - Automatically rotate keys with help of lambda. 
  * Key rotation helps in providing more security
