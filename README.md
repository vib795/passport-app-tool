# VFS Global Passport Application Tracker

A command-line tool to track passport application status from VFS Global's tracking portal.

**Available in both Python and Go!** See [COMPARISON.md](COMPARISON.md) to choose the right version for you.

## Features

- ✅ Automated form filling
- ✅ Interactive CAPTCHA solving
- ✅ Colored CLI output
- ✅ Error handling and validation
- ✅ Cross-platform support (Windows, macOS, Linux)
- ✅ Two implementations: Python (flexibility) and Go (performance)

## Prerequisites

- Python 3.8 or higher
- Google Chrome browser installed
- Internet connection

## Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd passport-app-tool
```

2. Install dependencies:
```bash
pip install -r requirements.txt
```

Or using a virtual environment (recommended):
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt
```

## Usage

### Basic Usage

```bash
python passport_tracker.py -r <REFERENCE_NUMBER> -d <DATE_OF_BIRTH>
```

### Example

```bash
python passport_tracker.py -r 25-2055287305 -d 12/09/1994
```

### Command-Line Options

| Option | Short | Description | Required |
|--------|-------|-------------|----------|
| `--reference` | `-r` | Application reference number | Yes |
| `--dob` | `-d` | Date of birth (DD/MM/YYYY format) | Yes |
| `--wait-time` | `-w` | Max wait time for CAPTCHA (seconds, default: 120) | No |
| `--headless` | - | Run in headless mode (not recommended) | No |

### Advanced Usage

Wait up to 3 minutes for CAPTCHA:
```bash
python passport_tracker.py -r 25-2055287305 -d 12/09/1994 --wait-time 180
```

## How It Works

1. **Opens Browser**: Launches Chrome and navigates to VFS tracking page
2. **Fills Form**: Automatically enters your reference number and date of birth
3. **CAPTCHA Solving**: Waits for you to manually solve the reCAPTCHA
4. **Submits Form**: Clicks submit after CAPTCHA is solved
5. **Displays Status**: Shows your application status in the terminal

## Important Notes

### About reCAPTCHA

The VFS Global website uses reCAPTCHA v2 to prevent automated access. This tool:
- Opens a visible browser window
- Waits for you to manually solve the CAPTCHA
- Automatically submits the form once CAPTCHA is solved

**You need to actively solve the CAPTCHA** when the browser window appears.

### Headless Mode

The `--headless` flag is provided but **not recommended** because:
- reCAPTCHA typically doesn't work in headless browsers
- Google's CAPTCHA system detects headless browsers
- You won't be able to see or solve the CAPTCHA

Only use headless mode if you're testing or if the website behavior changes.

## Troubleshooting

### Chrome Driver Issues

If you get errors about ChromeDriver:
```bash
pip install --upgrade webdriver-manager
```

The tool automatically downloads and manages ChromeDriver.

### CAPTCHA Timeout

If the CAPTCHA times out:
- Increase wait time: `--wait-time 300` (5 minutes)
- Make sure you're solving the CAPTCHA when the browser opens
- Check your internet connection

### "Element not found" Errors

This may happen if VFS Global changes their website structure:
1. Check if the website is accessible manually
2. Report the issue with screenshots
3. The tool may need updates to match new HTML structure

### SSL/Certificate Errors

If you encounter SSL errors:
```bash
pip install --upgrade certifi
```

## Development

### Project Structure

```
passport-app-tool/
├── passport_tracker.py    # Main CLI tool
├── requirements.txt       # Python dependencies
├── .gitignore            # Git ignore patterns
└── README.md             # This file
```

### Testing

Test with your own application details:
```bash
python passport_tracker.py -r YOUR_REF_NUMBER -d YOUR_DOB
```

### Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## Alternative: Go Implementation

If you prefer Go, see [`go_implementation/`](go_implementation/) directory for a Go version using chromedp.

**Benefits of Go version:**
- Single standalone binary (no dependencies)
- Faster execution and lower memory usage
- Easy cross-platform distribution
- See [COMPARISON.md](COMPARISON.md) for detailed comparison

**Quick start with Go:**
```bash
cd go_implementation
go build -o passport-tracker
./passport-tracker -r 25-2055287305 -d 12/09/1994
```

## Privacy & Security

- This tool runs locally on your machine
- No data is sent to third parties
- Only communicates with VFS Global's official website
- Your credentials are not stored

## Limitations

- Requires manual CAPTCHA solving (by design, to respect website's anti-bot measures)
- Needs Chrome browser installed
- Requires active internet connection
- Subject to VFS Global website changes

## License

MIT License - See LICENSE file for details

## Disclaimer

This tool is for personal use only. Always respect the website's terms of service and rate limits. The authors are not responsible for any misuse of this tool.

## Support

For issues or questions:
1. Check the Troubleshooting section
2. Open an issue on GitHub
3. Provide error messages and screenshots

## Changelog

### Version 1.0.0 (Initial Release)
- Basic tracking functionality
- Chrome automation with Selenium
- CLI interface with Click
- Manual CAPTCHA solving support
- Colored terminal output

---

Made with ❤️ for easier passport tracking
