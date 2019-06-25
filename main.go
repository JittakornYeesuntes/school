package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"school/todo"
)

func main() {
	r := gin.Default()

	api := r.Group("/api")
	api.GET("/todos", todo.GetTodosHandler)
	api.GET("/todos/:id", todo.GetTodosByIDHandler)
	api.POST("/todos", todo.PostTodosHandler)
	api.DELETE("/todos/:id", todo.DeleteTodosHandler)

	r.Run(":1234")
}
