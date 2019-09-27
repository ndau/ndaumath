Keyaddr WebAssembly
-------------------

This directory provides a simple wrapper interface to the `github.com/oneiro-ndev/ndaumath/keyaddr` package for WebAssembly consumers like node and web browsers.

WebAssembly is a low level language meant to provide near native execution speeds by limiting . WebAssembly is also a compile target for the Go langauge. Therefore, we can use the same code in our client libraries as we do on the blockchain, which helps keep our sensitive cryptographic results secure and correct.

Building
--------

Run `yarn build.sh`

Testing
-------

Run `yarn test`
