package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const vfsTrackingURL = "https://www.vfsvisaonline.com/Global-PassportTracking/"

var (
	referenceNumber string
	dateOfBirth     string
	waitTime        int
	headless        bool
)

type PassportTracker struct {
	referenceNumber string
	dateOfBirth     string
	waitTime        time.Duration
	headless        bool
}

func NewPassportTracker(ref, dob string, wait int, headless bool) *PassportTracker {
	return &PassportTracker{
		referenceNumber: ref,
		dateOfBirth:     dob,
		waitTime:        time.Duration(wait) * time.Second,
		headless:        headless,
	}
}

func (pt *PassportTracker) TrackApplication() (string, error) {
	// Setup Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", pt.headless),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, pt.waitTime+30*time.Second)
	defer cancel()

	var statusText string

	color.Cyan("Opening VFS Global tracking page...")

	err := chromedp.Run(ctx,
		chromedp.Navigate(vfsTrackingURL),
		chromedp.WaitVisible(`input[name="ApplicationNo"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="ApplicationNo"]`, pt.referenceNumber, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="DateOfBirth"]`, pt.dateOfBirth, chromedp.ByQuery),
	)

	if err != nil {
		return "", fmt.Errorf("error filling form: %w", err)
	}

	color.Yellow("\n" + strings.Repeat("=", 60))
	color.Yellow("Please solve the reCAPTCHA in the browser window...")
	color.Yellow(strings.Repeat("=", 60) + "\n")

	// Wait for CAPTCHA to be solved by polling for a value in the response field
	startTime := time.Now()
	captchaSolved := false

	for time.Since(startTime) < pt.waitTime {
		var captchaValue string
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.getElementById('g-recaptcha-response').value`, &captchaValue),
		)

		if err == nil && len(captchaValue) > 0 {
			captchaSolved = true
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	if !captchaSolved {
		return "", fmt.Errorf("CAPTCHA not solved within %v seconds", int(pt.waitTime.Seconds()))
	}

	// Small delay to ensure stability
	time.Sleep(1 * time.Second)

	color.Green("✓ CAPTCHA solved!\n")
	color.Cyan("Submitting form...")

	// Submit form
	err = chromedp.Run(ctx,
		chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second), // Wait for response
	)

	if err != nil {
		return "", fmt.Errorf("error submitting form: %w", err)
	}

	// Try to extract status message or error
	var foundStatus bool

	// First check for error messages
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`
			(function() {
				// Check for "Invalid Request" or error messages
				var headings = document.querySelectorAll('h1, h2, h3');
				for (var i = 0; i < headings.length; i++) {
					if (headings[i].textContent.includes('Invalid Request') ||
					    headings[i].textContent.includes('Error')) {
						return headings[i].textContent.trim();
					}
				}

				// Look for status message in paragraphs
				var paragraphs = document.querySelectorAll('p');
				for (var i = 0; i < paragraphs.length; i++) {
					var text = paragraphs[i].textContent.trim();
					if (text.length > 10 &&
					    (text.toLowerCase().includes('application') ||
					     text.toLowerCase().includes('passport') ||
					     text.toLowerCase().includes('reference'))) {
						return text;
					}
				}

				return '';
			})()
		`, &statusText),
	)

	if err == nil && len(statusText) > 0 {
		foundStatus = true
		// Check if it's an error message
		if strings.Contains(statusText, "Invalid Request") {
			return "", fmt.Errorf("form submission failed: %s - Please verify your reference number and date of birth are correct", statusText)
		}
	}

	if !foundStatus {
		return "", fmt.Errorf("could not extract status message from the page")
	}

	// Keep browser open for user to see
	color.Yellow("\nPress Ctrl+C to close the browser...")
	select {}

	return statusText, nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "passport-tracker",
		Short: "VFS Global Passport Application Status Tracker",
		Long: `A command-line tool to track passport application status from VFS Global.

Example:
  passport-tracker -r 25-2055287305 -d 12/09/1994`,
		Run: func(cmd *cobra.Command, args []string) {
			runTracker()
		},
	}

	rootCmd.Flags().StringVarP(&referenceNumber, "reference", "r", "", "Application reference number (required)")
	rootCmd.Flags().StringVarP(&dateOfBirth, "dob", "d", "", "Date of birth in DD/MM/YYYY format (required)")
	rootCmd.Flags().IntVarP(&waitTime, "wait-time", "w", 120, "Maximum time to wait for CAPTCHA solving (seconds)")
	rootCmd.Flags().BoolVar(&headless, "headless", false, "Run in headless mode (not recommended)")

	rootCmd.MarkFlagRequired("reference")
	rootCmd.MarkFlagRequired("dob")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runTracker() {
	color.Cyan("\n" + strings.Repeat("=", 60))
	color.Green("VFS Global Passport Application Tracker")
	color.Cyan(strings.Repeat("=", 60) + "\n")

	fmt.Printf("Reference Number: %s\n", color.YellowString(referenceNumber))
	fmt.Printf("Date of Birth: %s\n\n", color.YellowString(dateOfBirth))

	tracker := NewPassportTracker(referenceNumber, dateOfBirth, waitTime, headless)
	status, err := tracker.TrackApplication()

	if err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}

	color.Cyan("\n" + strings.Repeat("=", 60))
	color.Green("Application Status:")
	color.Cyan(strings.Repeat("=", 60))
	fmt.Printf("\n%s\n\n", status)
	color.Cyan(strings.Repeat("=", 60) + "\n")
}
