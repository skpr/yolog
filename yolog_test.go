package yolog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerLog(t *testing.T) {
	// Create a new logger that will always sample logs.
	logger := NewLogger("test", WithSampleRate(1.0))

	logger.SetAttr("skpr", "rocks").SetErrorCode(500)

	logger.SetAttrs("server", "backup", "rpc", "create")

	var buf bytes.Buffer

	logger.Log(&buf)

	assert.Contains(t, buf.String(), `"stream":"test"`)
	assert.Contains(t, buf.String(), `"error_code":500`)
	assert.Contains(t, buf.String(), `"skpr":"rocks"`)
	assert.Contains(t, buf.String(), `"server":"backup"`)
	assert.Contains(t, buf.String(), `"rpc":"create"`)
}
