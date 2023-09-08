package common

// User represents an individual user in the system.
type User struct {
	// ID is the unique identifier for the user
	ID int64 `json:"id"`
	// Name is the name of the user
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
