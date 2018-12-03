package main

import (
	"fmt"
	"sync"

	"../services"
)

// Database represents some storage able to store records
type Database interface {
	AddRecord(r services.Record)
}

// NewMemoryDatabase creates a simple database in memory
func NewMemoryDatabase() Database {
	return &memoryDatabase{
		mu:      sync.RWMutex{},
		records: map[string]services.Record{},
	}
}

type memoryDatabase struct {
	mu      sync.RWMutex
	records map[string]services.Record
}

func (db *memoryDatabase) AddRecord(r services.Record) {
	fmt.Printf("Adding record %#v to the database\n", r)
	db.mu.Lock()
	db.records[r.ID] = r
	db.mu.Unlock()
}
