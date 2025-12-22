package email

import (
	"context"
	"encoding/json"

	events "github.com/hawful70/platform-events/pkg/events"
)

type UserCreatedHandler struct {
	mailer Mailer
}

func NewUserCreatedHandler(mailer Mailer) *UserCreatedHandler {
	return &UserCreatedHandler{mailer: mailer}
}

func (h *UserCreatedHandler) Handle(ctx context.Context, value []byte) error {
	var evt events.UserCreated
	if err := json.Unmarshal(value, &evt); err != nil {
		return err
	}
	if evt.Type != events.UserCreatedType {
		return nil
	}
	return h.mailer.SendWelcome(ctx, evt.User.Email, evt.User.Username)
}
