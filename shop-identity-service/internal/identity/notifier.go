package identity

import "context"

type UserNotifier interface {
	UserCreated(ctx context.Context, user User) error
}

type noopNotifier struct{}

func (noopNotifier) UserCreated(ctx context.Context, user User) error {
	return nil
}

func NoopNotifier() UserNotifier {
	return noopNotifier{}
}
