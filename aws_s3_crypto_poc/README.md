# Proof of Concept for AWS S3 crypto vulnerabilities

This repository contains a proof of concept for three AWS S3 crypto SDK vulnerabilities:

 *  Hash of plaintext is included as metadata: This allows for an offline plaintext recovery attack in `O(plaintext_space)`.
 *  Unauthenticated AES-CBC mode: This allows for an online attack plaintext recovert attack requiring `128 * len(plaintext)` oracle queries on average.
 *  Unauthenticated metadata: Combined with the unauthenticated CBC support, this allows for an online attack against AES-GCM encrypted ciphertext by guessing 16 byte chunks of the plaintext.

This proof of concept mocks out the calls to AWS KMS and AWS S3, storing the objects in local RAM.

For a detailed description of these vulnerabitlities see [this security advisory](https://github.com/google/security-research/security/advisories/GHSA-76wf-9vgp-pj7w)

## Installation

```bash
go get github.com/sophieschmieg/exploits/aws_s3_crypto_poc
```

