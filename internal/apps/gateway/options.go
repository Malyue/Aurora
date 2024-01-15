package gateway

import "github.com/sirupsen/logrus"

/*
	use options func to init
*/

type Options struct {
	ConfigFilePath string
	Logger         *logrus.Logger
	BeforeStart    []func() error
	BeforeStop     []func() error
	AfterStart     []func() error
	AfterStop      []func() error
}

type OptionFunc func(*Options)

func WithConfig(filepath string) OptionFunc {
	return func(o *Options) {
		o.ConfigFilePath = filepath
	}
}

func WithBeforeStartFunc(beforeStartFunc []func() error) OptionFunc {
	return func(o *Options) {
		o.BeforeStart = beforeStartFunc
	}
}

func WithBeforeStopFunc(beforeStopFunc []func() error) OptionFunc {
	return func(o *Options) {
		o.BeforeStart = beforeStopFunc
	}
}

func WithAfterStartFunc(afterStartFunc []func() error) OptionFunc {
	return func(o *Options) {
		o.AfterStart = afterStartFunc
	}
}

func WithAfterStopFunc(afterStopFunc []func() error) OptionFunc {
	return func(o *Options) {
		o.AfterStart = afterStopFunc
	}
}

func WithLogger(logger *logrus.Logger) OptionFunc {
	return func(o *Options) {
		o.Logger = logger
	}
}
