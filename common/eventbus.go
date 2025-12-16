package common

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventHandler struct {
	Id          string
	HandlerFunc func(data interface{})
}

type EventBus interface {
	Subscribe(event string, handler ...func(interface{})) (reference interface{}, err error)
	Unsubscribe(event string, handler_id string) error
	Publish(event string, data interface{})
}

type DefaultEventBus struct {
	handlers map[string][]EventHandler
}

func NewDefaultEventBus() EventBus {
	return &DefaultEventBus{
		handlers: make(map[string][]EventHandler, 0),
	}
}

func (eb *DefaultEventBus) RegisterHandler(event string, handler func(interface{})) (handler_id string, err error) {

	handler_id = primitive.NewObjectID().Hex()

	eb.handlers[event] = append(eb.handlers[event], EventHandler{
		Id:          handler_id,
		HandlerFunc: handler,
	})

	return handler_id, nil
}

func (eb *DefaultEventBus) UnregisterHandler(event string, handler_id string) error {
	for i, h := range eb.handlers[event] {
		if h.Id == handler_id {
			eb.handlers[event] = append(eb.handlers[event][:i], eb.handlers[event][i+1:]...)
			break
		}
	}

	return nil
}

func (eb *DefaultEventBus) Publish(event string, data interface{}) {
	for _, handler := range eb.handlers[event] {
		go handler.HandlerFunc(data)
	}
}
