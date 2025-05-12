package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/InsafMin/go-web-calculator/internal/db"
	"github.com/InsafMin/go-web-calculator/pkg/types"
)

func HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	userID := types.GetUserID(r.Context())
	if userID == -1 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query("SELECT id, expression, status, result FROM expressions WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var exprs []db.Expression
	for rows.Next() {
		var e db.Expression
		err := rows.Scan(&e.ID, &e.Expression, &e.Status, &e.Result)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		exprs = append(exprs, e)
	}

	response := make([]types.Expression, len(exprs))
	for i, e := range exprs {
		response[i] = types.Expression{
			ID:     e.ID,
			Expr:   e.Expression,
			Status: e.Status,
			Result: e.Result,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]types.Expression{"expressions": response})
}

func HandleGetExpression(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	userID := types.GetUserID(r.Context())

	expr, err := db.GetExpressionByID(id)
	if err != nil {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	if expr.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]types.Expression{
		"expression": types.Expression{
			ID:     expr.ID,
			Expr:   expr.Expression,
			Status: expr.Status,
			Result: expr.Result,
		},
	})
}
