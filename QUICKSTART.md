# Quick Start Guide

Get started with the VFS Global Passport Tracker in 5 minutes!

## Choose Your Version

### Option 1: Python (Recommended for Beginners)

**Step 1:** Make sure you have Python 3.8+ installed
```bash
python --version
```

**Step 2:** Install dependencies
```bash
pip install -r requirements.txt
```

**Step 3:** Run the tracker
```bash
python passport_tracker.py -r YOUR_REFERENCE_NUMBER -d YOUR_DATE_OF_BIRTH
```

**Example:**
```bash
python passport_tracker.py -r 25-2055287305 -d 12/09/1994
```

### Option 2: Go (Recommended for Distribution)

**Step 1:** Make sure you have Go 1.21+ installed
```bash
go version
```

**Step 2:** Build the binary
```bash
cd go_implementation
go build -o passport-tracker
```

**Step 3:** Run the tracker
```bash
./passport-tracker -r YOUR_REFERENCE_NUMBER -d YOUR_DATE_OF_BIRTH
```

**Example:**
```bash
./passport-tracker -r 25-2055287305 -d 12/09/1994
```

## What Happens Next?

1. **Browser Opens**: A Chrome window will open automatically
2. **Form Pre-filled**: Your reference number and date of birth are already entered
3. **Solve CAPTCHA**: Click the "I'm not a robot" checkbox
4. **Wait**: The tool will automatically submit the form after CAPTCHA is solved
5. **View Status**: Your application status will be displayed in the terminal

## Common Issues

### "Chrome not found"
**Solution:** Install Google Chrome browser
- Windows/Mac: Download from https://www.google.com/chrome/
- Linux: `sudo apt install google-chrome-stable` (Ubuntu/Debian)

### "Module not found" (Python)
**Solution:** Make sure you installed dependencies
```bash
pip install -r requirements.txt
```

### "Permission denied" (Go on Linux/Mac)
**Solution:** Make the binary executable
```bash
chmod +x passport-tracker
```

### CAPTCHA doesn't appear
**Solution:** Make sure you're not using `--headless` mode
```bash
# ❌ Don't do this
python passport_tracker.py -r XXX -d XXX --headless

# ✅ Do this instead
python passport_tracker.py -r XXX -d XXX
```

## Tips

### Increase CAPTCHA Wait Time
If you need more time to solve the CAPTCHA:
```bash
# Python
python passport_tracker.py -r XXX -d XXX --wait-time 300

# Go
./passport-tracker -r XXX -d XXX -w 300
```

### Run Without Installing Python Globally
Use a virtual environment:
```bash
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt
python passport_tracker.py -r XXX -d XXX
```

## Need Help?

1. Check the main [README.md](README.md) for detailed documentation
2. See [COMPARISON.md](COMPARISON.md) to decide between Python and Go
3. Review troubleshooting sections in respective READMEs
4. Open an issue on GitHub with your error message

## Example Session

```
$ python passport_tracker.py -r 25-2055287305 -d 12/09/1994

============================================================
VFS Global Passport Application Tracker
============================================================

Reference Number: 25-2055287305
Date of birth: 12/09/1994

Opening VFS Global tracking page...
Entering application details...

============================================================
Please solve the reCAPTCHA in the browser window...
============================================================

✓ CAPTCHA solved!

Submitting form...

============================================================
Application Status:
============================================================

Your application reference number 25-2055287305 has been
submitted at VFS Indian Consular Application Center for processing

============================================================

Press Enter to close the browser...
```

## Next Steps

- Save your reference number for future checks
- Set up a cron job/scheduled task for automatic checking (advanced)
- Customize the code for your specific needs
- Share the tool with others who might find it useful

---

**Ready to track your passport?** Just run the command and follow the prompts!
