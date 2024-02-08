package demonAgent

import (
	"distributed-arithmetic-expression-evaluator/backend/models"
	"distributed-arithmetic-expression-evaluator/backend/queueMaster"
	"fmt"
	"time"
)

func QueueHandler() {
	for {
		gotExpr, expression := queueMaster.ExpressionsQueue.Dequeue()
		if gotExpr {
			answerCh := make(chan bool)
			go work(expression, answerCh)
			<-answerCh
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func work(expression models.Expression, answerCh chan bool) {
	fmt.Println(expression)

	// Здесь вы выполняете фактическую работу по обработке выражения
	// После завершения обработки отправьте сигнал в канал, чтобы сообщить, что работа завершена

	time.Sleep(10 * time.Second)
	answerCh <- true
}
