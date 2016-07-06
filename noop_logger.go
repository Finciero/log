package log

// NOOPLogger ...
type NOOPLogger struct{}

// NewNOOPContext returns a NOOPLogger, perfect for testing
func NewNOOPContext(_ ...interface{}) *NOOPLogger {
	return &NOOPLogger{}
}

// With does nothing and returns it self
func (n *NOOPLogger) With(_ ...interface{}) Logger {
	return n
}

// Log does nothing and returns nil
func (n *NOOPLogger) Log(_ ...interface{}) error {
	return nil
}

// Info does nothing and returns nil
func (n *NOOPLogger) Info(_ ...interface{}) error {
	return nil
}

// Error does nothing and returns nil
func (n *NOOPLogger) Error(_ error, _ ...interface{}) error {
	return nil
}

// Fatal does nothing and returns nil
func (n *NOOPLogger) Fatal(_ error, _ ...interface{}) error {
	return nil
}

// Warn does nothing and returns nil
func (n *NOOPLogger) Warn(_ ...interface{}) error {
	return nil
}
