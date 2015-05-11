package db

import (
	"github.com/jinzhu/gorm"
	"github.com/stellar/go-horizon/test"
	"log"
)

func OpenTestDatabase() gorm.DB {

	result, err := gorm.Open("postgres", test.DatabaseUrl())

	if err != nil {
		log.Panic(err)
	}
	result.LogMode(true)
	return result
}

func OpenStellarCoreTestDatabase() gorm.DB {

	result, err := gorm.Open("postgres", test.StellarCoreDatabaseUrl())

	if err != nil {
		log.Panic(err)
	}
	result.LogMode(true)
	return result
}

type mockDumpQuery struct{}
type mockStreamedQuery struct{}

func (q mockDumpQuery) Get() ([]interface{}, error) {
	return []interface{}{
		"hello",
		"world",
		"from",
		"go",
	}, nil
}

func (q mockDumpQuery) IsComplete(alreadyDelivered int) bool {
	return alreadyDelivered >= 4
}

type mockQuery struct {
	resultCount int
}

type mockResult struct {
	index int
}

func (q mockQuery) Get() ([]interface{}, error) {
	results := make([]interface{}, q.resultCount)

	for i := 0; i < q.resultCount; i++ {
		results[i] = mockResult{i}
	}

	return results, nil
}

func (q mockQuery) IsComplete(alreadyDelivered int) bool {
	return alreadyDelivered >= q.resultCount
}

type BrokenQuery struct {
	Err error
}

func (q BrokenQuery) Get() ([]interface{}, error) {
	return nil, q.Err
}

func (q BrokenQuery) IsComplete(alreadyDelivered int) bool {
	return alreadyDelivered > 0
}

func MustFirst(q Query) interface{} {
	result, err := First(q)

	if err != nil {
		panic(err)
	}

	return result
}

func MustResults(q Query) []interface{} {
	result, err := Results(q)

	if err != nil {
		panic(err)
	}

	return result
}
