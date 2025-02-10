package schedule

import (
	"fmt"

	"github.com/fgouvea/weather/weather-service/queue"
)

type Publisher struct {
	Saver     ScheduleSaver
	Publisher *queue.Publisher
}

func NewPublisher(saver ScheduleSaver, publisher *queue.Publisher) *Publisher {
	return &Publisher{
		Saver:     saver,
		Publisher: publisher,
	}
}

func (p *Publisher) Publish(schedule Schedule) error {
	schedule.Status = StatusProcessing

	err := p.Saver.Save(schedule)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSave, err)
	}

	return queue.Publish(p.Publisher, schedule)
}
