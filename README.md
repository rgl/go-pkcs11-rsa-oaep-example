[![Build status](https://github.com/rgl/go-pkcs11-rsa-oaep-example/workflows/Build/badge.svg)](https://github.com/rgl/go-pkcs11-rsa-oaep-example/actions?query=workflow%3ABuild)

# About

This is an example application that uses the [rgl/go-pkcs11-rsa-oaep](https://github.com/rgl/go-pkcs11-rsa-oaep) library.

# Usage

Execute the following instructions in a Ubuntu 20.04 terminal.

Install dependencies:

```bash
sudo apt-get install -y opensc
```

Build and execute:

```bash
export TEST_PKCS11_LIBRARY_PATH='/usr/lib/x86_64-linux-gnu/pkcs11/opensc-pkcs11.so'
export TEST_PKCS11_SO_PIN='3537363231383830'
export TEST_PKCS11_USER_PIN='648219'
export TEST_PKCS11_TOKEN_LABEL='test-token (UserPIN)'
export TEST_PKCS11_KEY_LABEL='test-rsa-2048'
go build -v
./go-pkcs11-rsa-oaep-example
```

You should see something alike:

```
2020/08/29 19:28:05 Loading /usr/lib/x86_64-linux-gnu/pkcs11/opensc-pkcs11.so...
2020/08/29 19:28:06 Opening session to test-token (UserPIN)...
2020/08/29 19:28:06 Getting the test-rsa-2048 key...
2020/08/29 19:28:06 Encrypting abracadabra...
2020/08/29 19:28:06 Decrypting 6283ca815967d1deedafe616e95c0f112098f23778e76e5a8c899164a8ed2f127bbc443abfe03399b22ad3bc765bdb13dd8343332e158c0a2f7f94692fff9f7931cbb5d1715e5546a436700d5cadbca6c68c7839ae56681ff803a2d9ee86cb0db7dbdd421f75bbb05cca93c868e7cbc55d5727497ef05541bf40cee6196687cd62ac8196ec2adc97e808d9f1ed105e6ceaf3ee019bdb2885ea6f3a826b5be1dd2414888f44f8f0900450b8ce87856073158ececa2e95d2bfd6e9dccd2712cb03d0849e7a3cf29bb7db8e9ad533c86a6abdd608199ba61077d5b818b7b097cde1975cf0123fbd5ad2a35a970134453aed66d734b4c44f0b1d1eee4c50dbc93ffa...
2020/08/29 19:28:07 Descrypted as abracadabra
```
