<?php

function encryptData($data)
{
    $input = json_encode($data);
    $secretData = encryptString($input);
    $encodedData = base64_encode($secretData);
    return $encodedData;
}

function decryptBase64Data($data)
{
    $decodedData = base64_decode($data);
    if ($decodedData === false) {
        // 不是 Base64 编码
        return null;
    }
    $date = new DateTime();
    $secretKey = genSecretKey($date, false);
    return decryptString($decodedData, $secretKey);
}

function decryptJsonBase64Data($data)
{
    $decodedData = base64_decode($data);
    if ($decodedData === false) {
        // 不是 Base64 编码
        return null;
    }
    return decryptJSONString($decodedData);
}

function decryptJSONString($input = "", $secretKey = "")
{
    $date = new DateTime();
    if ($secretKey == "") {
        $secretKey = genSecretKey($date, false);
    }
    $decrypted = decryptString($input, $secretKey);
    if (!isJSON($decrypted)) {
        $decrypted = decryptString($input, genSecretKey($date, true));
        if (!isJSON($decrypted)) {
            return null;
        }
    }
    return $decrypted;
}

function decryptString($input = "", $secretKey = "")
{
    if ($secretKey == "") {
        $date = new DateTime();
        $secretKey = genSecretKey($date, false);
    }
    $n = strlen($secretKey);
    $decrypted = "";
    for ($i = 0; $i < strlen($input); $i++) {
        $charCode = ord($input[$i]);
        $iv = ord($secretKey[$i % $n]);
        $decryptedCharCode = $charCode - $iv;
        $decrypted .= chr($decryptedCharCode);
    }
    return $decrypted;
}

function encryptString($input = "", $secretKey = "")
{
    if ($secretKey == "") {
        $secretKey = genSecretKey(null, false);
    }
    $n = strlen($secretKey);
    $encrypted = "";
    for ($i = 0; $i < strlen($input); $i++) {
        $charCode = ord($input[$i]);
        $iv = ord($secretKey[$i % $n]);
        $encryptedCharCode = $charCode + $iv;
        $encrypted .= chr($encryptedCharCode);
    }
    return $encrypted;
}

function getPreviousMinuteDate($date = null)
{
    if ($date == null) {
        $now = new DateTime();
        $date = $now;
    }
    $interval = new DateInterval('PT1M');
    $date->sub($interval);
    return $date;
}

function genSecretKey($date = null, $previousMinute = false)
{
    if ($date == null) {
        $now = new DateTime();
        $date = $now;
    }
    if ($previousMinute) {
        $date = getPreviousMinuteDate($date);
    }
    $timestampInSeconds = $date->getTimestamp();
    $timestampInMinutes = $timestampInSeconds - ($timestampInSeconds % 60);
    return strval($timestampInMinutes);
}

function isJSON($data = "")
{
    json_decode($data);
    return json_last_error() === JSON_ERROR_NONE;
}

function test()
{
    $encryData = encryptString("你好licat");
    print($encryData);
    echo "\n";
    $decryData = decryptString($encryData);
    print($decryData);
}
