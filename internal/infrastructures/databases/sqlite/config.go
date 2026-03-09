package sqlite

// Config holds the configuration for the SQLite database adapter.
type Config struct {
	// Path is the SQLite file path. Use ":memory:" for tests.
	Path string
}
