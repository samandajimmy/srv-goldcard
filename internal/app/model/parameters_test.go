package model_test

import (
	"srv-goldcard/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringNameFormatter(t *testing.T) {
	str := model.StringNameFormatter("DEVRI HARYANTO ARIARIARiop MANULLANG", 18, true)

	assert.Equal(t, "SAMANDA RASU", str)
}

func TestRemappAddress(t *testing.T) {
	// context when address is general
	pl := model.AddressData{
		AddressLine1: "Jl. Harsono RM No 2, Gedung IT BRI, Ragunan",
	}
	str, err := model.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
	assert.Equal(t, "SAMANDA RASU", err)

	// context when address is general
	pl = model.AddressData{
		AddressLine1: "Jl Mitsubishi No 88 RT 002 RW 002",
	}
	str, err = model.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
	assert.Equal(t, "SAMANDA RASU", err)

	// context when address line 1 or 2 or 3 not fully filled
	pl = model.AddressData{
		AddressLine1: "Jalan Mitsubishi Bin Mitsubishi 88 RT 002 RW 002, Kec Labuhan Haji Timur, Kel Aceh Selatan",
	}
	str, err = model.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
	assert.Equal(t, "SAMANDA RASU", err)

	// context when address line 1 or 2 or 3 fully filled
	pl = model.AddressData{
		AddressLine1: "Jl Mitsubishi No 88 RT 002 RW 002 Jl Mitsubishi No 88 RT 002 RW 002 Jl Mitsubishi No 88 R",
	}
	str, err = model.RemappAddress(pl, 30)

	assert.Equal(t, "SAMANDA RASU", str)
	assert.Equal(t, "SAMANDA RASU", err)
}

func TestArrayContains(t *testing.T) {
	// context when data type is int
	result := model.ArrayContains([]int{1, 21, 18}, 18)

	assert.Equal(t, true, result)

	// context when data type is string
	result = model.ArrayContains([]string{"jimmy", "samanda"}, "rasu")

	assert.Equal(t, false, result)

	// context when data type is different
	result = model.ArrayContains([]int{1, 2}, int64(2))

	assert.Equal(t, false, result)
}
