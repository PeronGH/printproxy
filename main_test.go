package main

import (
	"strings"
	"testing"
)

func TestBuildExports_Posix(t *testing.T) {
	cfg := Config{
		HTTP:    "127.0.0.1:2080",
		HTTPS:   "127.0.0.1:2080",
		NoProxy: []string{"localhost", "127.0.0.1"},
	}
	got := strings.Join(buildExports(cfg, Posix{}), "\n")
	want := "export http_proxy='http://127.0.0.1:2080'\n" +
		"export https_proxy='http://127.0.0.1:2080'\n" +
		"export no_proxy='localhost,127.0.0.1'"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestBuildExports_Pwsh(t *testing.T) {
	cfg := Config{
		HTTPS: "user:pa'ss@proxy:8080",
		SOCKS: "127.0.0.1:1080",
	}
	got := strings.Join(buildExports(cfg, Pwsh{}), "\n")
	want := "$env:https_proxy = 'http://user:pa''ss@proxy:8080'\n" +
		"$env:all_proxy = 'socks5://127.0.0.1:1080'"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestBuildExports_UnsetPosix(t *testing.T) {
	got := strings.Join(buildExports(Config{}, Posix{}), "\n")
	want := "unset http_proxy https_proxy ftp_proxy all_proxy no_proxy"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestBuildExports_UnsetPwsh(t *testing.T) {
	lines := buildExports(Config{}, Pwsh{})
	if len(lines) != len(proxyVars) {
		t.Fatalf("got %d lines, want %d", len(lines), len(proxyVars))
	}
	for i, name := range proxyVars {
		want := "Remove-Item Env:" + name + " -ErrorAction SilentlyContinue"
		if lines[i] != want {
			t.Errorf("line %d: got %q, want %q", i, lines[i], want)
		}
	}
}
