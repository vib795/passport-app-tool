# VFS Global Passport Tracker (Go Implementation)

This is a Go implementation of the VFS Global passport application tracker using chromedp.

## Prerequisites

- Go 1.21 or higher
- Google Chrome browser installed

## Installation

1. Navigate to this directory:
```bash
cd go_implementation
```

2. Install dependencies:
```bash
go mod download
```

3. Build the binary:
```bash
go build -o passport-tracker
```

Or build for different platforms:

**Linux:**
```bash
GOOS=linux GOARCH=amd64 go build -o passport-tracker-linux
```

**macOS:**
```bash
GOOS=darwin GOARCH=amd64 go build -o passport-tracker-macos
```

**Windows:**
```bash
GOOS=windows GOARCH=amd64 go build -o passport-tracker.exe
```

## Usage

Run directly with Go:
```bash
go run main.go -r 25-2055287305 -d 12/09/1994
```

Or use the compiled binary:
```bash
./passport-tracker -r 25-2055287305 -d 12/09/1994
```

### Command-Line Options

```
  -r, --reference string   Application reference number (required)
  -d, --dob string         Date of birth in DD/MM/YYYY format (required)
  -w, --wait-time int      Maximum time to wait for CAPTCHA solving (seconds) (default 120)
      --headless           Run in headless mode (not recommended)
  -v, --verbose            Show verbose chromedp debug output
  -h, --help              help for passport-tracker
```

### Examples

Basic usage:
```bash
./passport-tracker -r 25-2055287305 -d 12/09/1994
```

With custom wait time:
```bash
./passport-tracker -r 25-2055287305 -d 12/09/1994 -w 180
```

With verbose debug output (for troubleshooting):
```bash
./passport-tracker -r 25-2055287305 -d 12/09/1994 --verbose
```

## Advantages of Go Implementation

1. **Single Binary**: Compiles to a single executable with no dependencies
2. **Cross-Platform**: Easy to build for different operating systems
3. **Performance**: Faster startup and lower memory usage
4. **Distribution**: Easier to distribute (just share the binary)

## Comparison with Python Version

| Feature | Python | Go |
|---------|--------|-----|
| Setup Complexity | Medium (pip install) | Low (single binary) |
| Runtime Speed | Slower | Faster |
| Memory Usage | Higher | Lower |
| Distribution | Requires Python | Single binary |
| Development | More libraries | More verbose |

## Dependencies

- **chromedp**: Chrome DevTools Protocol for browser automation
- **cobra**: CLI framework for building commands
- **color**: Terminal color output

## Building for Production

Create optimized builds:

```bash
# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o passport-tracker-linux

# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o passport-tracker-macos

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o passport-tracker.exe
```

The `-ldflags="-s -w"` flag removes debug information and symbol table to reduce binary size.

## Notes

- This implementation uses chromedp which controls Chrome via DevTools Protocol
- Requires Chrome/Chromium to be installed on the system
- CAPTCHA must still be solved manually
- Browser window will remain open until you press Ctrl+C

## Troubleshooting

### Chrome Not Found

If you get "chrome not found" errors:
- Make sure Chrome is installed
- Set the `CHROME_BIN` environment variable to your Chrome executable path

### Build Errors

If you encounter build errors:
```bash
go mod tidy
go clean -modcache
go mod download
```

## License

MIT License - Same as main project
