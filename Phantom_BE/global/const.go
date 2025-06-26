package global

type ContextKey string

const(
	// VERSION used to identify artifact version
	VERSION = "v1"

	// Middleware error handling keys
	DBErrorKey          ContextKey = "db_error"
	DBTimeoutKey        ContextKey = "connection_timeout"
	ExternalTimeoutKey  ContextKey = "ExternalTimeout"
	RequestTimeoutKey   ContextKey = "RequestTimeout"
	
	// Request api limit times and duration
    CtLayout = "2006-01-02 15:04:05"
	SampleDataDir = "/opt/data/"
)