package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	_, _ = hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
