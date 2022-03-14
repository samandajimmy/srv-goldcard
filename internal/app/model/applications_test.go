package model_test

import (
	"srv-goldcard/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// personal information dummy data
var app = model.Applications{
	Status: "cacing",
}

func TestGetStatusDateKey(t *testing.T) {
	val := app.GetStatusDateKey()
	assert.Equal(t, "", val)
}
