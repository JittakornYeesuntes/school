package main

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"os"
	"net/http"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID int
	Title string
	Status string
}

func getTodosHandler(c *gin.Context){
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	stmt, err := db.Prepare("SELECT id, title, status FROM todos")
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	rows, err:= stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}

	todos := []Todo{}
	for rows.Next(){
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}

		todos = append(todos, t)
	}

	c.JSON(200, todos)
}

func main() {
	r := gin.Default()

	r.GET("/api/todos", getTodosHandler)

	r.Run(":1234")
}