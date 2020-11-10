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

func TestJoinAddress(t *testing.T) {
	pl := models.AddressData{
		AddressLine1: "Probolinggo Jl. Panglima Sudirman No. 75 Kota Probolinggo Jawa Timur",
		Subdistrict:  "KARANG TENGAH",
		Village:      "TUNJUNGSEKAR",
	}
	str := models.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
}
