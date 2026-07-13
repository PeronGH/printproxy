package main

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const settingsPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

func readConfig() (Config, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, settingsPath, registry.QUERY_VALUE)
	if err != nil {
		return Config{}, fmt.Errorf("open registry: %w", err)
	}
	defer key.Close()

	enabled, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		return Config{}, fmt.Errorf("read ProxyEnable: %w", err)
	}
	if enabled == 0 {
		return Config{}, nil
	}

	server, _, err := key.GetStringValue("ProxyServer")
	if err != nil {
		return Config{}, fmt.Errorf("read ProxyServer: %w", err)
	}
	override, _, _ := key.GetStringValue("ProxyOverride")

	cfg := Config{NoProxy: parseOverride(override)}
	if strings.Contains(server, "=") {
		for _, part := range strings.Split(server, ";") {
			k, v, ok := strings.Cut(strings.TrimSpace(part), "=")
			if !ok {
				continue
			}
			v = strings.TrimSpace(v)
			switch strings.ToLower(strings.TrimSpace(k)) {
			case "http":
				cfg.HTTP = v
			case "https":
				cfg.HTTPS = v
			case "ftp":
				cfg.FTP = v
			case "socks":
				cfg.SOCKS = v
			}
		}
	} else {
		cfg.HTTP = server
		cfg.HTTPS = server
		cfg.FTP = server
	}
	return cfg, nil
}

func parseOverride(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ";") {
		p = strings.TrimSpace(p)
		switch p {
		case "":
		case "<local>":
			out = append(out, "localhost", "127.0.0.1", "::1")
		default:
			out = append(out, p)
		}
	}
	return out
}
