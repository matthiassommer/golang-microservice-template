package utils

// closable is an interface for all closable streams and resources.
type closable interface {
	Close() error
}

// Close a closable instance, check for errors and print them to console.
func Close(c closable) {
	if err := c.Close(); err != nil {
		Log.Error(err)
	}
}
