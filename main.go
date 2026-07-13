package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HTTP    string
	HTTPS   string
	FTP     string
	SOCKS   string
	NoProxy []string
}

func main() {
	cfg, err := readConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "read proxy config: %v\n", err)
		os.Exit(1)
	}
	for _, line := range buildExports(cfg) {
		fmt.Println(line)
	}
}

func buildExports(cfg Config) []string {
	if cfg.HTTP == "" && cfg.HTTPS == "" && cfg.FTP == "" && cfg.SOCKS == "" {
		return []string{"unset http_proxy https_proxy ftp_proxy all_proxy no_proxy"}
	}
	var out []string
	emit := func(name, scheme, addr string) {
		if addr == "" {
			return
		}
		if !strings.Contains(addr, "://") {
			addr = scheme + "://" + addr
		}
		out = append(out, fmt.Sprintf("export %s=%s", name, shellQuote(addr)))
	}
	emit("http_proxy", "http", cfg.HTTP)
	emit("https_proxy", "http", cfg.HTTPS)
	emit("ftp_proxy", "http", cfg.FTP)
	emit("all_proxy", "socks5", cfg.SOCKS)
	if len(cfg.NoProxy) > 0 {
		out = append(out, fmt.Sprintf("export no_proxy=%s", shellQuote(strings.Join(cfg.NoProxy, ","))))
	}
	return out
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
