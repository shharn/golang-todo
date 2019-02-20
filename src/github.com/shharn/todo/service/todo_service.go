package service

import (
	"time"

	"github.com/pkg/errors"
	"github.com/shharn/todo/model"
	"github.com/shharn/todo/db"
)

const pageSize = 3

type TodoService interface {
	CreateTodo(model.Todo) error
	GetTodosByPageNo(int) ([]model.Todo, int)
	UpdateTodo(model.Todo) error
}

type SimpleTodoService struct {
	dbs db.Database
}

func (s SimpleTodoService) CreateTodo(todo model.Todo) error {
	now := time.Now().Format(time.RFC3339)
	todo.CreatedAt = now
	todo.ModifiedAt = now
	todo.IsComplete = false
	err := s.dbs.Create(todo)
	return err
}

func (s SimpleTodoService) GetTodosByPageNo(pageNo int) ([]model.Todo, int) {
	todos, _ := s.dbs.GetPagedList(pageNo, pageSize).([]model.Todo)
	totalCount := s.dbs.GetTotalCount()
	return todos, totalCount
}

func (s SimpleTodoService) UpdateTodo(todo model.Todo) error {
	oldTodo, _ := s.dbs.Get(todo.Id).(model.Todo)
	if !oldTodo.IsComplete && todo.IsComplete && !s.checkIfParentsAreComplete(todo){
		return errors.New("Cannot update todo")
	}
	err := s.dbs.Update(todo)
	return err
}

func (s SimpleTodoService) checkIfParentsAreComplete(todo model.Todo) bool {
	if todo.ParentIds == nil || len(todo.ParentIds) < 1 {
		return todo.IsComplete
	}

	result := true
	for _, pid := range todo.ParentIds {
		parentTodo, _ := s.dbs.Get(pid).(model.Todo)
		result = result && s.checkIfParentsAreComplete(parentTodo)
		if !result {
			return false
		}
	}
	return result
}

func NewSimpleTodoService() *SimpleTodoService {
	return &SimpleTodoService{
		dbs: &db.Instance,
	}
}