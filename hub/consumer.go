package hub

import "FsnotifyPlus/event"

type Consumer struct {
	Id   string           //Идентификатор подписчика
	Send chan event.Event //Канал для отправки сообщений подписчику
}

func NewConsumer() *Consumer {
	return &Consumer{
		"no_id",
		make(chan event.Event),
	}
}

func NewConsumerWithId(id string) *Consumer {
	return &Consumer{
		Id:   id,
		Send: make(chan event.Event),
	}
}
