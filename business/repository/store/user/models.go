package user

// User represents an individual User.
type User struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
}

// UpdateUser defines what information may be provided to modify an existing User.
type UpdateUser struct {
	Name  *string  `json:"name"`
	Email *string  `json:"email" validate:"omitempty,email"`
	Roles []string `json:"roles"`
}
