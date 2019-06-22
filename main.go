package main

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"os"
	"net/http"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Status string `json:"status"`
}

func getTodosHandler(c *gin.Context){
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	}
	defer db.Close()

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

	c.JSON(http.StatusOK, todos)
}

func getTodosByIDHandler(c *gin.Context){
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error open" : err.Error()})
		return
	}
	defer db.Close()

	stmt, err:= db.Prepare("SELECT id, title, status FROM todos WHERE id = $1")
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error prepare": err.Error()})
		return
	}

	id := c.Param("id")
	row := stmt.QueryRow(id)
	t := Todo{}
	err = row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error scan" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func postTodosHandler(c *gin.Context){
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error open" : err.Error()})
		return
	}
	defer db.Close()

	t := Todo{}
	if err = c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error bad request" : err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id, title, status")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error prepare" : err.Error()})
		return
	}

	row := stmt.QueryRow(t.Title, t.Status)
	if err = row.Scan(&t.ID, &t.Title, &t.Status) ; err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error scan" : err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func deleteTodosHandler(c *gin.Context){
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error open" : err.Error()})
		return
	}
	defer db.Close()

	id := c.Param("id")

	stmt, err := db.Prepare("DELETE FROM todos WHERE id = $1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error prepare" : err.Error()})
		return
	}

	if _, err = stmt.Exec(id) ; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error execute" : err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success")
}

func main() {
	r := gin.Default()

	api := r.Group("/api")
	api.GET("/todos", getTodosHandler)
	api.GET("/todos/:id", getTodosByIDHandler)
	api.POST("/todos", postTodosHandler)
	api.DELETE("/todos/:id", deleteTodosHandler)

	r.Run(":1234")
}