<p align="center">
  <img src="https://raw.githubusercontent.com/cpl/cryptor/master/docs/img/logo.png" width="500px" alt="Cryptor Logo"/>
</p>


<p align="center">
  <a href="https://cpl.li/cryptor">
    <img src="https://img.shields.io/badge/docs-cpl.li-informational.svg" alt="Official Docs" />
  </a>
  <a href="https://travis-ci.org/cpl/cryptor">
    <img src="https://img.shields.io/travis/cpl/cryptor/master.svg" alt="Travis CI" />
  </a>
  <a href="https://goreportcard.com/report/cpl.li/go/cryptor">
    <img src="https://goreportcard.com/badge/cpl.li/go/cryptor" alt="Go Report Card" />
  </a>
  <a href="https://coveralls.io/github/cpl/cryptor?branch=master">
    <img src="https://img.shields.io/coveralls/github/cpl/cryptor/master.svg" alt="Coverage Status" />
  </a>
</p>

---


## Introduction
Cryptor is an [overlay](https://en.wikipedia.org/wiki/Overlay_network#Over_the_Internet) [P2P](https://en.wikipedia.org/wiki/Peer-to-peer) network that values your privacy and anonymity above all else. When faced with a tradeoff between the user’s privacy and convenience, privacy comes on top.

To achieve it’s level of security and privacy, a few key concepts are enforced:
* **Obscurity is not security**. Do not assume that you’re secure because somebody does not know the *rules you play by*.
* **Silence is a virtue**. Do not respond to any invalid or unexpected requests, this allows you to hide from mass scanning attempts and some other *attacks*.
* **Less is more**. The less information you have to send, the better. The protocol only needs to send as little as possible.
* **Deceive**. Sending only requested data to your peers puts you at risk. Cryptor nodes will send random and decoy packets across the network to make the real ones harder to detect.
* **0 Pattern**. Some ISPs and companies use advanced heuristics and ML tools to detect certain unwanted traffic and block it. By having random packet sizes and encrypted payloads makes Cryptor harder to pinpoint.


For more information check out the [official documentation](https://cpl.li/cryptor/) or the [Discord server](https://discord.gg/vGQ76Uz).

## Install
Cryptor as a package only provides the *backend* stack for accessing the network and managing local resources. The Cryptor clients are developed under [`cryptor/cmd`](https://github.com/cpl/cryptor/tree/master/cmd). Each client has its own binary and installation routine.

Visit the [official documentation](https://cpl.li/cryptor/) for a list of all clients are their purpose/usage.

### Download Binaries
You can download pre-compiled binaries from the [GitHub Releases](https://github.com/cpl/cryptor/releases) page. Always check checksum of what you download, and only download official binaries from the GitHub repo.

### Go Get
If you have [Go installed](https://golang.org/doc/install), you can use `go get`. Running `go get cpl.li/go/cryptor` will download the entire Cryptor package and all clients. To download and install a specific client use `go get cpl.li/go/cryptor/cmd/<CLIENT NAME>`.


### Manually
Git clone [the project](https://github.com/cpl/cryptor) and use the `Makefile` provided for each individual client or package. For detailed instruction read the [official documentation](https://cpl.li/cryptor/).

## Usage
The main Cryptor client is `aegis`.
// TODO

## Documentation
* [Official documentation](https://cpl.li/cryptor)
* [Contribution guide](https://cpl.li/cryptor/docs/contribution/)

## References
* [RFC 5869 - HKDF](https://tools.ietf.org/html/rfc5869)
* [RFC 7693 - BLAKE2 MAC](https://tools.ietf.org/html/rfc7693)
* [RFC 2898 - PKCS](https://tools.ietf.org/html/rfc2898)
* [WireGuard Whitepaper](https://www.wireguard.com/papers/wireguard.pdf)
* [Noise Protocol Framework](https://noiseprotocol.org/noise.pdf)
* [Tor](https://www.torproject.org)
* [Freenet](https://freenetproject.org)

<p align=center>
  <img src="https://imgs.xkcd.com/comics/privacy_opinions.png">
</p>
