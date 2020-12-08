package models_test

import (
	"gade/srv-goldcard/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertMoneyToGold(t *testing.T) {
	card := models.Card{}
	gold := card.ConvertMoneyToGold(int64(15000000), int64(765741))

	assert.Equal(t, float64(20.8392), gold)
}

func TestSetGoldLimit(t *testing.T) {
	card := models.Card{}
	str := card.SetGoldLimit(int64(5000000), int64(1100000))

	assert.Equal(t, float64(6.0445), str)
}
