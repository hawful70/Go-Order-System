package events

const UserCreatedType = "user_created"

type UserCreated struct {
	Type string      `json:"type"`
	User UserPayload `json:"user"`
}

type UserPayload struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func NewUserCreated(id, email, username string) UserCreated {
	return UserCreated{
		Type: UserCreatedType,
		User: UserPayload{ID: id, Email: email, Username: username},
	}
}
