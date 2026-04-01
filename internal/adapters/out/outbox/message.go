package outbox

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Name           string     `gorm:"type:varchar(255);index"`
	Payload        []byte     `gorm:"type:bytea"`
	OccurredAtUtc  time.Time  `gorm:"column:occurred_at_utc"`
	ProcessedAtUtc *time.Time `gorm:"column:processed_at_utc"`
}

func (Message) TableName() string {
	return "outbox"
}
