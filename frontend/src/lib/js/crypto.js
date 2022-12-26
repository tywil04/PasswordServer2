import * as utils from "$lib/js/utils.js"

const config = {
  masterKey: {
    keyFunction: "PBKDF2",
    digest: "SHA-512",
    iterations: 100000,
  },
  masterHash: {
    keyFunction: "PBKDF2",
    digest: "SHA-512",
    iterations: 100000,
    size: 512,
  },
  databaseKey: {
    encryptionFunction: "AES-CBC",
    size: 256,
  },
  hash: {
    digest: "SHA-512",
  }
}

export async function importMasterKey(key) {
  return await crypto.subtle.importKey(
    "raw",
    key,
    {
      name: config.databaseKey.encryptionFunction,
      length: config.databaseKey.size,
    },
    true,
    ["wrapKey", "unwrapKey"]
  )
}

export async function exportKey(key) {
  return await crypto.subtle.exportKey("raw", key)
}

export async function generateMasterKey(masterPassword, email) {
  let masterPasswordKey = await crypto.subtle.importKey(
    "raw",
    utils.stringToArrayBuffer(masterPassword),
    config.masterKey.keyFunction,
    false,
    ["deriveKey", "deriveBits"],
  )

  return await crypto.subtle.deriveKey(
    {
      name: config.masterKey.keyFunction,
      hash: config.masterKey.digest,
      salt: utils.stringToArrayBuffer(email),
      iterations: config.masterKey.iterations,
    },
    masterPasswordKey,
    {
      name: config.databaseKey.encryptionFunction,
      length: config.databaseKey.size,
    },
    true,
    ["wrapKey", "unwrapKey"]
  );
}

export async function generateMasterHash(masterPassword, masterKey) {
  let masterKeyKey = await crypto.subtle.importKey(
    "raw",
    utils.stringToArrayBuffer(masterKey),
    config.masterHash.keyFunction,
    false,
    ["deriveKey", "deriveBits"],
  )

  return await crypto.subtle.deriveBits(
    {
      name: config.masterHash.keyFunction,
      hash: config.masterHash.digest,
      salt: utils.stringToArrayBuffer(masterPassword),
      iterations: config.masterHash.iterations,
    },
    masterKeyKey,
    config.masterHash.size,
  )
}

export async function generateDatabaseKey() {
  return await crypto.subtle.generateKey(
    {
      name: config.databaseKey.encryptionFunction, 
      length: config.databaseKey.size
    },
    true,
    ["encrypt", "decrypt"]
  )
}

export async function protectDatabaseKey(masterKey, databaseKey) {
  const iv = crypto.getRandomValues(new Uint8Array(16))
  const protectedDatabaseKey = await crypto.subtle.wrapKey(
    "raw", 
    databaseKey, 
    masterKey, 
    {
      name: config.databaseKey.encryptionFunction, 
      iv
    }
  )
  return [iv, protectedDatabaseKey]
}

export async function unprotectDatabaseKey(masterKey, protectedDatabaseKey) {
  const [iv, wrappedDatabaseKey] = protectedDatabaseKey
  return await crypto.subtle.unwrapKey(
    "raw", 
    wrappedDatabaseKey, 
    masterKey, 
    {
      name: config.databaseKey.encryptionFunction, 
      iv
    },
    {
      name: config.databaseKey.encryptionFunction, 
      length: config.databaseKey.size
    },
    false,
    ["encrypt", "decrypt"]
  )
}

export function randomUUID() {
  return crypto.randomUUID()
}

export async function hash(value) {
  var buffer = utils.stringToArrayBuffer(value)
  var hashBytes = await crypto.subtle.digest(config.hash.digest, buffer);
  return utils.arrayBufferToHex(hashBytes)
}