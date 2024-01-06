<?php

include_once __DIR__ . '/Cryptox.php';

$source = "李";
$encryptData = Cryptox::encrypt($source);
echo "密文:" . $encryptData . "\n";
$decodedData = Cryptox::decrypt($encryptData);
echo "明文:" . $decodedData . "\n";
