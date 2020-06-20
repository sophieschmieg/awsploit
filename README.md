# Proof of Concept for AWS S3 crypto vulnerabilities

This repository contains a proof of concept for three AWS S3 crypto SDK vulnerabilities:

 *  Hash of plaintext is included as metadata: This allows for an offline plaintext recovery attack in `O(plaintext_space)`.
 *  Unauthenticated AES-CBC mode: This allows for an online attack plaintext recovert attack requiring `128 * len(plaintext)` oracle queries on average.
 *  Unauthenticated metadata: Combined with the unauthenticated CBC support, this allows for an online attack against AES-GCM encrypted ciphertext by guessing 16 byte chunks of the plaintext.

This proof of concept mocks out the calls to AWS KMS and AWS S3, storing the objects in local RAM.

It first chooses a plaintext by concatenating four 16 byte segments randomly.
It then encrypts this plaintext using the AWS S3Crypto SDK in AES-GCM and checks that it decrypts successfully.

As a first attack, it uses the plaintext md5 hash to brute force search for the plaintext. The hash is stored as the "X-Amz-Meta-X-Amz-Unencrypted-Content-Md5" metadata.
As a second attack, it guesses the plaintext in chunks of 16 bytes, xors them into the ciphertext, downgrades the encryption to AES-CBC and checks via a padding oracle whether the xored ciphertext decrypts into `IV || counter`.
As a third attack, it prepares by encrypting the plaintext with CBC and runs a standard CBC padding oracle.
