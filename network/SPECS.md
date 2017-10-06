# Cryptor Network specifications

#### P2P
  * Discovery
  * Peer sharing

#### Chunk request
  `R0` is a peer requesting a chunk with hash `cH`
  `s` is a secret (n bytes of rand)
  `M` is the MAC adress of the machine
  SHA512(`s` + `M`) --> `sH` (secret hash used to id requests)
  * Protocol
    1. `R` computes `sH`
    2. `R` sends a request with `sH` as header for `cH` to its peers
    3. Each peer of `R` becomes now a `Rn` (requester of order n)
    4. If a peer has `cH`, then it will return it to `Rn-1` until `R0`
      * This uses TCP Chunk Sharing
    5. If the peer does not have the `cH`, then `Rn` will make a request to
    `Rn+1` for `cH` and also append a new `sH` (of `Rn`) to the list
      * If at any time, a peer gets a request containing it's `sH`, then
      that request will be ignored and it will stop propagation
  * Key observations
    * `s` should be re-generated at a set time interval
      * A copy of the old `sH` will be kept, but new requests will use new `s`
      * To strenghten peer privacy
