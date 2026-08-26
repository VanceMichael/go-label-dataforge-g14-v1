package config

import (
	"errors"
	"net"
)

func (c Config) Validate() error {
	if c.ListenAddr == "" {
		return errors.New("listen address required")
	}
	if _, _, e := net.SplitHostPort(c.ListenAddr); e != nil {
		return e
	}
	if c.DatabasePath == "" {
		return errors.New("database path required")
	}
	if c.SessionTTL <= 0 {
		return errors.New("session ttl must be positive")
	}
	if c.WorkerInterval <= 0 {
		return errors.New("worker interval must be positive")
	}
	if len(c.SessionSecret) < 8 {
		return errors.New("session secret too short")
	}
	return nil
}
func (c Config) IsEphemeral() bool {
	return c.DatabasePath == ":memory:" || c.DatabasePath == "file::memory:?cache=shared"
}
func (c Config) PublicURL() string {
	host, port, e := net.SplitHostPort(c.ListenAddr)
	if e != nil {
		return ""
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return "http://" + net.JoinHostPort(host, port)
}
