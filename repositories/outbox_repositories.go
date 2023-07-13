package repositories

import (
	"wa-blast/models"
)

// MessageRepository ...
type MessageRepository interface {
	InsertMessageOutbox(Data models.Outbox) error
}

type messageRepository struct{}

func initMessageRepository() MessageRepository {
	// Prepare statements
	var r messageRepository
	return &r
}

func (s *messageRepository) InsertMessageOutbox(Data models.Outbox) error {
	err := DBConnect.Table("outbox").Create(&Data).Error
	return err
}
