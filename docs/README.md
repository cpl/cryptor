![CRYPTOR](https://imgur.com/9435lX0.jpg)

[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/thee-engineer/cryptor)
[![Go Report Card](https://goreportcard.com/badge/github.com/thee-engineer/cryptor?style=flat-square)](https://goreportcard.com/report/github.com/thee-engineer/cryptor)
[![Coveralls status](https://img.shields.io/coveralls/github/thee-engineer/cryptor.svg?style=flat-square)](https://github.com/thee-engineer/cryptor)
[![Travis Build](https://img.shields.io/travis/thee-engineer/cryptor/master.svg?style=flat-square)](https://github.com/thee-engineer/cryptor)
![Ethereum](https://img.shields.io/badge/Ethereum-0xb296ae1bf5f88B7fE7327116bD9c77805Bc1b7Ef-blue.svg?style=flat-square)
[![HitCount](http://hits.dwyl.io/thee-engineer/cryptor.svg)](http://hits.dwyl.io/thee-engineer/cryptor)
###### Anonymous, P2P, secure file sharing.
---

## Description

**Cryptor** is a P2P network designed for sharing data without revealing one's true identity or the nature of the shared information. Data security is accomplished using asymetric encryption (AES256) and public-private key (RSA) for safe peer-to-peer communication and data integrity validation.

All local data stored on your machine is encrypted using your master password. The chunk design of the file sharing protocol also further improves security.

The P2P Protocol uses fake (or cover) request in order to avoid timing, MITM, 

## Install

In order to install the latest version of `cryptor` you will need to have [Go 1.5](https://golang.org/dl/ "download go") or higher installed. To check if you have Go installed and which verions, just run `go version` in your terminal.

The test suite will only run with Go 1.10 or higher.

```shell
# This will download the cryptor source files inside your $GOPATH/src/...
go get github.com/thee-engineer/cryptor

# Now you can build it yourself, you can use the Makefile or go build, up to you
cd $GOPATH/src/thee-engineer/cryptor

# This will put a binary of cryptor inside build/
make build
```

## Volunteer

The more nodes in the network, the faster and more secure it is for everyone using Cryptor. There are no risks involved. All communication between you and other peers is encrypted. Cryptor is also a "zero knowldge" system allowing you the freedom of Plausible deniability.