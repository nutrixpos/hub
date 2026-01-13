package common

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventChannel struct {
	Id      string
	Channel chan interface{}
}

type EventManager interface {
	Subscribe(event string) (eventChannel EventChannel, err error)
	Publish(event string, data interface{})
}

type DefaultEventManager struct {
	channels map[string][]EventChannel
}

func NewDefaultEventManager() EventManager {
	return &DefaultEventManager{
		channels: make(map[string][]EventChannel, 0),
	}
}

func (eb *DefaultEventManager) Subscribe(event string) (eventChannel EventChannel, err error) {

	channel_id := primitive.NewObjectID().Hex()

	eventChannel = EventChannel{
		Id:      channel_id,
		Channel: make(chan interface{}),
	}

	eb.channels[event] = append(eb.channels[event], eventChannel)

	return eventChannel, nil
}

func (eb *DefaultEventManager) Publish(event string, data interface{}) {
	for _, even_channel := range eb.channels[event] {
		even_channel.Channel <- data
	}
}
