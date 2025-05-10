# Multi Party Encryption
* It's a form of encryption where the secret key is sharded and shared between multiple parties
* Compromise of one node will not reveal the secret key
* There must be a quorum to reconstruct the master secret key
* Shamir's Secret Sharing -- An example. HashiCorp Vault uses this

NOTE: I could not make the python script working, but got the gist. Will revisit when I have time.

## Resources
* https://rya-sge.github.io/access-denied/2024/10/21/mpc-protocol-overview/
* YaoMillionaireProtocol
