package main

import (
	"fmt"
	"testing"
	"time"
)

func TestCryptox(t *testing.T) {
	data := "李"
	encryptStr := Encrypt(data, "")
	t.Log(encryptStr)
	decryptStr := Decrypt(encryptStr, false, "")
	t.Log(decryptStr)
}

func TestGenKey(t *testing.T) {
	date := time.Now()
	timestampInSeconds := date.Unix()
	fmt.Println(timestampInSeconds)
	fmt.Println(date.Minute())
	timestampInMinutes := timestampInSeconds - (timestampInSeconds % 60)
	timestampInMinutes = timestampInMinutes * int64(date.Minute()) / 10
	s := fmt.Sprint(timestampInMinutes)
	fmt.Println(s)
}
