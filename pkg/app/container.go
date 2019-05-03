package app

import "github.com/sirupsen/logrus"

// Container is a basic application container
type Container struct {
	cfg Config

	// panicCatchers allows control of panic handling
	panicCatchers []PanicCatcher

	logger logrus.FieldLogger
	closer func() error
}

func newContainer(cfg Config) *Container {
	return &Container{
		cfg: cfg,
	}
}

// Cfg returns app-level base configuration
func (c *Container) Cfg() Config {
	return c.cfg
}

// WithLogger sets logger instance
func (c *Container) WithLogger(logger logrus.FieldLogger) {
	c.logger = logger
}

// Logger returns app-level logger
func (c *Container) Logger() logrus.FieldLogger {
	if c.logger == nil {
		c.logger = logrus.New()
	}

	return c.logger
}

// SetCloserFunc enables service locator termination with callback
func (c *Container) SetCloserFunc(closer func() error) {
	c.closer = closer
}

// Close invokes service locator termination
func (c *Container) Close() error {
	if c.closer != nil {
		return c.closer()
	}

	return nil
}

// WithPanicCatcher sets control of panic handling
func (c *Container) WithPanicCatcher(panicCatchers []PanicCatcher) {
	c.panicCatchers = panicCatchers
}

// PanicCatcher returns app-level control of panic handling
func (c *Container) PanicCatcher() []PanicCatcher {
	return c.panicCatchers
}
