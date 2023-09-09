package helper

type IPBinding string

const (
	DefaultServerEndpoint           = "localhost"
	DefaultServerPort               = "9090"
	DefaultMemcachePort             = "11211"
	DefaultDatabasePort             = "5432"
	LocalHostBinding      IPBinding = "127.0.0.1"
)
