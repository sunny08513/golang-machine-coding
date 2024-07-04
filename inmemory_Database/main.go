package main

import (
	"fmt"
	"sync"
)

type Record map[string]interface{}
type Table map[string]Record

type InMemoryDB struct {
	tables map[string]Table
	mu     sync.RWMutex // Use a mutex for thread-safe operations
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		tables: make(map[string]Table),
	}
}

func (db *InMemoryDB) CreateTable(name string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.tables[name]; !exists {
		db.tables[name] = make(Table)
	}
}

func (db *InMemoryDB) Insert(tableName string, id string, record Record) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	table, exists := db.tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	table[id] = record
	return nil
}

func (db *InMemoryDB) Retrieve(tableName string, id string) (Record, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	table, exists := db.tables[tableName]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", tableName)
	}

	record, exists := table[id]
	if !exists {
		return nil, fmt.Errorf("record with id %s does not exist", id)
	}

	return record, nil
}

func (db *InMemoryDB) Update(tableName string, id string, updatedRecord Record) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	table, exists := db.tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	if _, exists := table[id]; !exists {
		return fmt.Errorf("record with id %s does not exist", id)
	}

	table[id] = updatedRecord
	return nil
}

func (db *InMemoryDB) Delete(tableName string, id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	table, exists := db.tables[tableName]
	if !exists {
		return fmt.Errorf("table %s does not exist", tableName)
	}

	if _, exists := table[id]; !exists {
		return fmt.Errorf("record with id %s does not exist", id)
	}

	delete(table, id)
	return nil
}

func main() {
	db := NewInMemoryDB()

	db.CreateTable("users")

	user1 := Record{"name": "Alice", "age": 30}
	user2 := Record{"name": "Bob", "age": 25}

	db.Insert("users", "1", user1)
	db.Insert("users", "2", user2)

	record, err := db.Retrieve("users", "1")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(record)
	}

	updatedUser1 := Record{"name": "Alice", "age": 31}
	db.Update("users", "1", updatedUser1)

	db.Delete("users", "2")
}
