package config

import (
	"errors"
)

var (
	ErrEnvironmentVariableNotFound = errors.New("no environment variable was found for the key specified")
)