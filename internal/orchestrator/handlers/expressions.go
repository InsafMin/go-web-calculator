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

	rows, err := db.DB.Query("SELECT id, expression, status, result, error_message FROM expressions WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var exprs []db.Expression
	for rows.Next() {
		var e db.Expression
		err := rows.Scan(&e.ID, &e.Expression, &e.Status, &e.Result, &e.ErrorMessage)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		exprs = append(exprs, e)
	}

	response := make([]types.Expression, len(exprs))
	for i, e := range exprs {
		var resultPtr *float64
		if e.Result.Valid {
			result := e.Result.Float64
			resultPtr = &result
		} else {
			resultPtr = nil
		}

		response[i] = types.Expression{
			ID:           e.ID,
			Expr:         e.Expression,
			Status:       e.Status,
			Result:       resultPtr,
			ErrorMessage: e.ErrorMessage.String,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]types.Expression{"expressions": response})
}

func HandleGetExpression(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/expressions/"):]
	expr, err := db.GetExpressionByID(id)
	if err != nil {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expression": expr.Expression,
		"status":     expr.Status,
		"result":     expr.Result.Float64,
		"error":      expr.ErrorMessage.String,
	})
}
