const cryptox = (function () {
    function encrypt(sourceData, secretKey) {
        if (!sourceData) return "";
        let jsonStr = JSON.stringify(sourceData)
        jsonStr = (jsonStr + "").trim()
        if (!jsonStr) return ""
        let base64Str = encodeBase64(jsonStr)
        secretKey = (secretKey || "").trim();
        secretKey = secretKey || genSecretKey()
        const unicodeArray = []
        const base64Array = base64Str.split('')
        const m = base64Array.length
        const secretArray = secretKey.split('')
        const n = secretArray.length
        for (let i = 0; i < m; i++) {
            const unicodeValue = base64Array[i].charCodeAt(0);
            const iv = secretArray[i % n].charCodeAt(0)
            unicodeArray.push(unicodeValue + iv)
        }
        const unicodeStr = unicodeArray.toString()
        const encryptData = encodeBase64(unicodeStr)
        return shuffleString(encryptData)
    }

    function decrypt(encryptString, validJSON, secretKey) {
        if (!encryptString) return "";
        encryptString = (encryptString + "").trim();
        if (!encryptString) return "";
        const base64String = unshuffleString(encryptString)
        const unicodeStr = decodeBase64(base64String);
        const unicodeArray = unicodeStr.split(',').map(Number);
        const m = unicodeArray.length
        secretKey = (secretKey || "").trim();
        const hasScretKey = secretKey && secretKey.length > 0
        secretKey = secretKey || genSecretKey();
        const base64Arr = []
        const secretArray = secretKey.split('')
        const n = secretArray.length
        for (let i = 0; i < m; i++) {
            const unicodeValue = unicodeArray[i];
            const iv = secretArray[i % n].charCodeAt(0)
            base64Arr[i] = String.fromCharCode(unicodeValue - iv);
        }
        const base64Str = base64Arr.join("")
        const decryptData = decodeBase64(base64Str)
        if (validJSON) {
            if (isJSON(decryptData) || hasScretKey) {
                return decryptData
            }
            secretKey = genSecretKey(undefined, true)
            return decrypt(encryptString, false, secretKey)
        }
        return decryptData
    }

    function decodeBase64(base64Str) {
        const binaryStr = atob(base64Str);
        const decoder = new TextDecoder("utf-8");
        const utf8Str = decoder.decode(new Uint8Array([...binaryStr].map(char => char.charCodeAt(0))));
        return utf8Str;
    }

    function encodeBase64(str) {
        const encoder = new TextEncoder();
        const data = encoder.encode(str);
        const base64Chars = btoa(String.fromCharCode.apply(null, data));
        return base64Chars;
    }

    function genSecretKey(date, previousMinute) {
        if (!date) {
            date = new Date();
        }
        if (previousMinute === true) {
            date.setMinutes(date.getMinutes() - 1);
        }
        const timestampInSeconds = Math.floor(date.getTime() / 1000);
        const timestampInMinutes = timestampInSeconds - (timestampInSeconds % 60);
        const s = (timestampInMinutes * 3 / 10).toString();
        return s.split('').reverse().join('') + s;
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

    function shuffleString(text) {
        const characters = text.split('');
        let left = 0;
        let right = characters.length - 1;
        while (left < right) {
            [characters[left], characters[right]] = [characters[right], characters[left]];
            left++;
            right--;
        }
        return characters.join('');
    }

    function unshuffleString(shuffledText) {
        return shuffleString(shuffledText);
    }

    return {
        encrypt: encrypt,
        decrypt: decrypt,
    }
})();

function test() {
    const a = "你好licat";
    const b = cryptox.encrypt(a);
    console.log(b);
    const c = cryptox.decrypt(b);
    console.log(c);
}

test();
