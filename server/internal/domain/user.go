package domain

type User struct {
	ID         string
	Email      string
	Name       string
	APIKeyHash string
	IsAdmin    bool
}
