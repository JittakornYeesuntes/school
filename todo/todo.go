package todo

import (
	"school/database"

	"github.com/gin-gonic/gin"
	"net/http"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func GetTodosHandler(c *gin.Context) {
	rows, err := database.SelectAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error database": err.Error()})
		return
	}

	todos := []Todo{}
	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error scan": err.Error()})
			return
		}

		todos = append(todos, t)
	}

	c.JSON(http.StatusOK, todos)
}

func GetTodosByIDHandler(c *gin.Context) {
	id := c.Param("id")
	row, err := database.SelectByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error database": err.Error()})
		return
	}

	t := Todo{}
	err = row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error scan": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t)
}

func PostTodosHandler(c *gin.Context) {
	t := Todo{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error bad request": err.Error()})
		return
	}

	row, err := database.InsertTodos(t.Title, t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error database": err.Error()})
		return
	}

	err = row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error scan": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func DeleteTodosHandler(c *gin.Context) {
	id := c.Param("id")

	err := database.DeleteByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error database": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "success")
}
