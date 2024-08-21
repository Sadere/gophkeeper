package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestZapLoggerInit(t *testing.T) {
	zapLogger, err := NewZapLogger("debug")

	assert.NotNil(t, zapLogger)
	assert.NoError(t, err)

	assert.Equal(t, zapLogger.Level(), zap.DebugLevel)
}

func TestZapLoggerInitInvalidLevel(t *testing.T) {
	zapLogger, err := NewZapLogger("invalid")

	assert.Nil(t, zapLogger)
	assert.Error(t, err)
}
