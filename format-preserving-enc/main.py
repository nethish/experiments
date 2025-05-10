import pyffx

# Define secret key and a tweak (optional)
key = b"mysecretpassword"
tweak = b"tweak123"

# Create an encryptor for numbers (radix=10)
fpe = pyffx.Integer(key, length=16)

# The original credit card number
plaintext = 4111111111111111

# Encrypt
ciphertext = fpe.encrypt(plaintext)
print(f"Encrypted: {ciphertext}")

# Decrypt
decrypted = fpe.decrypt(ciphertext)
print(f"Decrypted: {decrypted}")
