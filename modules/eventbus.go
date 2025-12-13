package modules

type EventBus interface {
	RegisterHandler(event string, handler func(interface{}))
	UnregisterHandler(event string, handler func(interface{}))
	Publish(event string, data interface{})
}

type DefaultEventBus struct {
	handlers map[string][]func(interface{})
}

func NewDefaultEventBus() EventBus {
	return &DefaultEventBus{
		handlers: make(map[string][]func(interface{})),
	}
}

func (eb *DefaultEventBus) RegisterHandler(event string, handler func(interface{})) {
	eb.handlers[event] = append(eb.handlers[event], handler)
}

func (eb *DefaultEventBus) UnregisterHandler(event string, handler func(interface{})) {
	for i, h := range eb.handlers[event] {
		if h == handler {
			eb.handlers[event] = append(eb.handlers[event][:i], eb.handlers[event][i+1:]...)
			break
		}
	}
}

func (eb *DefaultEventBus) Publish(event string, data interface{}) {
	for _, handler := range eb.handlers[event] {
		handler(data)
	}
}
