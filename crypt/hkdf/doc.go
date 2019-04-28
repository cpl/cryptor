/*
Package hkdf - HMAC Key Derivation Function

Inside this package are provided the functions, methods and structures needed
for implementing "HMAC-based Extract-and-Expand Key Derivation Function (HKDF)"
as defined in RFC 5869.

The hasing algorithm used is BLAKE2s as defined by RFC 7693.

The hmac.New is used as opposed to blake2s.New256(key) for the simple reason
that hmac.New seems to implement HMAC as defined in U.S. Federal Information
Processing Standards Publication 198.
*/
package hkdf // import "cpl.li/go/cryptor/crypt/hkdf"
