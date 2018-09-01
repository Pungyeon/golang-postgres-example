package main

import (
	"fmt"
	"strconv"

	"github.com/Pungyeon/pq-go/db"
)

func main() {
	database := db.NewPostgresTodoDB(PostgresConfig{"172.16.1.68", 5432, "postgres", "postgres", "test"})

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
