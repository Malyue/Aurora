package im

import "github.com/sirupsen/logrus"

type Options struct {
	ConfigFilePath string
	Logger         *logrus.Logger
	BeforeStart    []func() error
	BeforeStop     []func() error
	AfterStart     []func() error
	AfterStop      []func() error
}

type OptionFunc func(*Options)
