package types

type Expression struct {
	ID           string
	Expr         string
	Status       string
	Result       *float64
	ErrorMessage string
}
