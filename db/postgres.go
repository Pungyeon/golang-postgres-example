package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Pungyeon/pq-go/todo"
	_ "github.com/lib/pq"
)

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
func (pq PostgresTodoDB) Get(id int) todo.Todo {
	var t todo.Todo
	row := pq.connection.QueryRow(`SELECT * FROM todo WHERE uid=$1`, id)
	err := row.Scan(&t.UID, &t.Title, &t.Description, &t.Username, &t.Completed)
	if err != nil {
		log.Println(err)
		return t
	}
	return t

}

// GetAllUserTodos will return all todos which are tied to a specified user
func (pq PostgresTodoDB) GetAllUserTodos(user string) []todo.Todo {
	var todos []todo.Todo
	rows, err := pq.connection.Query("SELECT * FROM todo WHERE username=$1", user)
	if err != nil {
		log.Println(err)
		return []todo.Todo{}
	}
	for rows.Next() {
		var t todo.Todo
		err = rows.Scan(&t.UID, &t.Title, &t.Description, &t.Username, &t.Completed)
		if err != nil {
			continue
		}
		todos = append(todos, t)
	}
	return todos
}

// Insert a todo into a postgres database
func (pq PostgresTodoDB) Insert(t todo.Todo) (int, error) {
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
func (pq PostgresTodoDB) Put(t todo.Todo) error {
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
