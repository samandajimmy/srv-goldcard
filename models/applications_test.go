package models_test

import (
	"gade/srv-goldcard/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// personal information dummy data
var app = models.Applications{
	Status: "cacing",
}

func TestGetStatusDateKey(t *testing.T) {
	val := app.GetStatusDateKey()
	assert.Equal(t, "", val)
}
