import { expect } from 'chai'
import fs from 'fs'
import { promisify } from 'util'
require('./wasm_exec')
const readFile = promisify(fs.readFile)

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

before(done => {
  const go = new Go()
  instantiateStreaming(readFile('./keyaddr.wasm'), go.importObject)
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
        newKey: promisify(KeyaddrNS.newKey),
        exit: promisify(KeyaddrNS.exit)
      }
      done()
    })
    .catch(err => {
      console.log('something went wrong loading', err)
      done()
    })
})

const language = 'en'
const recoveryPhrase = 'eye eye eye eye eye eye eye eye eye eye eye eye'
const recoveryBytes = 'USolRKiVEqJUSolRKiVEqA=='
const privateKey =
  'npvta8jaftcjebhe9pi57hji5yrt3yc3f2gn3ih56rys38qxt52945vuf8xqu4jfkaaaaaaaaaaaacz6d28v6zwuqm6c7jt4yqcjk4ujvw53jqehafkm5xxvh39jjep58u7pw33dd7cc'
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
  it('it should error for incorrect number of arguments', done => {
    Keyaddr.wordsToBytes().catch(err => {
      expect(err).to.be.ok
      done()
    })
  })

  it('it gets recovery bytes from a recovery phrase', done => {
    Keyaddr.wordsToBytes(language, recoveryPhrase)
      .then(resp => {
        expect(resp).to.equal(recoveryBytes)
        done()
      })
      .catch(err => done(new Error(err)))
  })
  it('gets a new key from recovery bytes', done => {
    Keyaddr.newKey(recoveryBytes)
      .then(resp => {
        expect(resp).to.equal(privateKey)
        done()
      })
      .catch(err => done(new Error(err)))
  })

  it('derives a new address from the root private key', done => {
    Keyaddr.deriveFrom(privateKey, parentPath, childPath)
      .then(resp => {
        expect(resp).to.equal(firstChildPrivateKey)
        done()
      })
      .catch(err => done(new Error(err)))
  })

  it(`gets the address of the child's private key`, done => {
    Keyaddr.ndauAddress(firstChildPrivateKey)
      .then(resp => {
        expect(resp).to.equal(firstChildAddress)
        done()
      })
      .catch(err => done(new Error(err)))
  })

  it('gets a public key from a private one', done => {
    Keyaddr.toPublic(firstChildPrivateKey)
      .then(resp => {
        expect(resp).to.equal(firstChildPublicKey)
        done()
      })
      .catch(err => done(new Error(err)))
  })

  it(`gets grandchild's private key`, done => {
    Keyaddr.child(firstChildPrivateKey, firstGrandchildN)
      .then(resp => {
        expect(resp).to.equal(firstGrandchildPrivateKey)
        done()
      })
      .catch(err => done(new Error(err)))
  })

  it('signs a message', done => {
    Keyaddr.sign(firstGrandchildPrivateKey, msg)
      .then(resp => {
        expect(resp).to.equal(firstGrandchildSignature)
        done()
      })
      .catch(err => done(new Error(err)))
  })

  it('creates a hardened child private key', done => {
    Keyaddr.hardenedChild(firstChildPrivateKey, firstGrandchildN)
      .then(resp => {
        expect(resp).to.equal(firstHardenedGrandchildPrivateKey)
        done()
      })
      .catch(err => done(new Error(err)))
  })

  it('exits the wasm program', done => {
    Keyaddr.exit()
      .then(resp => {
        console.log(resp)
        done()
      })
      .catch(err => done(new Error(err)))
  })
  it('should no longer handle function calls', done => {
    Keyaddr.newKey(recoveryBytes)
      .then(resp => {
        expect(true).to.equal(false) // should never get here
        done()
      })
      .catch(err => {
        expect(err).to.be.ok
        done()
      })
  })
})
