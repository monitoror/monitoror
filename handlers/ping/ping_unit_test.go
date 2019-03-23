package ping

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/jsdidierlaurent/monitowall/renderings"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// ---------- Mocks ---------- //
type (
	PingMock struct{}
)

func (p *PingMock) Ping(hostname string) *HealthCheckResponse {
	return &HealthCheckResponse{
		Status:  SuccessStatus,
		Label:   hostname,
		Message: "1ms",
	}
}

// --------------------------- //

func TestGetPing_unit(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.QueryParams().Set("hostname", "gitlab.com")
	h := NewHandler(&PingMock{})

	var pingJSON = `{"status":"SUCCESS","label":"gitlab.com","message":"1ms"}`

	if assert.NoError(t, h.GetPing(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, pingJSON, strings.TrimSpace(rec.Body.String()))
	}
}
