package main

import (
	"testing"
)

func TestCryptox(t *testing.T) {
	data := "你好licat"
	encryptStr := Encrypt(data, "")
	t.Log(encryptStr)
	decryptStr := Decrypt(encryptStr, false, "")
	t.Log(decryptStr)
}
