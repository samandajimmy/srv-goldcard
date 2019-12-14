package models_test

import (
	"gade/srv-goldcard/models"
	"testing"
)

func TestNewBriAPI(t *testing.T) {
	models.NewBriAPI()
}

func mltplFunc(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}
