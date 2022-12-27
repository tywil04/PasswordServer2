import * as utils from "$lib/js/utils.js"

export const config = {
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
    size: 128,
  },
  hash: {
    digest: "SHA-512",
  }
}

function useConfig(testConfig) {
  if (testConfig === null) {
    return config 
  }
  return testConfig
}

export async function importMasterKey(key, optionalConfig=null) {
  let config = useConfig(optionalConfig)
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

export async function generateMasterKey(masterPassword, email, optionalConfig=null) {
  let config = useConfig(optionalConfig)
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

export async function generateMasterHash(masterPassword, masterKey, optionalConfig=null) {
  let config = useConfig(optionalConfig)
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

export async function generateDatabaseKey(optionalConfig) {
  let config = useConfig(optionalConfig)
  return await crypto.subtle.generateKey(
    {
      name: config.databaseKey.encryptionFunction, 
      length: config.databaseKey.size
    },
    true,
    ["encrypt", "decrypt"]
  )
}

export async function protectDatabaseKey(masterKey, databaseKey, optionalConfig=null) {
  let config = useConfig(optionalConfig)
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

export async function unprotectDatabaseKey(masterKey, protectedDatabaseKey, optionalConfig=null) {
  let config = useConfig(optionalConfig)
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

export async function hash(value, optionalConfig=null) {
  let config = useConfig(optionalConfig)
  var buffer = utils.stringToArrayBuffer(value)
  var hashBytes = await crypto.subtle.digest(config.hash.digest, buffer);
  return utils.arrayBufferToHex(hashBytes)
}