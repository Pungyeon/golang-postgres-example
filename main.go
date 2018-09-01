package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
)

// Todo struct descripting todo objects
type Todo struct {
	UID         int
	Title       string
	Description string
	Username    string // guid
	Completed   bool
}

// TodoDB interface for interacting with backend database
type TodoDB interface {
	Get(id int) Todo
	GetAllUserTodos(user string) []Todo
	Insert(t Todo) (int, error)
	Put(t Todo) error
	Delete(id int) error
}

// PostgresTodoDB is a connection to a PostgresSQL todo database
// and holds methods for interacting with Todos in a PostgresSQL
type PostgresTodoDB struct { // @Implements TodoDB
	connection *sql.DB
}

// NewPostgresTodoDB returns a new PostgresTodoDB pointer
func NewPostgresTodoDB(config PostgresConfig) TodoDB {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName,
	)
	sqlConnection, err := sql.Open("postgres", connectionString)

	if err != nil {
		panic(err)
	}
	return &PostgresTodoDB{
		connection: sqlConnection,
	}
}

// PostgresConfig holds the configuration for a postgres connection
type PostgresConfig struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

// Get returns a todo from a postgres database
func (pq PostgresTodoDB) Get(id int) Todo {
	var todo Todo
	row := pq.connection.QueryRow(`SELECT * FROM todo WHERE uid=$1`, id)
	err := row.Scan(&todo.UID, &todo.Title, &todo.Description, &todo.Username, &todo.Completed)
	if err != nil {
		log.Println(err)
		return todo
	}
	return todo

}

// GetAllUserTodos will return all todos which are tied to a specified user
func (pq PostgresTodoDB) GetAllUserTodos(user string) []Todo {
	var todos []Todo
	rows, err := pq.connection.Query("SELECT * FROM todo WHERE username=$1", user)
	if err != nil {
		log.Println(err)
		return []Todo{}
	}
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.UID, &todo.Title, &todo.Description, &todo.Username, &todo.Completed)
		if err != nil {
			continue
		}
		todos = append(todos, todo)
	}
	return todos
}

// Insert a todo into a postgres database
func (pq PostgresTodoDB) Insert(t Todo) (int, error) {
	stmt, err := pq.connection.Prepare(`INSERT INTO todo(title, description, username, completed) VALUES($1,$2,$3,$4) returning uid;`)
	if err != nil {
		return 0, err
	}
	var uid int
	stmt.QueryRow(
		t.Title, t.Description, t.Username, t.Completed,
	).Scan(&uid)
	return uid, err
}

// Put edits a todo in a postgres database
func (pq PostgresTodoDB) Put(t Todo) error {
	stmt, err := pq.connection.Prepare("UPDATE todo SET title=$1 description=$2 completed=$3 where uid=$4")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(t.Title, t.Description, t.Completed, t.UID)
	return err
}

// Delete removes a todo from a postgres database
func (pq PostgresTodoDB) Delete(id int) error {
	_, err := pq.connection.Exec("DELETE FROM todo WHERE uid=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	db := NewPostgresTodoDB(PostgresConfig{"172.16.1.68", 5432, "postgres", "postgres", "test"})

	id, err := db.Insert(Todo{
		Title:       "Implement Postgres SQL",
		Description: "We need to ensure that everything is working via. the repository pattern",
		Username:    "Pungy",
		Completed:   false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted Todo was given UID: " + strconv.Itoa(id))
	/* 	todos := db.GetAllUserTodos("Pungy")
	   	for _, t := range todos {
	   		fmt.Println(t)
	   	} */
	db.Delete(id)
}
