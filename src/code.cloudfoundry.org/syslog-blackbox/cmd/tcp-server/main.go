package main

import (
	"log"
	"net/http"

	"code.cloudfoundry.org/syslog-blackbox/pkg/stat"
	"code.cloudfoundry.org/syslog-blackbox/pkg/web"

	"code.cloudfoundry.org/go-envstruct"
	"code.cloudfoundry.org/syslog-blackbox/pkg/syslog"
)

func main() {
	cfg := loadConfig()

	c := stat.NewCounter()
	r := web.NewRouter(c.Counts)
	l := syslog.NewListener(cfg.SyslogAddr, c.Add)

	l.Run(false)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, r))
}

// Config stores all configuration for the tcp-server.
type Config struct {
	SyslogAddr string `env:"SYSLOG_ADDR"`
	HTTPAddr   string `env:"HTTP_ADDR"`
}

func loadConfig() Config {
	var cfg Config
	if err := envstruct.Load(&cfg); err != nil {
		log.Fatalf("failed to load config from environment: %s", err)
	}

	return cfg
}
