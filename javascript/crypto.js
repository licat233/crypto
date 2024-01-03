function decodeBase64(base64String) {
    if (typeof window === "undefined" && typeof Buffer !== 'undefined') {
        return Buffer.from(base64String, 'base64').toString('utf-8');
    }
    return window.atob(base64String)
}

function encodeBase64(str) {
    if (typeof window === "undefined" && typeof Buffer !== 'undefined') {
        return Buffer.from(str, 'utf-8').toString('base64');
    }

    return window.btoa(str)
}

function encryptString(inputString, secretKey) {
    if (!inputString) return ""
    secretKey = secretKey || genSecretKey()
    const n = secretKey.length
    const byteArray = []
    for (let i = 0; i < inputString.length; i++) {
        const charCode = inputString.charCodeAt(i);
        const iv = secretKey.charCodeAt(i % n)
        const encryptedCharCode = charCode + iv;
        byteArray.push(encryptedCharCode)
    }
    const unicodeStr = byteArray.toString()
    const encryptData = encodeBase64(unicodeStr)
    return encryptData
}

function decryptString(inputString, secretKey) {
    if (!inputString) return ""
    secretKey = secretKey || genSecretKey()
    inputString = decodeBase64(inputString)
    const byteArray = inputString.split(',').map(Number);
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

function getPreviousMinuteDate(date) {
    const now = date || new Date();
    now.setMinutes(now.getMinutes() - 1);
    return now;
}

function isJSON(str) {
    try {
        JSON.parse(str);
        return true;
    } catch (e) {
        return false;
    }
}

function reverseString(input) {
    return input.split('').reverse().join('');
}

function genSecretKey(date, isPreviousMinute) {
    if (!date) {
        date = new Date();
    }
    if (isPreviousMinute === true) {
        date = getPreviousMinuteDate(date)
    }
    const timestampInSeconds = Math.floor(date.getTime() / 1000);
    const timestampInMinutes = timestampInSeconds - (timestampInSeconds % 60);
    const s = (timestampInMinutes * 3 / 10).toString();

    return reverseString(s) + s;
}

//test
function test() {
    var str = "你好licat";
    var encryptData = encryptString(str);
    console.log(encryptData);

    var decryptData = decryptString(encryptData);
    console.log(decryptData);
}