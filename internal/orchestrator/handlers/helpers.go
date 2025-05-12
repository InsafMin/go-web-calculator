package handlers

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/InsafMin/go-web-calculator/pkg/calculator"
	"github.com/InsafMin/go-web-calculator/pkg/errors"
	"github.com/InsafMin/go-web-calculator/pkg/types"
)

func ParseExpression(expr string, exprID string) ([]*types.Task, error) {
	tokens, err := calculator.Tokenize(expr)
	if err != nil {
		return nil, err
	}
	rpn, err := calculator.ToRPN(tokens)
	if err != nil {
		return nil, err
	}
	bracketLevels := getBracketLevels(tokens, rpn)

	var tasks []*types.Task
	var stack []string
	taskMap := make(map[string]string)
	taskCounter := 1

	for _, token := range rpn {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, token)
		} else if calculator.IsOperator(rune(token[0])) {
			if len(stack) < 2 {
				return nil, errors.ErrInvalidExpression
			}
			arg2 := stack[len(stack)-1]
			arg1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			taskID := exprID + "-" + strconv.Itoa(taskCounter)
			if val, ok := taskMap[arg1]; ok {
				arg1 = val
			}
			if val, ok := taskMap[arg2]; ok {
				arg2 = val
			}

			priority := calculator.Priority(token) + bracketLevels[token]
			task := &types.Task{
				ID:            taskID,
				Arg1:          parseNumber(arg1),
				Arg2:          parseNumber(arg2),
				Operation:     token,
				ExpressionID:  exprID,
				Priority:      priority,
				OperationTime: getOperationTime(token),
				Done:          make(chan bool),
			}
			tasks = append(tasks, task)
			resultKey := "task-" + taskID
			taskMap[resultKey] = taskID
			stack = append(stack, resultKey)
			taskCounter++
		}
	}

	return tasks, nil
}

func getOperationTime(operation string) time.Duration {
	var timeMs int
	var err error
	switch operation {
	case "+":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	case "-":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	case "*":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	case "/":
		timeMs, err = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	default:
		return 0
	}
	if err != nil {
		switch operation {
		case "+":
			return 100 * time.Millisecond
		case "-":
			return 100 * time.Millisecond
		case "*":
			return 200 * time.Millisecond
		case "/":
			return 200 * time.Millisecond
		default:
			return 0
		}
	}
	return time.Duration(timeMs) * time.Millisecond
}

func getBracketLevels(tokens []string, rpn []string) map[string]int {
	bracketLevels := make(map[string]int)
	currentLevel := 0
	tokenLevels := make(map[string]int)
	for _, token := range tokens {
		if token == "(" {
			currentLevel += 2
		} else if token == ")" {
			currentLevel -= 2
		} else if calculator.IsOperator(rune(token[0])) {
			tokenLevels[token] = currentLevel
		}
	}
	for _, token := range rpn {
		if calculator.IsOperator(rune(token[0])) {
			bracketLevels[token] = tokenLevels[token]
		}
	}
	return bracketLevels
}

func parseNumber(s string) float64 {
	if strings.HasPrefix(s, "task-") {
		return 0
	}
	num, _ := strconv.ParseFloat(s, 64)
	return num
}
