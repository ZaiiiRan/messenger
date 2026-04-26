package implkafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/producers/models"
)

func (p *Producer) ProduceEmailCodeTask(ctx context.Context, email string, c *code.Code) error {
	msg := models.CodeMessageFromDomain(c, email)

	value, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("ProduceEmailCodeTask: marshal: %w", err)
	}

	return p.Produce(ctx, Message{
		Key:   c.GetUserID(),
		Value: string(value),
	})
}
