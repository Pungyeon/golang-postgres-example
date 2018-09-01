package db

import "github.com/Pungyeon/pq-go/todo"

// TodoDB interface for interacting with backend database
type TodoDB interface {
	Get(id int) todo.Todo
	GetAllUserTodos(user string) []todo.Todo
	Insert(t todo.Todo) (int, error)
	Put(t todo.Todo) error
	Delete(id int) error
}
