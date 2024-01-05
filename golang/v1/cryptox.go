package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Encrypt(sourceData interface{}, secretKey string) string {
	if sourceData == nil {
		return ""
	}

	jsonStr, _ := json.Marshal(sourceData)
	base64Str := base64.StdEncoding.EncodeToString(jsonStr)
	unicodeSlice := getUnicodeCodes(base64Str)
	m := len(unicodeSlice)

	secretKey = strings.TrimSpace(secretKey)
	hasSecretKey := len(secretKey) > 0
	if !hasSecretKey {
		date := time.Now()
		secretKey = genSecretKey(&date, false)
	}

	secretUnicodesArr := getUnicodeCodes(secretKey)
	n := len(secretUnicodesArr)
	unicodeArray := make([]int, m)
	for i, base64Unicode := range unicodeSlice {
		iv := secretUnicodesArr[i%n]
		unicodeArray[i] = base64Unicode + iv
	}

	unicodeStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(unicodeArray)), ","), "[]")
	encryptData := base64.StdEncoding.EncodeToString([]byte(unicodeStr))
	return shuffleString(encryptData)
}

func Decrypt(encryptString string, validJSON bool, secretKey string) string {
	if encryptString == "" {
		return ""
	}

	encryptString = strings.TrimSpace(encryptString)
	if encryptString == "" {
		return ""
	}

	base64Str := unshuffleString(encryptString)
	unicodeStr, _ := base64.StdEncoding.DecodeString(base64Str)
	unicodeSlice := stringToIntSlice(strings.Split(string(unicodeStr), ","))

	m := len(unicodeSlice)

	secretKey = strings.TrimSpace(secretKey)
	hasSecretKey := len(secretKey) > 0
	if !hasSecretKey {
		date := time.Now()
		secretKey = genSecretKey(&date, false)
	}

	secretUnicodesArr := getUnicodeCodes(secretKey)
	n := len(secretUnicodesArr)

	base64Arr := make([]byte, m)
	for i, base64Unicode := range unicodeSlice {
		iv := secretUnicodesArr[i%n]
		base64Arr[i] = byte(base64Unicode - iv)
	}

	base64Str = string(base64Arr)
	decryptData, _ := base64.StdEncoding.DecodeString(base64Str)

	if validJSON {
		if isJSON(decryptData) || hasSecretKey {
			return string(decryptData)
		}
		secretKey = genSecretKey(nil, true)
		return Decrypt(encryptString, false, secretKey)
	}

	return string(decryptData)
}

func getPreviousMinuteDate(date *time.Time) time.Time {
	if date == nil {
		now := time.Now()
		date = &now
	}
	return date.Add(time.Minute * -1)
}

func genSecretKey(date *time.Time, previousMinute bool) string {
	if date == nil {
		now := time.Now()
		date = &now
	}
	if previousMinute {
		*date = getPreviousMinuteDate(date)
	}
	timestampInSeconds := date.Unix()
	timestampInMinutes := timestampInSeconds - (timestampInSeconds % 60)
	return fmt.Sprint(timestampInMinutes)
}

func isJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}

func getUnicodeCodes(str string) []int {
	unicode := make([]int, len(str))

	for i, char := range str {
		unicode[i] = int(char)
	}

	return unicode
}

func shuffleString(text string) string {
	characters := strings.Split(text, "")
	left := 0
	right := len(characters) - 1

	for left < right {
		characters[left], characters[right] = characters[right], characters[left]
		left++
		right--
	}

	return strings.Join(characters, "")
}

func unshuffleString(shuffledText string) string {
	return shuffleString(shuffledText)
}

func stringToIntSlice(strSlice []string) []int {
	intSlice := make([]int, len(strSlice))
	for i, str := range strSlice {
		intValue, _ := strconv.Atoi(str)
		intSlice[i] = intValue
	}
	return intSlice
}
