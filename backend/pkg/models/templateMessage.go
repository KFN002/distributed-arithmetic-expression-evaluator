package models

type TemplateMessage struct {
	Expression *Expression
	Message    string
}

func (m *TemplateMessage) AddData(msg string, expression *Expression) {
	m.Expression = expression
	m.Message = msg
}

func (m *TemplateMessage) ChangeExpression(expression *Expression) {
	m.Expression = expression
}

func (m *TemplateMessage) ChangeMessage(msg string) {
	m.Message = msg
}

func CreateNewTemplateMessage() *TemplateMessage {
	return &TemplateMessage{}
}
