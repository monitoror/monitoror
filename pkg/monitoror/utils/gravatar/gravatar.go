package gravatar

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const GRAVATAR_URL = "https://www.gravatar.com/avatar/%s?d=blank"

func GetGravatarUrl(email string) string {
	return fmt.Sprintf(GRAVATAR_URL, getMD5Hash(email))
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
