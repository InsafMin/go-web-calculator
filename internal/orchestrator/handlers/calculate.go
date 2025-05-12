package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/InsafMin/go-web-calculator/internal/db"
	"github.com/InsafMin/go-web-calculator/pkg/types"
)

func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	userID := types.GetUserID(r.Context())
	if userID == -1 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())

	dbExpr := db.Expression{
		ID:         id,
		UserID:     userID,
		Expression: req.Expression,
		Status:     "pending",
	}

	if err := db.SaveExpression(&dbExpr); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
