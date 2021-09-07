package eventbus

import "sync"

type EventData struct {
	Data    interface{}
	Subject string
}

type EventChannel chan EventData

type EventChannelEntries []EventChannel

// NewEventBus creates new EventBus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: map[string]EventChannelEntries{},
	}
}

type EventBus struct {
	subscribers map[string]EventChannelEntries
	rm          sync.RWMutex
}

// Subscribe subscribes a channel to a subject
func (eb *EventBus) Subscribe(subject string, ch EventChannel) {
	eb.rm.Lock()
	defer eb.rm.Unlock()
	if prev, found := eb.subscribers[subject]; found {
		eb.subscribers[subject] = append(prev, ch)
	} else {
		eb.subscribers[subject] = append([]EventChannel{}, ch)
	}

}

// Publish sends the data object to all subscribers that are listening to a subject
func (eb *EventBus) Publish(subject string, data interface{}) {
	eb.rm.RLock()
	defer eb.rm.RUnlock()
	if chans, found := eb.subscribers[subject]; found {
		for _, ch := range chans {
			data := EventData{Data: data, Subject: subject}
			ch <- data
		}
	}
}
