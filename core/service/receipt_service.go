package service

import "github.com/RA341/receipt-processor-challenge/models"

var ()

type ReceiptService struct {
	db Database
}

func NewReceiptService(db Database) *ReceiptService {
	return &ReceiptService{db: db}
}

func (s *ReceiptService) GetPointsById(transactionId string) (totalPoints int64, err error) {
	return s.db.GetPointById(transactionId)
}

func (s *ReceiptService) NewReceipt(receipt models.Receipt) (transactionId string, err error) {
	finalPoints := calculatePoints(
		&receipt,
		defaultPointRules...,
	)

	pointId, err := s.db.CreatePoint(finalPoints)
	if err != nil {
		return "", err
	}

	return pointId, nil
}
