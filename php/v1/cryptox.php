<?php
class Cryptox
{
    public static function encrypt($sourceData, $secretKey = "")
    {
        if (empty($sourceData)) return "";
        $jsonStr = json_encode($sourceData, JSON_UNESCAPED_UNICODE);
        $jsonStr = trim($jsonStr);
        if (empty($jsonStr)) return "";
        $base64Arr = self::stringToArray($jsonStr);
        $m = count($base64Arr);
        $secretKey = trim($secretKey);
        $hasSecretKey = $secretKey && strlen($secretKey) > 0;
        if (!$hasSecretKey) $secretKey = self::genSecretKey();
        echo $secretKey;
        $secretKeyArr = self::stringToArray($secretKey);
        $n = count($secretKeyArr);
        $unicodeArray = [];
        for ($i = 0; $i < $m; $i++) {
            $unicodeValue = self::encodeUnicode($base64Arr[$i]);
            $iv = self::encodeUnicode($secretKeyArr[$i % $n]);
            $unicodeArray[] = $unicodeValue + $iv;
        }

        $unicodeStr = implode(",", $unicodeArray);
        $encryptData = base64_encode($unicodeStr);
        return self::shuffleString($encryptData);
    }

    public static function decrypt($encryptString, $validJSON = false, $secretKey = "")
    {
        $encryptString = trim((string)$encryptString);
        if (empty($encryptString)) return "";

        $base64Str = self::unshuffleString($encryptString);
        $unicodeStr = base64_decode($base64Str);
        $unicodeArray = array_map('intval', explode(',', $unicodeStr));
        $m = count($unicodeArray);

        $secretKey = trim($secretKey);
        $hasSecretKey = $secretKey && strlen($secretKey) > 0;
        if (!$hasSecretKey) $secretKey = self::genSecretKey();

        $secretKeyArr = self::stringToArray($secretKey);
        $n = count($secretKeyArr);

        $base64Arr = [];
        for ($i = 0; $i < $m; $i++) {
            $unicodeValue = $unicodeArray[$i];
            $iv = self::encodeUnicode($secretKeyArr[$i % $n]);
            $base64Arr[$i] = self::decodeUnicode($unicodeValue - $iv);
        }

        $decryptData = implode("", $base64Arr);

        if ($validJSON) {
            if (self::isJSON($decryptData) || $hasSecretKey) {
                return $decryptData;
            }
            $secretKey = self::genSecretKey(null, true);
            return self::decrypt($encryptString, false, $secretKey);
        }

        return $decryptData;
    }

    private static function getPreviousMinuteDate($timestamp = null)
    {
        if ($timestamp == null) {
            $timestamp = time();
        }
        return strtotime("-1 minute", $timestamp);
    }

    private static function genSecretKey($timestamp = null, $previousMinute = false)
    {
        if ($timestamp == null) {
            $timestamp = time();
        }
        if ($previousMinute) {
            $timestamp = self::getPreviousMinuteDate($timestamp);
        }
        $timestampInMinutes = $timestamp - ($timestamp % 60);
        $minutes = date('i', $timestamp);
        $timestampInMinutes = $timestampInMinutes * $minutes / 10;
        $s = strval($timestampInMinutes);
        return strrev($s) . $s;
    }

    private static function isJSON($data = "")
    {
        json_decode($data);
        return json_last_error() === JSON_ERROR_NONE;
    }

    private static function shuffleString($text)
    {
        $characters = str_split($text);
        $left = 0;
        $right = count($characters) - 1;
        while ($left < $right) {
            [$characters[$left], $characters[$right]] = [$characters[$right], $characters[$left]];
            $left++;
            $right--;
        }
        return implode('', $characters);
    }

    private static function unshuffleString($shuffledText)
    {
        return self::shuffleString($shuffledText);
    }

    private static function encodeUnicode($str)
    {
        return mb_ord($str, 'UTF-8');
    }

    private static function decodeUnicode($unicode)
    {
        return mb_chr($unicode, 'UTF-8');
    }

    private static function stringToArray($str)
    {
        return mb_str_split($str, 1, 'UTF-8');
    }
}
