package lister

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yourusername/hidemyemail-generator/internal/api"
	"github.com/yourusername/hidemyemail-generator/internal/output"
	"github.com/yourusername/hidemyemail-generator/pkg/models"
)

type Lister struct {
	client    *api.Client
	formatter *output.Formatter
}

func NewLister(client *api.Client, formatter *output.Formatter) *Lister {
	return &Lister{
		client:    client,
		formatter: formatter,
	}
}

func (l *Lister) List(labelQuery string, activeOnly bool, inactiveOnly bool) error {
	listResp, err := l.client.ListEmails()
	if err != nil {
		if isNetworkError(err) {
			l.formatter.Error("Network error: failed to connect to iCloud API")
			return fmt.Errorf("network error: failed to connect to iCloud API")
		}
		if isTimeoutError(err) {
			l.formatter.Error("Request timeout: iCloud API did not respond within 10 seconds")
			return fmt.Errorf("request timeout: iCloud API did not respond within 10 seconds")
		}
		return err
	}

	if !listResp.Success {
		if errorMsg, ok := listResp.Error.(string); ok {
			l.formatter.Error(errorMsg)
		} else {
			l.formatter.Error("Failed to list emails: unknown error")
		}
		return nil
	}

	filteredEmails, err := l.filterEmails(listResp.Result.HMEEmails, labelQuery, activeOnly, inactiveOnly)
	if err != nil {
		l.formatter.Error(fmt.Sprintf("Invalid regex pattern: %v", err))
		return fmt.Errorf("invalid regex pattern: %w", err)
	}

	if len(filteredEmails) == 0 {
		l.formatter.Log("No emails found matching the specified filters", output.INFO)
		return nil
	}

	// Sort by creation date (newest first)
	sort.Slice(filteredEmails, func(i, j int) bool {
		return filteredEmails[i].CreateTimestamp > filteredEmails[j].CreateTimestamp
	})

	l.displayTable(filteredEmails)
	
	// Save to files
	if err := l.saveToFiles(filteredEmails); err != nil {
		l.formatter.Error(fmt.Sprintf("Failed to save to files: %v", err))
	} else {
		l.formatter.Success("Saved to results/emails_list.txt and results/emails_full.txt")
	}
	
	return nil
}

func (l *Lister) displayTable(emails []models.Email) {
	headers := []string{"Email", "Label", "Note", "IsActive", "Created"}
	rows := make([][]string, 0, len(emails))
	
	for _, email := range emails {
		isActiveStr := "false"
		if email.IsActive {
			isActiveStr = "true"
		}
		
		// Convert timestamp (milliseconds) to human-readable format
		createdTime := time.Unix(email.CreateTimestamp/1000, 0)
		createdStr := createdTime.Format("2006-01-02 15:04:05")
		
		row := []string{
			email.HME,
			email.Label,
			email.Note,
			isActiveStr,
			createdStr,
		}
		rows = append(rows, row)
	}
	
	l.formatter.Table(headers, rows)
}

// saveToFiles saves emails to two files: emails_list.txt (just emails) and emails_full.txt (full info)
func (l *Lister) saveToFiles(emails []models.Email) error {
	// Create results directory
	outputDir := "results"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}
	
	// File 1: Just email addresses
	emailsList := make([]string, 0, len(emails))
	for _, email := range emails {
		emailsList = append(emailsList, email.HME)
	}
	
	listPath := fmt.Sprintf("%s/emails_list.txt", outputDir)
	err := os.WriteFile(listPath, []byte(strings.Join(emailsList, "\n")+"\n"), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", listPath, err)
	}
	
	// File 2: Full information with aligned columns
	var fullInfo strings.Builder
	
	// Calculate max widths for each column
	maxEmailLen := len("Email")
	maxLabelLen := len("Label")
	maxNoteLen := len("Note")
	maxActiveLen := len("IsActive")
	maxCreatedLen := len("Created")
	
	for _, email := range emails {
		if len(email.HME) > maxEmailLen {
			maxEmailLen = len(email.HME)
		}
		if len(email.Label) > maxLabelLen {
			maxLabelLen = len(email.Label)
		}
		if len(email.Note) > maxNoteLen {
			maxNoteLen = len(email.Note)
		}
	}
	
	// Write header
	headerFormat := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n", 
		maxEmailLen, maxLabelLen, maxNoteLen, maxActiveLen, maxCreatedLen)
	fullInfo.WriteString(fmt.Sprintf(headerFormat, "Email", "Label", "Note", "IsActive", "Created"))
	
	// Write separator line
	totalWidth := maxEmailLen + maxLabelLen + maxNoteLen + maxActiveLen + maxCreatedLen + 12 // 12 for " | " separators
	fullInfo.WriteString(strings.Repeat("-", totalWidth) + "\n")
	
	// Write data rows
	rowFormat := fmt.Sprintf("%%-%ds | %%-%ds | %%-%ds | %%-%ds | %%-%ds\n", 
		maxEmailLen, maxLabelLen, maxNoteLen, maxActiveLen, maxCreatedLen)
	
	for _, email := range emails {
		createdTime := time.Unix(email.CreateTimestamp/1000, 0)
		createdStr := createdTime.Format("2006-01-02 15:04:05")
		isActiveStr := "false"
		if email.IsActive {
			isActiveStr = "true"
		}
		
		fullInfo.WriteString(fmt.Sprintf(rowFormat,
			email.HME,
			email.Label,
			email.Note,
			isActiveStr,
			createdStr,
		))
	}
	
	fullPath := fmt.Sprintf("%s/emails_full.txt", outputDir)
	err = os.WriteFile(fullPath, []byte(fullInfo.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", fullPath, err)
	}
	
	return nil
}

func (l *Lister) filterEmails(emails []models.Email, labelQuery string, activeOnly bool, inactiveOnly bool) ([]models.Email, error) {
	var labelRegex *regexp.Regexp
	var err error
	if labelQuery != "" {
		labelRegex, err = regexp.Compile(labelQuery)
		if err != nil {
			return nil, err
		}
	}

	filtered := make([]models.Email, 0, len(emails))
	for _, email := range emails {
		if labelRegex != nil && !labelRegex.MatchString(email.Label) {
			continue
		}
		if activeOnly && !email.IsActive {
			continue
		}
		if inactiveOnly && email.IsActive {
			continue
		}
		filtered = append(filtered, email)
	}

	return filtered, nil
}

func isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "no such host") ||
		strings.Contains(errMsg, "no route to host") ||
		strings.Contains(errMsg, "network is unreachable") ||
		strings.Contains(errMsg, "connection reset")
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "deadline exceeded") ||
		strings.Contains(errMsg, "timed out")
}
