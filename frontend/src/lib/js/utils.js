export function processProps(props, blocklist) {
    return Object.fromEntries(Object.entries(props).filter(([key, value]) => {
        if (!blocklist.includes(key)) {
            return [key, value]
        }
    }))
}

export function stringToArrayBuffer(string) {
    return new TextEncoder().encode(string)
}
  
export function arrayBufferToHex(byteArray) {
    return [...new Uint8Array(byteArray)].map(x => x.toString(16).padStart(2, '0')).join('');
}

export function hexToArrayBuffer(hex) {
    var typedArray = new Uint8Array(hex.match(/[\da-f]{2}/gi).map((h) => parseInt(h, 16)))
    return typedArray.buffer
}