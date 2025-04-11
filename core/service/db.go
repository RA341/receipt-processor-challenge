package service

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
)

type Database interface {
	CreatePoint(totalPoints int64) (transactionId string, err error)
	GetPointById(transactionId string) (totalPoints int64, err error)
}

type FranklyWeHaveNoIdeaWhereYourDataIsDB struct {
	pointsTable *sync.Map
}

func NewDB() (*FranklyWeHaveNoIdeaWhereYourDataIsDB, error) {
	return &FranklyWeHaveNoIdeaWhereYourDataIsDB{pointsTable: &sync.Map{}}, nil
}

func (f *FranklyWeHaveNoIdeaWhereYourDataIsDB) CreatePoint(totalPoints int64) (transactionId string, err error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	transactionId = newUUID.String()
	f.pointsTable.Store(transactionId, totalPoints)

	return transactionId, err
}

func (f *FranklyWeHaveNoIdeaWhereYourDataIsDB) GetPointById(transactionId string) (totalPoints int64, err error) {
	tmpPoint, ok := f.pointsTable.Load(transactionId)
	if !ok {
		return 0, fmt.Errorf("unable to find points for: %s", transactionId)
	}

	return tmpPoint.(int64), nil
}
