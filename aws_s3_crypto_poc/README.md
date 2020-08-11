# Proof of Concept for AWS S3 crypto vulnerabilities

This repository contains a proof of concept for three AWS S3 crypto SDK vulnerabilities:

 *  Hash of plaintext is included as metadata: This allows for an offline plaintext recovery attack in `O(plaintext_space)`. For a detailed description of these vulnerabitlities see [this security advisory](https://github.com/google/security-research/security/advisories/GHSA-76wf-9vgp-pj7w)
 *  Unauthenticated AES-CBC mode: This allows for an online attack plaintext recovert attack requiring `128 * len(plaintext)` oracle queries on average. For a detailed description of these vulnerabitlities see [this security advisory](https://github.com/google/security-research/security/advisories/GHSA-f5pg-7wfw-84q9)
 *  Unauthenticated metadata: Combined with the unauthenticated CBC support, this allows for an online attack against AES-GCM encrypted ciphertext by guessing 16 byte chunks of the plaintext. For a detailed description of these vulnerabitlities see [this security advisory](https://github.com/google/security-research/security/advisories/GHSA-7f33-f4f5-xwgw)

This proof of concept mocks out the calls to AWS KMS and AWS S3, storing the objects in local RAM.


## Installation

```bash
go get github.com/sophieschmieg/exploits/aws_s3_crypto_poc
go run github.com/sophieschmieg/exploits/aws_s3_crypto_poc
```

