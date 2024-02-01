package domain

type AggregateRoot interface {
	AddLocalEvent(eventData any)
	GetLocalEvents() []any
	ClearLocalEvents()

	AddDistributedEvent(eventData any)
	GetDistributedEvents() []any
	ClearDistributedEvents()
}

type BasicAggregateRoot struct {
	LocalEvents       []any
	DistributedEvents []any
}

func (entity *BasicAggregateRoot) AddLocalEvent(eventData any) {
	if entity.LocalEvents == nil {
		entity.LocalEvents = make([]any, 0, 1)
	}
	entity.LocalEvents = append(entity.LocalEvents, eventData)
}

func (entity *BasicAggregateRoot) GetLocalEvents() []any {
	if entity.LocalEvents == nil {
		return make([]any, 0)
	}
	return entity.LocalEvents
}

func (entity *BasicAggregateRoot) ClearLocalEvents() {
	if entity.LocalEvents != nil {
		entity.LocalEvents = entity.LocalEvents[:0]
	}
}

func (entity *BasicAggregateRoot) AddDistributedEvent(eventData any) {
	if entity.DistributedEvents == nil {
		entity.DistributedEvents = make([]any, 0, 1)
	}
	entity.DistributedEvents = append(entity.DistributedEvents, eventData)
}

func (entity *BasicAggregateRoot) GetDistributedEvents() []any {
	if entity.DistributedEvents == nil {
		return make([]any, 0)
	}
	return entity.DistributedEvents
}

func (entity *BasicAggregateRoot) ClearDistributedEvents() {
	if entity.DistributedEvents != nil {
		entity.DistributedEvents = entity.DistributedEvents[:0]
	}
}
