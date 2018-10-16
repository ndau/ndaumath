# ndaumath
Definitive math libraries for calculations related to ndau


## Purpose

This is intended to be the canonical location for all of the mathematical calculations for ndau.

The reason for this is that we are not just adding and subtracting -- we have EAI and SIB calculations that involve fractions, and we have to be very careful about things like overflow and precision.

We also use this library to define a robust test suite that is intended to be independent of implementation language. Thus, the same test suite can be used to prove that the math is correct no matter how it is implemented.

Finally, this repository should contain the ultimate reference documentation for ndau's mathematics.

## Packages

The packages in this library contain the low-level math functions used by ndau.

### Address

 Addresses in ndau are similar but not identical to bitcoin addresses. In
 general, we were trying to create an address that was reasonably short
 (typeable by a human), identifiable as an ndau address (and not confused with
 addresses from other kinds of cryptocurrencies), had a built-in checksum
 capability, was case-insensitive, and allowed for a family of related address
 types. Its implementation details are hidden from normal uses.


 ### B32

 B32 is an implementation of a base32 encoding; this is preferable to base64 in
 that it is case-insensitive and avoids 1, 0, i, and l as problematic characters
 in many typefaces. Addresses use this library, as do keys.

 ### Bitset256

 Bitset256 is an implementation of a 256-bit bitset, indexed by a single byte.
 It's intended to be quite fast (certainly faster than a variable-width
 implementation).

### Constants

A collection of the key constants in the ndau universe.

### EAI

 A careful implementation of the math behind EAI. EAI is complex and
 additionally is dependent on non-integral values. This library implements EAI
 using ratios of 64-bit numbers, using 128-bit math for some intermediate
 calculations to avoid overflow errors.

### Key

Support for HD key derivation and manipulation derived from source code for a
bitcoin-specific implementation.

### KeyAddr

An implementation of a client-side Key generation and manipulation library in Go, that uses
gomobile to generate Java and Objective-C code for use by Android and IoS applications.

### ndauErr

Defines a couple of error types used by ndaumath libraries.

### Signature

Implementation of a generic concept of signatures so that ndau can someday have new signature types
added if and when the existing signature types become obsolete.

### Signed

An implementation of 64-bit signed math with overflows, errors, and a couple of special operations
not usually found in such a library: MulDiv (with 128-bit intermediate operations), DivMod, and Exp
that calculates e to a ratio rather than a floating point value.

The point of this library is to define calculations in a way that is guaranteed to be reproduceable
and exact, in other languages and on other hardware.

### Types

Defines some basic types for ndau -- the quanity of ndau, the way timestamps are represented, etc.

### Unsigned

The equivalent of the Signed library, only Unsigned.

### Words

An implementation of the BIP-0039 words-to-bits technique, using the same wordlist. It supports
multiple languages in the API but so far only English is supported.


## Commands

None of the code in the cmd subdirectory is required to make ndau work, and in some cases the
code is obsolete and is only being kept around for experimental purposes.

The ones documented below, however, may sometimes be useful, and may eventually evolve into
full-fledged applications.

### addrtool

addrtool can generate and validate ndau addresses.

### sigtool

sigtool can generate ndau keypairs, sign blocks of data, and validate signatures.

