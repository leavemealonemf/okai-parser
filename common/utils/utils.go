package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
)

func ReverseBytes(hexString string) []byte {
	bytes, _ := hex.DecodeString(hexString)

	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return bytes
}

func HexToDec(hexString string) int64 {
	dec, _ := strconv.ParseInt(hexString, 16, 64)
	return dec
}

func BytesToHexString(bytes []byte) string {
	encoded := hex.EncodeToString(bytes)
	return encoded
}

func HexToBytes(hexString string) []byte {
	bytes, _ := hex.DecodeString(hexString)
	return bytes
}

func JsonStringify(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Failed to JsonStringify")
	}

	return string(b), nil
}
