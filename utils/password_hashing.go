package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/paiduay/queq-hospital-api/config"
)

// PasswordHashing - Hash Password
// mode: prod, dev
func PasswordHashing(password string) *string {
	h := sha256.New()
	h.Write([]byte(config.Configs.PasswordSalt + password))
	hashedString := fmt.Sprintf("%x", h.Sum(nil))
	return &hashedString
}

// FileHashing - hash a content of file
func FileHashing(key interface{}) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return "", err
	}
	strOfBytes := convert(buf.Bytes()[:])
	strOfByteHashed := PasswordHashing(strOfBytes)
	return *strOfByteHashed, nil
}

// convert - change bytes array to string
func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, "")
}
