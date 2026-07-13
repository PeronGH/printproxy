# printproxy

Prints the system proxy settings as POSIX-shell `export` lines on Windows
and macOS, so tools like `curl`, `git`, and `wget` honor them. Originally
built for [busybox-w32](https://frippery.org/busybox/) ash.

## Usage

POSIX shells (bash, ash, zsh, …):

```sh
eval "$(printproxy)"      # or printproxy.exe on Windows
```

PowerShell:

```powershell
printproxy.exe -shell pwsh | Invoke-Expression
```

When no proxy is enabled, the program prints `unset` / `Remove-Item` lines
instead, so the same command clears the variables.

## Example

```console
$ printproxy
export http_proxy='http://127.0.0.1:2080'
export https_proxy='http://127.0.0.1:2080'
export ftp_proxy='http://127.0.0.1:2080'
export no_proxy='localhost,127.0.0.1,::1'
```

Settings are read from:

- **Windows**: `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`
- **macOS**: `SCDynamicStoreCopyProxies` (SystemConfiguration framework)

## Install

Grab the binary for your platform from
[Releases](../../releases/latest) and drop it on your `PATH`:

- `printproxy-windows-amd64.exe` / `printproxy-windows-arm64.exe`
- `printproxy-darwin-amd64` / `printproxy-darwin-arm64`

## Build

```sh
# Windows (any host)
GOOS=windows GOARCH=amd64 go build .

# macOS (requires cgo + Xcode CLT)
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build .
```
