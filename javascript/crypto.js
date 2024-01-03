// import { Buffer } from 'buffer';

function encryptData(data) {
    const input = JSON.stringify(data)
    return encryptString(input)
}

function decryptJSONString(input, date) {
    if (!input) return "";
    let decrypted = decryptString(input, genSecretKey(date || new Date()))
    if (!isJSON(decrypted)) {
        if (date) {
            return ""
        }
        decrypted = decryptString(input, genSecretKey(date, true))
        if (!isJSON(decrypted)) {
            return ""
        }
    }
    return decrypted;
}

function decodeBase64(base64String) {
    if (typeof window === "undefined" && typeof Buffer !== 'undefined') {
        const buffer = Buffer.from(base64String, 'base64');
        return buffer.toString('utf-8');
    }
    return window.atob(base64String)
}

function encodeBase64(str) {
    if (typeof window === "undefined" && typeof Buffer !== 'undefined') {
        return Buffer.from(str, 'utf-8').toString('base64');
    }
    return window.btoa(str)
}

function encryptString(input, secretKey) {
    if (!input) return ""
    secretKey = secretKey || genSecretKey()
    const n = secretKey.length
    const byteArray = []
    for (let i = 0; i < input.length; i++) {
        const charCode = input.charCodeAt(i);
        const iv = secretKey.charCodeAt(i % n)
        const encryptedCharCode = charCode + iv;
        byteArray.push(encryptedCharCode)
    }
    const unicodeStr = byteArray.toString()
    const encryptData = encodeBase64(unicodeStr)
    return encryptData
}

function decryptString(input, secretKey) {
    if (!input) return ""
    secretKey = secretKey || genSecretKey()
    input = decodeBase64(input)
    const byteArray = input.split(',').map(Number);
    const n = secretKey.length;
    const res = []
    for (let i = 0; i < byteArray.length; i++) {
        const byte = byteArray[i];
        const iv = secretKey.charCodeAt(i % n);
        const decryptedByte = byte - iv;
        res[i] = String.fromCharCode(decryptedByte);
    }

    return res.join("")
}

function genSecretKey(date, previousMinute) {
    if (!date) {
        date = new Date();
    }
    if (previousMinute === true) {
        date = getPreviousMinuteDate(date)
    }
    const timestampInSeconds = Math.floor(date.getTime() / 1000);
    const timestampInMinutes = timestampInSeconds - (timestampInSeconds % 60);
    const s = (timestampInMinutes * 3 / 10).toString();

    return reverseString(s) + s;
}

function reverseString(input) {
    return input.split('').reverse().join('');
}

function isJSON(str) {
    try {
        JSON.parse(str);
        return true;
    } catch (e) {
        return false;
    }
}

function getPreviousMinuteDate(date) {
    const now = date || new Date();
    now.setMinutes(now.getMinutes() - 1);
    return now;
}