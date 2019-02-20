package db

import (
	"log"
	"sync"

	"github.com/shharn/todo/model"
)

type Database interface {
	Create(interface{}) error
	Update(interface{}) error
	GetAll() interface{}
	GetPagedList(int, int) interface{}
	Get(int) interface{}
	GetTotalCount() int
}

type TodoMemoryDatabase struct {
	nextId int
	todos map[int]model.Todo
	rwLock *sync.RWMutex
}

func (db *TodoMemoryDatabase) Create(todo interface{}) error {
	db.rwLock.Lock()
	defer db.rwLock.Unlock()
	cTodo, _ := todo.(model.Todo)
	cTodo.Id = db.nextId
	db.nextId++
	db.todos[cTodo.Id] = cTodo
	log.Printf("Created todo - %v", cTodo)
	return nil
}

func (db *TodoMemoryDatabase) Update(todo interface{}) error {
	db.rwLock.Lock()
	defer db.rwLock.Unlock()
	cTodo, _ := todo.(model.Todo)
	db.todos[cTodo.Id] = cTodo
	return nil
}

func (db *TodoMemoryDatabase) GetAll() interface{} {
	db.rwLock.RLock()
	defer db.rwLock.RUnlock()
	result := []model.Todo{}
	for _, todo := range db.todos {
		result = append(result, todo)
	}
	return result
}

func (db *TodoMemoryDatabase) GetPagedList(pageNo, pageSize int) interface{} {
	db.rwLock.RLock()
	defer db.rwLock.RUnlock()
	offset := getOffset(pageNo, pageSize)
	upperLimit := getUpperLimit(offset, pageSize, db.GetTotalCount())
	result := []model.Todo{}
	for i := offset; i < upperLimit; i++ {
		result = append(result, db.todos[i])
	}
	return result
}

func (db *TodoMemoryDatabase) Get(id int) interface{} {
	db.rwLock.RLock()
	defer db.rwLock.RUnlock()
	if todo, ok := db.todos[id]; ok {
		return todo
	}
	return model.Todo{}
}

func (db *TodoMemoryDatabase) GetTotalCount() int {
	db.rwLock.RLock()
	defer db.rwLock.RUnlock()
	return db.nextId - 1
}

func getOffset(pageNo, pageSize int) int {
	return (pageNo - 1) * pageSize + 1
}

func getUpperLimit(offset, pageSize, total int) int {
	maybeResult := offset + pageSize
	if maybeResult <= total {
		return maybeResult 
	}
	return total + 1
}

var Instance = TodoMemoryDatabase{
	nextId: 1,
	todos: map[int]model.Todo{},
	rwLock: &sync.RWMutex{},
}
