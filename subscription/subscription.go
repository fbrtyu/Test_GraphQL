package subscription

import (
	"context"
	"fmt"
	"ozon-test/graph/model"
	"time"
)

// Есть идеи как реализовать получение уведомлений, только от определенных постов, но по времени не успеваю реализовать

// Структура для ответа в уведомлении
type NewComment struct {
	Text string `json:"Text"`
}

// Канал, в который пишутся все новые комментарии
var Publicch = make(chan NewComment)

func Subonpost(ctx context.Context) <-chan *model.Comment {
	ch := make(chan *model.Comment)

	// Бесконечный цикл каждую секунду присылает новые данные из канала Publicch в канал ch клиенту
	go func() {
		defer close(ch)
		for {
			time.Sleep(1 * time.Second)
			newcomment := <-Publicch
			comment := &model.Comment{
				Text: newcomment.Text,
			}
			select {
			case <-ctx.Done():
				fmt.Println("Subscription Closed")
				return
			case ch <- comment:
			}
		}
	}()
	return ch
}
