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
	unicodeSlice := encodeUnicode(string(jsonStr))
	m := len(unicodeSlice)

	secretKey = strings.TrimSpace(secretKey)
	hasSecretKey := len(secretKey) > 0
	if !hasSecretKey {
		date := time.Now()
		secretKey = genSecretKey(&date, false)
	}

	secretUnicodesArr := encodeUnicode(secretKey)
	n := len(secretUnicodesArr)
	unicodeArray := make([]int, m)
	for i, base64Unicode := range unicodeSlice {
		iv := secretUnicodesArr[i%n]
		unicodeArray[i] = base64Unicode + iv
	}

	unicodeStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(unicodeArray)), ","), "[]")
	encryString := base64Encode(unicodeStr)
	return shuffleString(encryString)
}

func Decrypt(encryptString string, validJSON bool, secretKey string) string {
	encryptString = strings.TrimSpace(encryptString)
	if encryptString == "" {
		return ""
	}

	base64Str := unshuffleString(encryptString)
	unicodeStr := base64Decode(base64Str)
	unicodeSlice := stringToIntSlice(strings.Split(unicodeStr, ","))

	secretKey = strings.TrimSpace(secretKey)
	hasSecretKey := len(secretKey) > 0
	if !hasSecretKey {
		date := time.Now()
		secretKey = genSecretKey(&date, false)
	}

	secretUnicodesArr := encodeUnicode(secretKey)
	n := len(secretUnicodesArr)

	var stringArr []string
	for i, base64Unicode := range unicodeSlice {
		iv := secretUnicodesArr[i%n]
		stringArr = append(stringArr, decodeUnicode(base64Unicode-iv))
	}

	decryptData := strings.Join(stringArr, "")

	if validJSON {
		if isJSON([]byte(decryptData)) || hasSecretKey {
			return string(decryptData)
		}
		secretKey = genSecretKey(nil, true)
		return Decrypt(encryptString, false, secretKey)
	}

	return string(decryptData)
}

func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) string {
	data, _ := base64.StdEncoding.DecodeString(str)
	return string(data)
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
	timestampInMinutes = timestampInMinutes * int64(date.Minute()) / 10
	s := fmt.Sprint(timestampInMinutes)
	return reverseString(s) + s
}

func isJSON(data []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(data, &js) == nil
}

func encodeUnicode(str string) []int {
	var unicodes []int
	for _, char := range str {
		unicodes = append(unicodes, int(char))
	}

	return unicodes
}

func decodeUnicode(unicode int) string {
	return string(rune(unicode))
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

func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
