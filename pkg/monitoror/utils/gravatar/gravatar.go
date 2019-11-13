package gravatar

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const GravatarURL = "https://www.gravatar.com/avatar/%s?d=blank"

func GetGravatarURL(email string) string {
	return fmt.Sprintf(GravatarURL, getMD5Hash(email))
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	_, _ = hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
