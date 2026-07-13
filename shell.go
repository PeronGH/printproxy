package main

import (
	"fmt"
	"strings"
)

type Shell interface {
	Set(name, value string) string
	UnsetAll(names []string) []string
}

func pickShell(name string) (Shell, error) {
	switch name {
	case "sh":
		return Posix{}, nil
	case "pwsh":
		return Pwsh{}, nil
	}
	return nil, fmt.Errorf("unknown shell %q (want sh or pwsh)", name)
}

type Posix struct{}

func (Posix) Set(name, value string) string {
	return "export " + name + "=" + posixQuote(value)
}

func (Posix) UnsetAll(names []string) []string {
	return []string{"unset " + strings.Join(names, " ")}
}

func posixQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

type Pwsh struct{}

func (Pwsh) Set(name, value string) string {
	return "$env:" + name + " = " + pwshQuote(value)
}

func (Pwsh) UnsetAll(names []string) []string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = "Remove-Item Env:" + n + " -ErrorAction SilentlyContinue"
	}
	return out
}

func pwshQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
