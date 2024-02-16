package redis

// Defines the types of configuration.
const (
	TypeRedis    = "redis"
	TypeSentinel = "sentinel"
)

// An Config is the structure format for Redis vault secret.
type Config struct {
	Type               string   `json:"type"`
	SentinelMasterName string   `json:"sentinel_master_name"`
	SentinelHosts      []string `json:"sentinel_hosts"`
	SentinelPassword   string   `json:"sentinel_password"`
	Host               string   `json:"host"`
	Password           string   `json:"password"`
	Database           int      `json:"database"`
	TLS                bool     `json:"tls"`
	InsecureSkipVerify bool     `json:"insecure_skip_verify"`
}
