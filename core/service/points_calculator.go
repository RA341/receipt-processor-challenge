package service

import (
	"github.com/RA341/receipt-processor-challenge/models"
	u "github.com/RA341/receipt-processor-challenge/utils"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var (
	defaultPointRules = []calculationOpts{
		pointsForRetailerName(),
		pointsForRoundTotal(),
		pointsForTotalMultipleOfQuarter(),
		pointsPerTwoItems(),
		pointsForItemDescriptionLength(),
		pointsForOddPurchaseDay(),
		pointsForPurchaseTimeBetween2And4PM(),
	}
)

type calculationOpts func(receipt *models.Receipt) int64

func calculatePoints(receipt *models.Receipt, pointsCalculationsOpts ...calculationOpts) int64 {
	var points int64 = 0
	for _, opt := range pointsCalculationsOpts {
		points += opt(receipt)
	}

	return points
}

// Rule 1: One point for every alphanumeric character in the retailer name.
func pointsForRetailerName() calculationOpts {
	return func(receipt *models.Receipt) int64 {
		var points int64 = 0
		for _, r := range receipt.Retailer {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				points++
			}
		}
		return points
	}
}

// Rule 2: 50 points if the total is a round dollar amount with no cents.
func pointsForRoundTotal() calculationOpts {
	return func(receipt *models.Receipt) int64 {
		if strings.HasSuffix(receipt.Total, ".00") {
			_, err := strconv.ParseFloat(receipt.Total, 64)
			if err == nil { // if valid number
				return int64(50)
			}
		}
		return 0
	}
}

// Rule 3: 25 points if the total is a multiple of 0.25.
func pointsForTotalMultipleOfQuarter() calculationOpts {
	return func(receipt *models.Receipt) int64 {
		totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
		if err != nil {
			return 0
		}
		// Convert to cents for modulo check
		totalInCents := int64(math.Round(totalFloat*100 + 0.00001)) // Add epsilon for precision
		if totalInCents >= 0 && totalInCents%25 == 0 {
			return int64(25)
		}
		return 0
	}
}

// Rule 4: 5 points for every two items on the receipt.
func pointsPerTwoItems() calculationOpts {
	return func(receipt *models.Receipt) int64 {
		if receipt.Items == nil {
			return 0
		}

		numberOfPairs := int64(len(receipt.Items) / 2)
		points := numberOfPairs * 5.0
		return points
	}
}

// Rule 5: If the trimmed length of the item description is a multiple of 3,
// multiply the price by 0.2 and round up.
func pointsForItemDescriptionLength() calculationOpts {
	return func(receipt *models.Receipt) int64 {
		if receipt.Items == nil {
			return 0
		}
		var rulePoints int64 = 0
		for _, item := range receipt.Items {
			trimmedDesc := strings.TrimSpace(item.ShortDescription)
			descLen := len(trimmedDesc)

			if descLen > 0 && descLen%3 == 0 {
				priceFloat, err := strconv.ParseFloat(item.Price, 64)
				if err != nil {
					slog.Warn("Could not parse price for item. Skipping....",
						slog.String("price", item.Price),
						slog.String("shortDescription", item.ShortDescription),
						u.ErrLog(err),
					)
					continue
				}
				itemPoints := int64(math.Ceil(priceFloat * 0.2))
				rulePoints += itemPoints
			}
		}
		return rulePoints
	}
}

// Rule 6: 6 points if the day in the purchase date is odd.
func pointsForOddPurchaseDay() calculationOpts {
	return func(receipt *models.Receipt) int64 {
		layout := "2006-01-02" // YYYY-MM-DD
		purchaseDate, err := time.Parse(layout, receipt.PurchaseDate)
		if err != nil {
			slog.Warn("Could not parse purchase date",
				slog.String("purchaseDate", receipt.PurchaseDate),
				u.ErrLog(err),
			)
			return 0
		}
		day := purchaseDate.Day()
		if day%2 != 0 {
			return int64(6)
		}
		return 0
	}
}

// Rule 7: 10 points if the time of purchase is after 2:00pm (14:00) and before 4:00pm (16:00).
func pointsForPurchaseTimeBetween2And4PM() calculationOpts {
	return func(receipt *models.Receipt) int64 {
		layout := "15:04" // HH:MM (24-hour)
		purchaseTime, err := time.Parse(layout, receipt.PurchaseTime)
		if err != nil {
			slog.Warn("Could not parse purchase time",
				slog.String("purchaseTime", receipt.PurchaseTime),
				u.ErrLog(err),
			)
			return 0
		}
		// Define boundary times
		time1400, _ := time.Parse(layout, "14:00")
		time1600, _ := time.Parse(layout, "16:00")

		if purchaseTime.After(time1400) && purchaseTime.Before(time1600) {
			return int64(10)
		}
		return 0
	}
}
