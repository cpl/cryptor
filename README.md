![CRYPTOR](https://imgur.com/DhcXZ7E.jpg)

[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/thee-engineer/cryptor)
[![Go Report Card](https://goreportcard.com/badge/github.com/thee-engineer/cryptor?style=flat-square)](https://goreportcard.com/report/github.com/thee-engineer/cryptor)
[![Codecov](https://img.shields.io/codecov/c/github/thee-engineer/cryptor.svg?style=flat-square)]()
[![Travis Build](https://img.shields.io/travis/thee-engineer/cryptor/master.svg?style=flat-square)](https://github.com/thee-engineer/cryptor)
![Ethereum](https://img.shields.io/badge/Ethereum-0xb296ae1bf5f88B7fE7327116bD9c77805Bc1b7Ef-blue.svg?style=flat-square)
[![HitCount](http://hits.dwyl.io/thee-engineer/cryptor.svg)](http://hits.dwyl.io/thee-engineer/cryptor)
###### Anonymous, P2P, secure file sharing.
---


## Description
`Cryptor` is a P2P network designed for sharing data without revealing one's true identity or the nature of the shared files. This is accomplished using AES256 encrypted chunks of files and sharing them across the network. No node will know what package it stores + everything is encrypted junk without the `tail chunk` and the encryption key.

P2P communication uses Elliptic Curve key exchange in order to generate a new symetric AES256 block cipher, used in sending requests between peers.

The initial idea came from [Dragos A. Radu](https://github.com/dragosthealex) and provided a prototype implementation in Python.

## Architecture
// TODO: provide diagram

## Install
In order to install the latest version of `cryptor` you will need to have [Go 1.5](https://golang.org/dl/ "download go") or higher installed. To check if you have Go installed and which verions, just run `go version` in your terminal.

```
# This will download the cryptor source files inside your $GOPATH/src/...
go get github.com/thee-engineer/cryptor
```

```
# In order to install the CLI tool, use the following:
go install github.com/thee-engineer/cryptor/cmd/cryptor-cli
```

Now you can run `cryptor-cli --help` to see usage information.
