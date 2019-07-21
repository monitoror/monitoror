package gravatar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGravatarUrl(t *testing.T) {
	gravatarUrl := GetGravatarUrl("test@gmail.com")
	assert.Equal(t, "https://www.gravatar.com/avatar/1aedb8d9dc4751e229a335e371db8058?d=blank", gravatarUrl)
}
