package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataSource string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSource)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	migration := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        login TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS expressions (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expression TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    result REAL,
    error_message TEXT
	);`
	_, err = DB.Exec(migration)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

type Expression struct {
	ID           string
	UserID       int
	Expression   string
	Status       string
	Result       sql.NullFloat64
	ErrorMessage sql.NullString
}

func SaveExpression(expr *Expression) error {
	_, err := DB.Exec(
		"INSERT INTO expressions (id, user_id, expression, status) VALUES (?, ?, ?, ?)",
		expr.ID, expr.UserID, expr.Expression, expr.Status,
	)
	return err
}

func GetFirstPendingExpression() *Expression {
	row := DB.QueryRow("SELECT id, user_id, expression FROM expressions WHERE status = 'pending'")
	e := &Expression{}
	err := row.Scan(&e.ID, &e.UserID, &e.Expression)
	if err != nil {
		return nil
	}
	return e
}

func GetExpressionByID(id string) (*Expression, error) {
	e := &Expression{}
	err := DB.QueryRow("SELECT id, user_id, expression, status, result, error_message FROM expressions WHERE id = ?", id).Scan(
		&e.ID, &e.UserID, &e.Expression, &e.Status, &e.Result, &e.ErrorMessage,
	)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func UpdateExpressionStatus(id string, status string, result float64) error {
	_, err := DB.Exec("UPDATE expressions SET status = ?, result = ? WHERE id = ?", status, result, id)
	return err
}

func UpdateExpressionError(id string, errMsg string) error {
	_, err := DB.Exec(
		"UPDATE expressions SET status = 'error', error_message = ? WHERE id = ?", errMsg, id,
	)
	return err
}
