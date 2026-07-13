package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const settingsPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

func main() {
	key, err := registry.OpenKey(registry.CURRENT_USER, settingsPath, registry.QUERY_VALUE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open registry: %v\n", err)
		os.Exit(1)
	}
	defer key.Close()

	enabled, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read ProxyEnable: %v\n", err)
		os.Exit(1)
	}

	if enabled == 0 {
		fmt.Println("unset http_proxy https_proxy ftp_proxy all_proxy no_proxy")
		return
	}

	server, _, err := key.GetStringValue("ProxyServer")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read ProxyServer: %v\n", err)
		os.Exit(1)
	}

	override, _, _ := key.GetStringValue("ProxyOverride")

	for _, line := range buildExports(server, override) {
		fmt.Println(line)
	}
}

func buildExports(server, override string) []string {
	var out []string
	proxies := parseProxyServer(server)
	emit := func(name, scheme, addr string) {
		if addr == "" {
			return
		}
		if !strings.Contains(addr, "://") {
			addr = scheme + "://" + addr
		}
		out = append(out, fmt.Sprintf("export %s=%s", name, shellQuote(addr)))
	}
	emit("http_proxy", "http", proxies["http"])
	emit("https_proxy", "http", proxies["https"])
	emit("ftp_proxy", "http", proxies["ftp"])
	emit("all_proxy", "socks5", proxies["socks"])
	if np := convertNoProxy(override); np != "" {
		out = append(out, fmt.Sprintf("export no_proxy=%s", shellQuote(np)))
	}
	return out
}

func parseProxyServer(s string) map[string]string {
	m := map[string]string{}
	if strings.Contains(s, "=") {
		for _, part := range strings.Split(s, ";") {
			k, v, ok := strings.Cut(strings.TrimSpace(part), "=")
			if !ok {
				continue
			}
			m[strings.ToLower(strings.TrimSpace(k))] = strings.TrimSpace(v)
		}
		return m
	}
	m["http"] = s
	m["https"] = s
	m["ftp"] = s
	return m
}

func convertNoProxy(s string) string {
	var out []string
	for _, p := range strings.Split(s, ";") {
		p = strings.TrimSpace(p)
		switch p {
		case "":
			continue
		case "<local>":
			out = append(out, "localhost", "127.0.0.1", "::1")
		default:
			out = append(out, p)
		}
	}
	return strings.Join(out, ",")
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
