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
	dataMap *sync.Map
}

func NewDB() (*FranklyWeHaveNoIdeaWhereYourDataIsDB, error) {
	return &FranklyWeHaveNoIdeaWhereYourDataIsDB{dataMap: &sync.Map{}}, nil
}

func (f *FranklyWeHaveNoIdeaWhereYourDataIsDB) CreatePoint(totalPoints int64) (transactionId string, err error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	transactionId = newUUID.String()
	f.dataMap.Store(transactionId, totalPoints)

	return transactionId, err
}

func (f *FranklyWeHaveNoIdeaWhereYourDataIsDB) GetPointById(transactionId string) (totalPoints int64, err error) {
	tmpPoint, ok := f.dataMap.Load(transactionId)
	if !ok {
		return 0, fmt.Errorf("unable to find points for: %s", transactionId)
	}

	return tmpPoint.(int64), nil
}
