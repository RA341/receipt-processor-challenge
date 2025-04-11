package service

import (
	"github.com/RA341/receipt-processor-challenge/models"
	"testing"
)

type TestCase struct {
	receipt        models.Receipt
	expectedPoints int64
}

var (
	testMap = map[string]TestCase{
		"test 1": {
			receipt: models.Receipt{
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
			},
			expectedPoints: 28,
		},
		"test 2": {
			receipt: models.Receipt{
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
			},
			expectedPoints: 109,
		},
	}
)

func TestCalculatePointsAll(t *testing.T) {
	for name, test := range testMap {
		t.Run(name, func(t *testing.T) {
			runTest(t, test)
		})
	}
}

func TestCalculations1(t *testing.T) {
	testCase := testMap["test 1"]
	runTest(t, testCase)
}

func TestCalculations2(t *testing.T) {
	testCase := testMap["test 2"]
	runTest(t, testCase)
}

func runTest(t *testing.T, testCase TestCase) {
	result := calculatePoints(&testCase.receipt, defaultPointRules...)
	if result != testCase.expectedPoints {
		t.Fatalf("Expected %v but got %v", testCase.expectedPoints, result)
	}
}
