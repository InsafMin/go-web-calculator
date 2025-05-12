package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	db "github.com/InsafMin/go-web-calculator/internal/db"
	handler "github.com/InsafMin/go-web-calculator/internal/orchestrator/handlers"
)

func setupTestServer() *httptest.Server {
	// ✅ Инициализируем БД в памяти
	db.InitDB(":memory:")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", handler.HandleRegister)
	mux.HandleFunc("/api/v1/login", handler.HandleLogin)
	mux.HandleFunc("/api/v1/calculate", handler.HandleCalculate)
	mux.HandleFunc("/api/v1/expressions", handler.HandleGetExpressions)
	mux.HandleFunc("/api/v1/expressions/", handler.HandleGetExpression)

	return httptest.NewServer(mux)
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	ID string `json:"id"`
}

func getTestToken(t *testing.T) string {
	ts := setupTestServer()
	defer ts.Close()

	registerURL := ts.URL + "/api/v1/register"
	loginURL := ts.URL + "/api/v1/login"

	registerReq := RegisterRequest{
		Login:    "testuser",
		Password: "testpass",
	}
	jsonData, _ := json.Marshal(registerReq)

	resp, err := http.Post(registerURL, "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	jsonData, _ = json.Marshal(registerReq)
	resp, err = http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NoError(t, err)
	assert.NotEmpty(t, loginResp.Token)

	return loginResp.Token
}

func TestHandleCalculate(t *testing.T) {
	ts := setupTestServerWithDummyAuth()
	defer ts.Close()

	calculateURL := ts.URL + "/api/v1/calculate"
	expr := CalculateRequest{
		Expression: "2+2*2",
	}
	jsonExpr, _ := json.Marshal(expr)

	req, err := http.NewRequest("POST", calculateURL, bytes.NewBuffer(jsonExpr))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var calcResp CalculateResponse
	err = json.NewDecoder(resp.Body).Decode(&calcResp)
	require.NoError(t, err)
	assert.NotEmpty(t, calcResp.ID)
}

func setupTestServerWithDummyAuth() *httptest.Server {
	db.InitDB(":memory:")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", handler.HandleRegister)
	mux.HandleFunc("/api/v1/login", handler.HandleLogin)
	mux.HandleFunc("/api/v1/calculate", DummyAuthMiddleware(handler.HandleCalculate))

	return httptest.NewServer(mux)
}

func DummyAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user_id", 1)))
	}
}
