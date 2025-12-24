package main

import (
	"context"
	"fmt"
	"io"
	"log"
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
	verbose         bool
)

type PassportTracker struct {
	referenceNumber string
	dateOfBirth     string
	waitTime        time.Duration
	headless        bool
	verbose         bool
}

func NewPassportTracker(ref, dob string, wait int, headless, verbose bool) *PassportTracker {
	return &PassportTracker{
		referenceNumber: ref,
		dateOfBirth:     dob,
		waitTime:        time.Duration(wait) * time.Second,
		headless:        headless,
		verbose:         verbose,
	}
}

func (pt *PassportTracker) TrackApplication() (string, error) {
	// Suppress chromedp verbose logging unless verbose flag is set
	if !pt.verbose {
		log.SetOutput(io.Discard)
	}

	// Setup Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", pt.headless),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create context with custom logging
	var ctx context.Context
	var ctxCancel context.CancelFunc
	if pt.verbose {
		ctx, ctxCancel = chromedp.NewContext(allocCtx)
	} else {
		ctx, ctxCancel = chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...interface{}) {}))
	}
	defer ctxCancel()

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

	// Wait for CAPTCHA to be solved
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`#g-recaptcha-response`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second), // Give user time to see the form
	)

	if err != nil {
		return "", fmt.Errorf("CAPTCHA not solved: %w", err)
	}

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

	// Extract status
	err = chromedp.Run(ctx,
		chromedp.Text(`p`, &statusText, chromedp.ByQuery),
	)

	if err != nil {
		return "", fmt.Errorf("error extracting status: %w", err)
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
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose chromedp debug output")

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

	tracker := NewPassportTracker(referenceNumber, dateOfBirth, waitTime, headless, verbose)
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
