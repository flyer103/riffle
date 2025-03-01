package serving

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

// ServerOptions contains the options for the server
type ServerOptions struct {
	Port         int           `json:"port"`
	DBPath       string        `json:"dbPath"`
	LogLevel     string        `json:"logLevel"`
	EnablePprof  bool          `json:"enablePprof"`
	MetricsPort  int           `json:"metricsPort"`
	RateLimit    int           `json:"rateLimit"`
	EnableCORS   bool          `json:"enableCORS"`
	CORSOrigins  []string      `json:"corsOrigins"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
}

// NewServerOptions creates a new ServerOptions with default values
func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		Port:         8080,
		DBPath:       "./riffle.db",
		LogLevel:     "info",
		EnablePprof:  false,
		MetricsPort:  0,
		RateLimit:    100,
		EnableCORS:   false,
		CORSOrigins:  []string{"*"},
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// AddFlags adds flags to the given FlagSet
func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&o.Port, "port", o.Port, "Port to listen on")
	fs.StringVar(&o.DBPath, "db-path", o.DBPath, "Path to the SQLite database file")
	fs.StringVar(&o.LogLevel, "log-level", o.LogLevel, "Log level (debug, info, warn, error)")
	fs.BoolVar(&o.EnablePprof, "enable-pprof", o.EnablePprof, "Enable pprof debugging endpoints")
	fs.IntVar(&o.MetricsPort, "metrics-port", o.MetricsPort, "Port for Prometheus metrics (0 to disable)")
	fs.IntVar(&o.RateLimit, "rate-limit", o.RateLimit, "Rate limit in requests per second (0 to disable)")
	fs.BoolVar(&o.EnableCORS, "enable-cors", o.EnableCORS, "Enable CORS")
	fs.StringSliceVar(&o.CORSOrigins, "cors-origins", o.CORSOrigins, "Allowed CORS origins")
	fs.DurationVar(&o.ReadTimeout, "read-timeout", o.ReadTimeout, "HTTP server read timeout")
	fs.DurationVar(&o.WriteTimeout, "write-timeout", o.WriteTimeout, "HTTP server write timeout")
}

// Complete completes the options
func (o *ServerOptions) Complete() error {
	return nil
}

// Validate validates the options
func (o *ServerOptions) Validate() error {
	if o.Port < 1 || o.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if o.MetricsPort < 0 || o.MetricsPort > 65535 {
		return fmt.Errorf("metrics port must be between 0 and 65535")
	}

	if o.RateLimit < 0 {
		return fmt.Errorf("rate limit must be greater than or equal to 0")
	}

	if o.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be greater than 0")
	}

	if o.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be greater than 0")
	}

	return nil
}
