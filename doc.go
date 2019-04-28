/*
Package cryptor is an [overlay](https://en.wikipedia.org/wiki/Overlay_network#Over_the_Internet) [P2P](https://en.wikipedia.org/wiki/Peer-to-peer) network that values your privacy and anonymity above all else. When faced with a tradeoff between the user’s privacy and convenience, privacy comes on top.

To achieve it’s level of security and privacy, a few key concepts are enforced:
* **Obscurity is not security**. Do not assume that you’re secure because somebody doesn’t know the *rules you play by*.
* **Silence is a virtue**. Do not respond to any invalid or unexpected requests, this allows you to hide from mass scanning attempts and some other *attacks*.
* **Less is more**. The less information you have to send, the better. The protocol only needs to send as little as possible.
* **Deceive**. Sending only requested data to your peers puts you at risk. Cryptor nodes will send random and decoy packets across the network to make the real ones harder to detect.
* **0 Pattern**. Some ISPs and companies use advanced heuristics and ML tools to detect certain unwanted traffic and block it. By having random packet sizes and encrypted payloads makes Cryptor harder to pinpoint.


For more information check out the [official documentation](https://cpl.li/cryptor/) or the [Discord server](https://discord.gg/vGQ76Uz).
*/
package cryptor // import "cpl.li/go/cryptor"
