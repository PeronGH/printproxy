# printproxy

Prints Windows proxy settings as POSIX-shell `export` lines. Designed for
[busybox-w32](https://frippery.org/busybox/) ash, where sourcing them makes
tools like `curl`, `git`, and `wget` honor the system proxy.

## Usage

```sh
eval "$(printproxy.exe)"
```

When the system proxy is disabled, the program prints `unset ...` instead,
so the same command clears the variables.

## Example

```console
$ printproxy.exe
export http_proxy='http://127.0.0.1:2080'
export https_proxy='http://127.0.0.1:2080'
export ftp_proxy='http://127.0.0.1:2080'
export no_proxy='localhost,127.0.0.1,::1'
```

Settings are read from
`HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`.

## Install

Grab `printproxy-windows-amd64.exe` (or `-arm64.exe`) from
[Releases](../../releases/latest) and drop it on your `PATH`.

## Build

```sh
GOOS=windows GOARCH=amd64 go build .
```
