# ndau key format

The keys for ndau are designed to be:
* able to grow as keys and cryptography evolve
* distinct from keys used for other cryptocurrencies
* private and public keys are easily identifiable
* human-readable
* encoded in a way that they could be read aloud and manually typed without confusion
* able to detect minor typos

Essentially, keys are converted to a binary form, and then converted to a base32 representation where case doesn't matter, the characters `i`, `o`, `1` and `0` are omitted, and includes a checksum.


## Algorithm

### Build a binary key

The instructions below are for non-extended keys ONLY. The "extra" bytes that apply to HD key trees are ignored in these instructions; the system supports doing this for extended keys but the explanation of how they are serialized is more complex.

* Determine the key type (0=null key, 1=ed25519, 2=secp256k1)
* Determine the key length:
    * ed25519 public keys are 32 bytes
    * ed25519 private keys are 64 bytes
    * secp256k1 public keys are 33 bytes
    * secp256k1 private keys are 32 bytes
* Build the 5-byte prefix, which consists of `0x92`, the key type, `0xc4`, keylen+1, keylen:
    * for ed25519 public keys it is `9201c42120`
    * for ed25519 private keys it is `9201c44140`
    * for secp256k1 public keys it is `9202c42221`
    * for secp256k1 private keys it is `9202c42120`
* Build the "packed" byte array, which concatenates:
    * The prefix
    * The bytes of the key
* Calculate the checksum of this byte array. This is done by:
    * Calculating the number of bytes in the checksum, which is the number of bytes needed to pad the length
    of the array to a multiple of 5 -- but that number must be a minimum of 3. (So the checksum length will be between 3 and 7 bytes). For the keys below, that will be:
        * ed25519 public keys need 3 bytes
        * ed25519 private keys need 6 bytes
        * secp256k1 public keys need 7 bytes
        * secp256k1 private keys need 3 bytes
    * calculate the checksum as the trailing n bytes of the sha224 checksum of the input bytes
```go
func cksumN(input []byte, n byte) []byte {
	sum := sha256.Sum224(input)
	return sum[sha256.Size224-int(n):]
}
```
    * append the checksum to the end of the input bytes, making its length a multiple of 5
    * convert the input byte stream to base32 using the alphabet `abcdefghijkmnpqrstuvwxyz23456789`
    * add "npub" for public keys and "npvt" for private keys to the front of the string

### Examples

*Ed25519Public*
 raw key: `9e3e08c194fae465c70e0ac0487f63b00dd969c13e45417731d1ab12768f8636` (len 32)
  packed: `9201c421209e3e08c194fae465c70e0ac0487f63b00dd969c13e45417731d1ab12768f8636` (len 37)
  result: `npuba8jadtbbecrd6cgbuv7qi3qhb2fnaud9nq2a5ymj2e9eksmzghi4yevyt8ddnuxhgw33b8ed`

*Ed25519Private*
 raw key: `db531a15b6444c71b790aae7ae9dd6b3f87561319399af6e490ae1078cda9e03cd5c569c7160d694a9a2607944a267b7cd5fed312d3ef854cc46b5cf92857f2e` (len 64)
  packed: `9201c44140db531a15b6444c71b790aae7ae9dd6b3f87561319399af6e490ae1078cda9e03cd5c569c7160d694a9a2607944a267b7cd5fed312d3ef854cc46b5cf92857f2e` (len 69)
  result: `npvtayjadtcbidpxggsxy3ce26pzucxqrmw74439s7mbggj3vm5qjefqcb6n5krahvk6k4qhc2gyuuw4e2d3iutgrp8pm9yvcmj89bkn2txx38jik93qvwy9fxad`

*Secp256k1Public*
 raw key: `02d235e59fc697cb3a2a416c60f1c4af1edf68fdd21001c97091a5f45db4267f1d` (len 33)
  packed: `9202c4222102d235e59fc697cb3a2a416c60f1c4af1edf68fdd21001c97091a5f45db4267f1d` (len 38)
  result: `npuba4jaftbceebpeprfv9djru34fjay22ht2uzt7z5i9zjbaaqjqci4m7c7ysvh8hpwj3zwvgpn`

*Secp256k1Private*
 raw key: `72422691b240a4fd9b7ee6f31be97042ed38d2b8e126fe2d865958a9e2533e4e` (len 32)
  packed: `9202c4212072422691b240a4fd9b7ee6f31be97042ed38d2b8e126fe2d865958a9e2533e4e` (len 37)
  result: `npvta8jaftbbeb3eejwtyjakj9n5r5vrgg9jqbbq4qguzdsup9tps3nxtkrckn9e643gumuxyt7z`
# ndau signature format
