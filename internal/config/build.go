package config

import "fmt"

const (
	Database = "test"
)

var (
	version  = "dev"
	revision = "none"
	date     = "unknown"
)

func Version() string {
	return fmt.Sprintf("%s (revision %.7s @ %s)", version, revision, date)
}
