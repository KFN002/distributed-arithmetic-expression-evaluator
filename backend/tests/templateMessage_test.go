package tests

import (
	"github.com/KFN002/distributed-arithmetic-expression-evaluator.git/backend/pkg/models"
	"testing"
)

func TestTemplateMessage_AddData(t *testing.T) {
	expr := &models.Expression{ID: 1, Expression: "2+3", Status: "pending", Result: nil, CreatedAt: "01-01-2022 10:00:00", FinishedAt: nil, UserID: 1}
	msg := "Test message"

	m := &models.TemplateMessage{}
	m.AddData(msg, expr)

	if m.Expression != expr {
		t.Errorf("AddData() failed to set Expression correctly. Got: %v, Expected: %v", m.Expression, expr)
	}

	if m.Message != msg {
		t.Errorf("AddData() failed to set Message correctly. Got: %s, Expected: %s", m.Message, msg)
	}
}

func TestTemplateMessage_ChangeExpression(t *testing.T) {
	expr1 := &models.Expression{ID: 1, Expression: "2+3", Status: "pending", Result: nil, CreatedAt: "01-01-2022 10:00:00", FinishedAt: nil, UserID: 1}
	expr2 := &models.Expression{ID: 2, Expression: "3*4", Status: "pending", Result: nil, CreatedAt: "01-01-2022 10:00:00", FinishedAt: nil, UserID: 1}

	m := &models.TemplateMessage{Expression: expr1}
	m.ChangeExpression(expr2)

	if m.Expression != expr2 {
		t.Errorf("ChangeExpression() failed to change Expression. Got: %v, Expected: %v", m.Expression, expr2)
	}
}

func TestTemplateMessage_ChangeMessage(t *testing.T) {
	msg1 := "Test message 1"
	msg2 := "Test message 2"

	m := &models.TemplateMessage{Message: msg1}
	m.ChangeMessage(msg2)

	if m.Message != msg2 {
		t.Errorf("ChangeMessage() failed to change Message. Got: %s, Expected: %s", m.Message, msg2)
	}
}

func TestCreateNewTemplateMessage(t *testing.T) {
	m := models.CreateNewTemplateMessage()

	if m.Expression != nil {
		t.Errorf("CreateNewTemplateMessage() failed to initialize Expression. Got: %v, Expected: nil", m.Expression)
	}

	if m.Message != "" {
		t.Errorf("CreateNewTemplateMessage() failed to initialize Message. Got: %s, Expected: ''", m.Message)
	}
}
