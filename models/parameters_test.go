package models_test

import (
	"gade/srv-goldcard/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringNameFormatter(t *testing.T) {
	str := models.StringNameFormatter("SAMANDA RASU", 14, false)

	assert.Equal(t, "SAMANDA RASU", str)
}
