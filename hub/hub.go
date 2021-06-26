package hub

import (
	"FsnotifyPlus/event"
)

//Бестиповая шина для передачи сообщений
type Hub struct {
	Consumers  map[*Consumer]bool //Подписчики, который имеют канал для отправки сообщений
	Broadcast  chan event.Event   //Канал для передачи сообщений в шину. Полученные сообщения добавляются в буфер и рассылаются подписчикам
	Register   chan *Consumer     //Регистрация нового подписчика
	Unregister chan *Consumer     //Отписка от события для конкретного подписчика
}

func NewHub() *Hub {
	return &Hub{
		Consumers:  make(map[*Consumer]bool),
		Broadcast:  make(chan event.Event),
		Register:   make(chan *Consumer),
		Unregister: make(chan *Consumer),
	}
}

//Метод, который запускает работу шины. Обрабатывает все основные операции
//TODO следует запустить в отделльной горутине
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Consumers[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Consumers[client]; ok {
				delete(h.Consumers, client)
				close(h.Unregister)
			}
		case sendingEvent := <-h.Broadcast:
			for consumer := range h.Consumers {
				consumer.Send <- sendingEvent
			}
		}
	}
}
