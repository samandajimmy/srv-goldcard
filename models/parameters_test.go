package models_test

import (
	"gade/srv-goldcard/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringNameFormatter(t *testing.T) {
	str := models.StringNameFormatter("DEVRI HARYANTO ARIARIARiop MANULLANG", 18, true)

	assert.Equal(t, "SAMANDA RASU", str)
}

func TestRemappAddress(t *testing.T) {
	pl := models.AddressData{
		AddressLine1: "Jl. Harsono RM No 2, Gedung IT BRI, Ragunan",
	}
	str, err := models.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
	assert.Equal(t, "SAMANDA RASU", err)
}
