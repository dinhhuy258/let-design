package entity

type ScheduledEvent struct {
	Message string
}

type MessageBusRepository interface {
	PublishScheduledEvent(topic string, event ScheduledEvent) error
}
