<?php
class Cryptox
{
    public static function encrypt($sourceData, $secretKey = "")
    {
        if (empty($sourceData)) return "";
        $jsonStr = json_encode($sourceData, JSON_UNESCAPED_UNICODE);
        $jsonStr = trim($jsonStr);
        if (empty($jsonStr)) return "";
        $base64Str = base64_encode($jsonStr);
        $base64Arr = str_split($base64Str);
        $m = count($base64Arr);

        $secretKey = trim($secretKey);
        $hasSecretKey = $secretKey && strlen($secretKey) > 0;
        if (!$hasSecretKey) $secretKey = self::genSecretKey();
        $secretKeyArr = str_split($secretKey);
        $n = count($secretKeyArr);
        $unicodeArray = [];
        for ($i = 0; $i < $m; $i++) {
            $unicodeValue = ord($base64Arr[$i]);
            $iv = ord($secretKeyArr[$i % $n]);
            $unicodeArray[] = $unicodeValue + $iv;
        }

        $unicodeStr = implode(",", $unicodeArray);
        $encryptData = base64_encode($unicodeStr);
        return self::shuffleString($encryptData);
    }

    public static function decrypt($encryptString, $validJSON = false, $secretKey = "")
    {
        if (empty($encryptString)) return "";
        $encryptString = trim((string)$encryptString);
        if (empty($encryptString)) return "";

        $base64Str = self::unshuffleString($encryptString);
        $unicodeStr = base64_decode($base64Str);
        $unicodeArray = array_map('intval', explode(',', $unicodeStr));
        $m = count($unicodeArray);

        $secretKey = trim($secretKey);
        $hasSecretKey = $secretKey && strlen($secretKey) > 0;
        if (!$hasSecretKey) $secretKey = self::genSecretKey();

        $secretKeyArr = str_split($secretKey);
        $n = count($secretKeyArr);

        $base64Arr = [];
        for ($i = 0; $i < $m; $i++) {
            $unicodeValue = $unicodeArray[$i];
            $iv = ord($secretKeyArr[$i % $n]);
            $base64Arr[$i] = chr($unicodeValue - $iv);
        }

        $base64Str = implode("", $base64Arr);
        $decryptData = base64_decode($base64Str);

        if ($validJSON) {
            if (self::isJSON($decryptData) || $hasSecretKey) {
                return $decryptData;
            }
            $secretKey = self::genSecretKey(null, true);
            return self::decrypt($encryptString, false, $secretKey);
        }

        return $decryptData;
    }

    private static function getPreviousMinuteDate($date = null)
    {
        if ($date == null) {
            $now = new DateTime();
            $date = $now;
        }
        $interval = new DateInterval('PT1M');
        $date->sub($interval);
        return $date;
    }

    private static function genSecretKey($date = null, $previousMinute = false)
    {
        if ($date == null) {
            $now = new DateTime();
            $date = $now;
        }
        if ($previousMinute) {
            $date = self::getPreviousMinuteDate($date);
        }
        $timestampInSeconds = $date->getTimestamp();
        $timestampInMinutes = $timestampInSeconds - ($timestampInSeconds % 60);
        return strval($timestampInMinutes);
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
}
