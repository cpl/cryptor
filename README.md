![CRYPTOR](https://imgur.com/9435lX0.jpg)

[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/thee-engineer/cryptor)
[![Go Report Card](https://goreportcard.com/badge/github.com/thee-engineer/cryptor?style=flat-square)](https://goreportcard.com/report/github.com/thee-engineer/cryptor)
[![Codecov](https://img.shields.io/codecov/c/github/thee-engineer/cryptor.svg?style=flat-square)]()
[![Travis Build](https://img.shields.io/travis/thee-engineer/cryptor/master.svg?style=flat-square)](https://github.com/thee-engineer/cryptor)
![Ethereum](https://img.shields.io/badge/Ethereum-0xb296ae1bf5f88B7fE7327116bD9c77805Bc1b7Ef-blue.svg?style=flat-square)
[![HitCount](http://hits.dwyl.io/thee-engineer/cryptor.svg)](http://hits.dwyl.io/thee-engineer/cryptor)
###### Privacy, Anonimity, Freedom.
---


## Description
`Cryptor` is a P2P network designed around protecting the user's identity, freedom and privacy.
The purpose of `Cryptor` is to provide a platform on top of which other applications can be
integrated. A file sharing protocol is beeing designed, which will allow users to share
data across the network: safley and without reavealing the origin of the package or the
idenity of the people requesting the package.

The initial concept came from [Dragos A. Radu](https://github.com/dragosthealex) & [Daniel Hodgson](https://github.com/DanielHodgson) who provided a prototype implementation in Python.

## Our promise
`Cryptor` stands for freedom (we believe the internet should be accesible to everyone to do
whatever they wish), anonimity (it's your right and choice to protect your identity) and privacy
(whatever you do should be for "your eyes only"). These three points will alway be respected by
`Cryptor`. We want users to feel safe while using `Cryptor`.

## Install
In order to install the latest version of `cryptor` you will need to have [Go 1.5](https://golang.org/dl/ "download go") or higher installed. To check if you have Go installed and which verions, just run `go version` in your terminal.

```
# In order to install the CLI tool, use the following:
# REMOVED, will be added in a future release
go install github.com/thee-engineer/cryptor/cmd/cryptor-cli
```

```
# In order to install the WEB tool, use the following:
# WIP, under construction
go install github.com/thee-engineer/cryptor/cmd/cryptor-web
```

<!-- Now you can run `cryptor-cli --help` to see usage information  or `cryptor-web`
to open the dashboard in your web browser. -->

```
go get github.com/thee-engineer/cryptor
# This will download the cryptor source files inside your $GOPATH/src/...

cd $GOPATH/src/thee-engineer/cryptor
# Now you can build it yourself, you can use the Makefile or go build, up to you

make build
# This will put a binary of cryptor inside build/
```