#!/usr/bin/env python3
"""
VFS Global Passport Application Status Tracker
A command-line tool to check passport application status from VFS Global.
"""

import sys
import time
import click
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from selenium.common.exceptions import TimeoutException, NoSuchElementException
from webdriver_manager.chrome import ChromeDriverManager
from colorama import init, Fore, Style

# Initialize colorama
init(autoreset=True)

VFS_TRACKING_URL = "https://www.vfsvisaonline.com/Global-PassportTracking/"


class PassportTracker:
    """Handles passport application tracking from VFS Global."""

    def __init__(self, headless=False, wait_time=30):
        """
        Initialize the tracker.

        Args:
            headless: Run browser in headless mode (not recommended due to CAPTCHA)
            wait_time: Maximum wait time for CAPTCHA solving (seconds)
        """
        self.headless = headless
        self.wait_time = wait_time
        self.driver = None

    def _setup_driver(self):
        """Setup and configure Chrome WebDriver."""
        chrome_options = Options()

        if self.headless:
            chrome_options.add_argument("--headless")
            print(f"{Fore.YELLOW}Warning: Headless mode may not work with reCAPTCHA{Style.RESET_ALL}")

        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument("--disable-blink-features=AutomationControlled")
        chrome_options.add_experimental_option("excludeSwitches", ["enable-automation"])
        chrome_options.add_experimental_option('useAutomationExtension', False)

        # User agent to appear more like a real browser
        chrome_options.add_argument(
            "user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
            "(KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
        )

        try:
            service = Service(ChromeDriverManager().install())
            self.driver = webdriver.Chrome(service=service, options=chrome_options)
            self.driver.execute_script(
                "Object.defineProperty(navigator, 'webdriver', {get: () => undefined})"
            )
        except Exception as e:
            print(f"{Fore.RED}Error setting up Chrome driver: {e}{Style.RESET_ALL}")
            sys.exit(1)

    def _wait_for_captcha_solution(self):
        """Wait for user to solve the CAPTCHA."""
        print(f"\n{Fore.CYAN}{'='*60}{Style.RESET_ALL}")
        print(f"{Fore.YELLOW}Please solve the reCAPTCHA in the browser window...{Style.RESET_ALL}")
        print(f"{Fore.CYAN}{'='*60}\n{Style.RESET_ALL}")

        try:
            # Wait for the reCAPTCHA response token to appear
            WebDriverWait(self.driver, self.wait_time).until(
                lambda driver: driver.execute_script(
                    "return document.getElementById('g-recaptcha-response').value.length > 0"
                )
            )
            print(f"{Fore.GREEN}✓ CAPTCHA solved!{Style.RESET_ALL}\n")
            time.sleep(1)  # Small delay to ensure stability
            return True
        except TimeoutException:
            print(f"{Fore.RED}✗ CAPTCHA not solved within {self.wait_time} seconds{Style.RESET_ALL}")
            return False

    def track_application(self, reference_number, date_of_birth):
        """
        Track passport application status.

        Args:
            reference_number: Application reference number
            date_of_birth: Date of birth in DD/MM/YYYY format

        Returns:
            Status message or None if failed
        """
        try:
            self._setup_driver()

            print(f"{Fore.CYAN}Opening VFS Global tracking page...{Style.RESET_ALL}")
            self.driver.get(VFS_TRACKING_URL)

            # Wait for page to load
            wait = WebDriverWait(self.driver, 10)

            # Fill in the reference number
            print(f"{Fore.CYAN}Entering application details...{Style.RESET_ALL}")
            ref_input = wait.until(
                EC.presence_of_element_located((By.NAME, "ApplicationNo"))
            )
            ref_input.clear()
            ref_input.send_keys(reference_number)

            # Fill in date of birth
            dob_input = self.driver.find_element(By.NAME, "DateOfBirth")
            dob_input.clear()
            dob_input.send_keys(date_of_birth)

            # Wait for user to solve CAPTCHA
            if not self._wait_for_captcha_solution():
                return None

            # Click submit button
            print(f"{Fore.CYAN}Submitting form...{Style.RESET_ALL}")
            submit_button = self.driver.find_element(
                By.XPATH, "//button[contains(text(), 'SUBMIT') or @type='submit']"
            )
            submit_button.click()

            # Wait for results
            time.sleep(3)

            # Try to extract status message
            try:
                # Look for the blue text message
                status_element = wait.until(
                    EC.presence_of_element_located(
                        (By.XPATH, "//p[contains(@style, 'color') or contains(text(), 'application reference')]")
                    )
                )
                status_text = status_element.text

                # Also check for any other status information
                try:
                    additional_info = self.driver.find_element(
                        By.XPATH, "//div[contains(@class, 'status') or contains(@class, 'result')]"
                    ).text
                    if additional_info and additional_info != status_text:
                        status_text += f"\n\n{additional_info}"
                except NoSuchElementException:
                    pass

                return status_text

            except TimeoutException:
                # If no status found, check for error messages
                try:
                    error_msg = self.driver.find_element(
                        By.XPATH, "//div[contains(@class, 'error') or contains(@class, 'alert')]"
                    ).text
                    return f"Error: {error_msg}"
                except NoSuchElementException:
                    # Return the page source for debugging
                    return "Could not find status message. The page may have changed."

        except Exception as e:
            print(f"{Fore.RED}Error during tracking: {e}{Style.RESET_ALL}")
            return None
        finally:
            if self.driver:
                input(f"\n{Fore.YELLOW}Press Enter to close the browser...{Style.RESET_ALL}")
                self.driver.quit()


@click.command()
@click.option(
    '--reference', '-r',
    required=True,
    help='Application reference number (e.g., 25-2055287305)'
)
@click.option(
    '--dob', '-d',
    required=True,
    help='Date of birth in DD/MM/YYYY format (e.g., 12/09/1994)'
)
@click.option(
    '--headless',
    is_flag=True,
    default=False,
    help='Run in headless mode (not recommended, CAPTCHA requires manual solving)'
)
@click.option(
    '--wait-time', '-w',
    default=120,
    type=int,
    help='Maximum time to wait for CAPTCHA solving in seconds (default: 120)'
)
def main(reference, dob, headless, wait_time):
    """
    VFS Global Passport Application Status Tracker

    Track your passport application status from VFS Global.

    Example usage:
        python passport_tracker.py -r 25-2055287305 -d 12/09/1994
    """
    print(f"\n{Fore.CYAN}{'='*60}{Style.RESET_ALL}")
    print(f"{Fore.GREEN}VFS Global Passport Application Tracker{Style.RESET_ALL}")
    print(f"{Fore.CYAN}{'='*60}{Style.RESET_ALL}\n")

    print(f"Reference Number: {Fore.YELLOW}{reference}{Style.RESET_ALL}")
    print(f"Date of Birth: {Fore.YELLOW}{dob}{Style.RESET_ALL}\n")

    tracker = PassportTracker(headless=headless, wait_time=wait_time)
    status = tracker.track_application(reference, dob)

    if status:
        print(f"\n{Fore.CYAN}{'='*60}{Style.RESET_ALL}")
        print(f"{Fore.GREEN}Application Status:{Style.RESET_ALL}")
        print(f"{Fore.CYAN}{'='*60}{Style.RESET_ALL}")
        print(f"\n{status}\n")
        print(f"{Fore.CYAN}{'='*60}{Style.RESET_ALL}\n")
    else:
        print(f"\n{Fore.RED}Failed to retrieve application status.{Style.RESET_ALL}\n")
        sys.exit(1)


if __name__ == "__main__":
    main()
