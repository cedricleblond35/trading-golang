package config

import "fmt"

const (
	Database = "trading"
)

var (
	version  = "dev"
	revision = "none"
	date     = "unknown"
)

func Version() string {
	return fmt.Sprintf("%s (revision %.7s @ %s)", version, revision, date)
}
