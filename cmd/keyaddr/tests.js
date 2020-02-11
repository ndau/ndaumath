/* ----- ---- --- -- -
 * Copyright 2020 The Axiom Foundation. All Rights Reserved.
 *
 * Licensed under the Apache License 2.0 (the "License").  You may not use
 * this file except in compliance with the License.  You can obtain a copy
 * in the file LICENSE in the source distribution or at
 * https://www.apache.org/licenses/LICENSE-2.0.txt
 * - -- --- ---- -----
 */


import chai, { expect } from 'chai'
import fs from 'fs'
import { promisify } from 'util'
import chaiAsPromised from 'chai-as-promised'
require('./wasm_exec')
const readFile = promisify(fs.readFile)
chai.use(chaiAsPromised)

const toUint8Array = b => {
  var u = new Uint8Array(b.length)
  for (var i = 0; i < b.length; ++i) {
    u[i] = b[i]
  }
  return u
}

const instantiateStreaming = (source, importObject) => {
  importObject = importObject || {}
  return source
    .then(response => Promise.resolve(toUint8Array(response)))
    .then(arrayBuffer => Promise.resolve(new WebAssembly.Module(arrayBuffer)))
    .then(mod => {
      return WebAssembly.instantiate(mod, importObject).then(instance => {
        return {
          module: mod,
          instance: instance
        }
      })
    })
}

before(() => {
  const go = new Go()
  return instantiateStreaming(readFile('./keyaddr.wasm'), go.importObject)
    .then(function (result) {
      go.run(result.instance)
    })
    .then(() => {
      global.Keyaddr = {
        newKey: promisify(KeyaddrNS.newKey),
        wordsToBytes: promisify(KeyaddrNS.wordsToBytes),
        deriveFrom: promisify(KeyaddrNS.deriveFrom),
        ndauAddress: promisify(KeyaddrNS.ndauAddress),
        toPublic: promisify(KeyaddrNS.toPublic),
        child: promisify(KeyaddrNS.child),
        sign: promisify(KeyaddrNS.sign),
        hardenedChild: promisify(KeyaddrNS.hardenedChild),
        wordsFromPrefix: promisify(KeyaddrNS.wordsFromPrefix),
        isPrivate: promisify(KeyaddrNS.isPrivate),
        fromString: promisify(KeyaddrNS.fromString),
        wordsFromBytes: promisify(KeyaddrNS.wordsFromBytes),
        exit: promisify(KeyaddrNS.exit)
      }
    })
})

const language = 'en'
const recoveryPhrase = 'eye eye eye eye eye eye eye eye eye eye eye eye'
const recoveryBytes = 'USolRKiVEqJUSolRKiVEqA=='
const recoveryBytesDifferent = 'USolRKiVEqJUTolRKiVEqA=='
const recoveryPhraseDifferent =
  'eye eye eye eye eye eye eye spell eye eye eye exotic'
const privateKey =
  'npvta8jaftcjebhe9pi57hji5yrt3yc3f2gn3ih56rys38qxt52945vuf8xqu4jfkaaaaaaaaaaaacz6d28v6zwuqm6c7jt4yqcjk4ujvw53jqehafkm5xxvh39jjep58u7pw33dd7cc'
const badPrivateKey = 'foo' + privateKey
const parentPath = `/`
const childPath = `/44'/20036'/100/1`
const firstChildPrivateKey =
  'npvta8jaftcjea9en4tr26txweh3fkejikikubnn7mthymvir3292cquaphxr2egybb9y9asaaaaaecjj4t7hddnv9ebxpym82qugb7j5uagph248cx9wjuttq4y4k9zs43vvnn3z9tw'
const firstChildAddress = 'ndakj49v6nnbdq3yhnf8f2j6ivfzicedvfwtunckivfsw9qt'
const firstChildPublicKey =
  'npuba4jaftckeebyrmpkw4ap7jae22wyb83ncdseuwpvfibunh5fi8id6vs2he86jwieh856caaaaaasjfhkhw6npur6sgxy3r5b4i2hxhqia3w9dm2kz8tgkgf5k5jm88c5htkhhpt8'
const firstGrandchildN = 1
const firstGrandchildPrivateKey =
  'npvta8jaftcjecame82cpnjyjidck3yam94xsixuns994m7i28rwb5tet3pxredtabjzyxuaaaaaagc9jpvb2as73vizj34tcnhgfdum475u34rtmmzdhvfrad8krkhsc8maq9y7avm2'
const msg = 'bmRhdSBpcyBncmVhdAo='
const firstGrandchildSignature =
  'ayjaftcggbcaeidngksig436aeyij65qbu8tq7va2famh4we2f5urbk57v8hg4pj6wbcadmy39j2we2uqsn8rc8rhycznfagqdfrkcf3pkdstmu9xxhkeyi6s78ad42k'
const firstHardenedGrandchildPrivateKey =
  'npvta8jaftcjedampmhj9tybxp78m3dp67fppc53cmezvyu3ree9qnj7ywsvq7cqgbjzyxuiaaaaahx7w9s5zktt9efze3xtaysg2vydrnrpq68wp2cpuu7ep4veibgq3rvtjc2nsuqy'

describe('Keyaddr', () => {
  before(() => {
    global.KeyaddrLogLevel = global.KeyaddrLogLevelDebug
  })

  describe('wordsToBytes', () => {
    it('it should error for incorrect number of arguments', () => {
      expect(Keyaddr.wordsToBytes()).to.eventually.be.rejected
    })

    it('it gets recovery bytes from a recovery phrase', async () => {
      const bytes = await Keyaddr.wordsToBytes(language, recoveryPhrase)
      expect(bytes).to.equal(recoveryBytes)
    })
  })

  describe('newKey', () => {
    it('gets a new key from recovery bytes', async () => {
      const key = await Keyaddr.newKey(recoveryBytes)
      expect(key).to.equal(privateKey)
    })
    it('errors with bad recovery bytes', async () => {
      return await expect(Keyaddr.newKey(recoveryBytes + '42')).to.eventually.be
        .rejected
    })
  })

  describe('deriveFrom', () => {
    it('derives a new key from the root private key', async () => {
      const key = await Keyaddr.deriveFrom(privateKey, parentPath, childPath)
      expect(key).to.equal(firstChildPrivateKey)
    })
    it('errors with missing arguments', async () => {
      return await expect(Keyaddr.deriveFrom()).to.eventually.be.rejected
    })
    it('errors with bad private key', async () => {
      return await expect(
        Keyaddr.deriveFrom(badPrivateKey, parentPath, childPath)
      ).to.eventually.be.rejected
    })
    it('errors with bad parentPath', async () => {
      return await expect(Keyaddr.deriveFrom(privateKey, 'foo', childPath)).to
        .eventually.be.rejected
    })
    it('errors with bad childPath', async () => {
      return await expect(Keyaddr.deriveFrom(privateKey, parentPath, 'foo')).to
        .eventually.be.rejected
    })
    it('errors with test case data', async () => {
      return await expect(
        Keyaddr.deriveFrom('ZWEQAwQFBgcICQoLDA0ODw==', '/', '/1')
      ).to.eventually.be.rejected
    })
  })

  describe('ndauAddress', () => {
    it(`gets the address of the child's private key`, async () => {
      const address = await Keyaddr.ndauAddress(firstChildPrivateKey)
      expect(address).to.equal(firstChildAddress)
    })
    it(`errors with a bad private key`, async () => {
      return await expect(Keyaddr.ndauAddress(badPrivateKey)).to.eventually.be
        .rejected
    })
  })

  describe('toPublic', () => {
    it('gets a public key from a private one', async () => {
      const pubKey = await Keyaddr.toPublic(firstChildPrivateKey)
      expect(pubKey).to.equal(firstChildPublicKey)
    })
    it(`errors with a bad private key`, async () => {
      return await expect(Keyaddr.toPublic(badPrivateKey)).to.eventually.be
        .rejected
    })
  })

  describe('child', () => {
    it(`gets grandchild's private key`, async () => {
      const key = await Keyaddr.child(firstChildPrivateKey, firstGrandchildN)
      expect(key).to.equal(firstGrandchildPrivateKey)
    })
    it(`errors with a bad private key`, async () => {
      return await expect(Keyaddr.child(badPrivateKey)).to.eventually.be
        .rejected
    })
  })

  describe('sign', () => {
    it('signs a message', async () => {
      const sig = await Keyaddr.sign(firstGrandchildPrivateKey, msg)
      expect(sig).to.equal(firstGrandchildSignature)
    })
    it(`errors with a bad private key`, async () => {
      return await expect(Keyaddr.sign(badPrivateKey, msg)).to.eventually.be
        .rejected
    })
  })

  describe('hardenedChild', () => {
    it('creates a hardened child private key', async () => {
      const key = await Keyaddr.hardenedChild(
        firstChildPrivateKey,
        firstGrandchildN
      )
      expect(key).to.equal(firstHardenedGrandchildPrivateKey)
    })
    it(`errors with a bad private key`, async () => {
      return await expect(
        Keyaddr.hardenedChild(badPrivateKey, firstGrandchildN)
      ).to.eventually.be.rejected
    })
  })

  describe('wordsFromPrefix', () => {
    it('gets list of words from a prefix', async () => {
      const words = await Keyaddr.wordsFromPrefix('en', 'gir', 100)
      expect(words).to.equal('giraffe girl')
    })
    it('gets a truncated list of words', async () => {
      const resp = await Keyaddr.wordsFromPrefix('en', 'g', 2)
      expect(resp.split(' ').length).to.equal(2)
    })
  })

  describe('isPrivate', () => {
    it('tests a public key for privacy', async () => {
      const isPrivate = await Keyaddr.isPrivate(firstChildPublicKey)
      expect(isPrivate).to.equal(false)
    })
    it('tests a private key for privacy', async () => {
      const isPrivate = await Keyaddr.isPrivate(firstChildPrivateKey)
      expect(isPrivate).to.equal(true)
    })
    it('tests a private key for privacy', async () => {
      expect(Keyaddr.isPrivate(badPrivateKey)).to.eventually.be.rejected
    })
  })

  describe('fromString', () => {
    it('creates a key from a public key string', async () => {
      const key = await Keyaddr.fromString(firstChildPublicKey)
      expect(key).to.deep.equal({
        key: firstChildPublicKey
      })
    })
    it('creates a key from a private key string', async () => {
      const key = await Keyaddr.fromString(firstChildPrivateKey)
      expect(key).to.deep.equal({
        key: firstChildPrivateKey
      })
    })
    it('errors trying to create a key from a bad string', () => {
      return expect(Keyaddr.fromString(badPrivateKey)).to.eventually.be.rejected
    })
  })

  it('converts bytes to words', async () => {
    const words = await Keyaddr.wordsFromBytes('en', recoveryBytes)
    expect(words).to.equal(recoveryPhrase)
  })
  it('converts bytes to words with different bytes', async () => {
    const words = await Keyaddr.wordsFromBytes('en', recoveryBytesDifferent)
    expect(words).to.equal(recoveryPhraseDifferent)
  })
  it('errors for bad bytes', () => {
    expect(Keyaddr.wordsFromBytes('en', 'foobar')).to.eventually.be.rejected
  })
})

describe('simple memory test', () => {
  it('should not run out of memory', async () => {
    const bytes = await Keyaddr.wordsToBytes(language, recoveryPhrase)
    const key = await Keyaddr.newKey(bytes)
    const oldLogLevel = global.KeyaddrLogLevel // save previous log level
    global.KeyaddrLogLevel = global.KeyaddrLogLevelError // turn off excessive logs
    for (let i = 1; i < 1000; i++) {
      const newKey = await Keyaddr.deriveFrom(key, '/', `/44'/20036'/100/${i}`)
      const addy = await Keyaddr.ndauAddress(newKey)
    }
    global.KeyaddrLogLevel = oldLogLevel // restore log level
    return Promise.resolve()
  })
})

// these have to go on the bottom for...obvious reasons
describe('exiting the wasm application', () => {
  it('exits the wasm program', async () => {
    const msg = await Keyaddr.exit()
    console.log(msg)
  })
  it('should no longer handle function calls', async () => {
    expect(Keyaddr.newKey(recoveryBytes)).to.be.rejected
  })
})
