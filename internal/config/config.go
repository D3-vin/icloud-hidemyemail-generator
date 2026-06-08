package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds the application configuration
type Config struct {
	CookieFile   string // Path to the cookie file
	OutputFile   string // Path to the output file for generated emails
	NoOutputFile bool   // Whether to skip writing to output file
	Cookies      string // Loaded and sanitized cookie string
}

// LoadConfig creates a new Config instance and loads cookies from the specified file
// Default cookie file path is "cookies.txt" in the current working directory
func LoadConfig(cookieFile string) (*Config, error) {
	// Use default path if not specified
	if cookieFile == "" {
		cookieFile = "cookies.txt"
	}

	cfg := &Config{
		CookieFile:   cookieFile,
		OutputFile:   "emails.txt", // Default output file
		NoOutputFile: false,
	}

	// Load cookies from file
	if err := cfg.LoadCookies(); err != nil {
		return nil, err
	}

	// Validate cookie format
	if err := cfg.ValidateCookies(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadCookies reads the cookie file, skips comment lines, and strips whitespace
// Comment lines start with "//"
// Strips CR (ASCII 13), LF (ASCII 10), space (ASCII 32), and tab (ASCII 9) characters
func (c *Config) LoadCookies() error {
	// Get absolute path for better error messages
	absPath, err := filepath.Abs(c.CookieFile)
	if err != nil {
		absPath = c.CookieFile // Fallback to relative path
	}

	// Check if file exists
	if _, err := os.Stat(c.CookieFile); os.IsNotExist(err) {
		return fmt.Errorf("Cookie file not found at path: %s", absPath)
	}

	// Read file contents
	data, err := os.ReadFile(c.CookieFile)
	if err != nil {
		// Check for permission errors
		if os.IsPermission(err) {
			return fmt.Errorf("Permission denied: cannot read cookie file at path: %s", absPath)
		}
		return fmt.Errorf("Failed to read cookie file: %w", err)
	}

	// Parse file contents line by line
	lines := strings.Split(string(data), "\n")
	var cookieParts []string

	for _, line := range lines {
		// Strip whitespace from line
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip comment lines starting with "//"
		if strings.HasPrefix(line, "//") {
			continue
		}

		// Add non-comment line to cookie parts
		cookieParts = append(cookieParts, line)
	}

	// Join all non-comment lines
	cookieString := strings.Join(cookieParts, "")

	// Strip CR, LF, space, and tab characters from the cookie string
	cookieString = sanitizeCookieString(cookieString)

	// Check if cookie string is empty after processing
	if cookieString == "" {
		return fmt.Errorf("Cookie file is empty or contains only comments")
	}

	c.Cookies = cookieString
	return nil
}

// ValidateCookies checks that the cookie string contains valid KEY=VALUE pairs
// Expected format: KEY1=VALUE1;KEY2=VALUE2 or KEY1=VALUE1 KEY2=VALUE2
func (c *Config) ValidateCookies() error {
	// Check if cookies contain at least one KEY=VALUE pair
	if !strings.Contains(c.Cookies, "=") {
		return fmt.Errorf("Invalid cookie format: expected KEY=VALUE pairs separated by semicolons")
	}

	// Split by semicolon or whitespace to find pairs
	// Replace semicolons with spaces for uniform splitting
	normalized := strings.ReplaceAll(c.Cookies, ";", " ")
	parts := strings.Fields(normalized)

	// Check if at least one part contains "="
	hasValidPair := false
	for _, part := range parts {
		if strings.Contains(part, "=") {
			hasValidPair = true
			break
		}
	}

	if !hasValidPair {
		return fmt.Errorf("Invalid cookie format: expected KEY=VALUE pairs separated by semicolons")
	}

	return nil
}

// sanitizeCookieString removes CR, LF, space, and tab characters from the cookie string
// CR = ASCII 13 (\r), LF = ASCII 10 (\n), Space = ASCII 32, Tab = ASCII 9 (\t)
func sanitizeCookieString(s string) string {
	// Remove CR (carriage return)
	s = strings.ReplaceAll(s, "\r", "")
	// Remove LF (line feed)
	s = strings.ReplaceAll(s, "\n", "")
	// Remove space characters
	s = strings.ReplaceAll(s, " ", "")
	// Remove tab characters
	s = strings.ReplaceAll(s, "\t", "")

	return s
}
