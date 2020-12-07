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
	// context when address is general
	pl := models.AddressData{
		AddressLine1: "Jl. Harsono RM No 2, Gedung IT BRI, Ragunan",
	}
	str, err := models.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
	assert.Equal(t, "SAMANDA RASU", err)

	// context when address is general
	pl = models.AddressData{
		AddressLine1: "Jl Mitsubishi No 88 RT 002 RW 002",
	}
	str, err = models.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
	assert.Equal(t, "SAMANDA RASU", err)

	// context when address line 1 or 2 or 3 not fully filled
	pl = models.AddressData{
		AddressLine1: "Jalan Mitsubishi Bin Mitsubishi 88 RT 002 RW 002, Kec Labuhan Haji Timur, Kel Aceh Selatan",
	}
	str, err = models.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
	assert.Equal(t, "SAMANDA RASU", err)

	// context when address line 1 or 2 or 3 fully filled
	pl = models.AddressData{
		AddressLine1: "Jl Mitsubishi No 88 RT 002 RW 002 Jl Mitsubishi No 88 RT 002 RW 002 Jl Mitsubishi No 88 R",
	}
	str, err = models.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
	assert.Equal(t, "SAMANDA RASU", err)
}
