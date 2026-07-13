package main

import (
	"flag"
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

var proxyVars = []string{"http_proxy", "https_proxy", "ftp_proxy", "all_proxy", "no_proxy"}

func main() {
	shellFlag := flag.String("shell", "sh", "output syntax: sh or pwsh")
	flag.Parse()

	sh, err := pickShell(*shellFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	cfg, err := readConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "read proxy config: %v\n", err)
		os.Exit(1)
	}
	for _, line := range buildExports(cfg, sh) {
		fmt.Println(line)
	}
}

func buildExports(cfg Config, sh Shell) []string {
	if cfg.HTTP == "" && cfg.HTTPS == "" && cfg.FTP == "" && cfg.SOCKS == "" {
		return sh.UnsetAll(proxyVars)
	}
	var out []string
	emit := func(name, scheme, addr string) {
		if addr == "" {
			return
		}
		if !strings.Contains(addr, "://") {
			addr = scheme + "://" + addr
		}
		out = append(out, sh.Set(name, addr))
	}
	emit("http_proxy", "http", cfg.HTTP)
	emit("https_proxy", "http", cfg.HTTPS)
	emit("ftp_proxy", "http", cfg.FTP)
	emit("all_proxy", "socks5", cfg.SOCKS)
	if len(cfg.NoProxy) > 0 {
		out = append(out, sh.Set("no_proxy", strings.Join(cfg.NoProxy, ",")))
	}
	return out
}
