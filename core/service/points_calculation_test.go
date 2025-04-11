package service

import (
	"github.com/RA341/receipt-processor-challenge/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculations1(t *testing.T) {
	rec := models.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			}, {
				ShortDescription: "Emils Cheese Pizza",
				Price:            "12.25",
			}, {
				ShortDescription: "Knorr Creamy Chicken",
				Price:            "1.26",
			}, {
				ShortDescription: "Doritos Nacho Cheese",
				Price:            "3.35",
			}, {
				ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				Price:            "12.00",
			},
		},
		Total: "35.35",
	}

	result := calculatePoints(&rec, defaultPointRules...)
	expected := int64(28)
	assert.Equal(t, expected, result)
}

func TestCalculations2(t *testing.T) {
	rec := models.Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []models.Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
		Total: "9.00",
	}

	result := calculatePoints(&rec, defaultPointRules...)
	expected := int64(109)
	assert.Equal(t, expected, result)
}
