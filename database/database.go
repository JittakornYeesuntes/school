package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

func Connect() (*sql.DB, error) {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return conn, err
	}
	return conn, err
}

func SelectAll() (*sql.Rows, error) {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	stmt, err := conn.Prepare("SELECT id, title, status FROM todos")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return rows, err
	}

	return rows, err
}

func SelectByID(id string) (*sql.Row, error) {
	conn, err := Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stmt, err := conn.Prepare("SELECT id, title, status FROM todos WHERE id = $1")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(id)
	return row, nil
}

func InsertTodos(title string, status string) (*sql.Row, error) {
	conn, err := Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stmt, err := conn.Prepare("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id, title, status")
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(title, status)
	return row, nil
}

func DeleteByID(id string) error {
	conn, err := Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	stmt, err := conn.Prepare("DELETE FROM todos WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
