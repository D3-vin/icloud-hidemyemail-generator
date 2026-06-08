package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yourusername/hidemyemail-generator/internal/config"
	"github.com/yourusername/hidemyemail-generator/internal/output"
	"github.com/yourusername/hidemyemail-generator/pkg/models"
)

type APIClient interface {
	GenerateEmail() (*models.GenerateResponse, error)
	ReserveEmail(hme, label, note string) (*models.ReserveResponse, error)
	ListEmails() (*models.ListResponse, error)
}

type Generator struct {
	client    APIClient
	formatter *output.Formatter
	config    *config.Config
}

func NewGenerator(client APIClient, formatter *output.Formatter, config *config.Config) *Generator {
	return &Generator{
		client:    client,
		formatter: formatter,
		config:    config,
	}
}

// Generate creates email addresses sequentially (like Python version)
func (g *Generator) Generate(label string, count int) ([]string, error) {
	if count < 1 || count > 100 {
		return nil, fmt.Errorf("Error: --count must be between 1 and 100")
	}

	var successfulEmails []string
	var errors []error

	// Generate sequentially (simpler, like Python async)
	for i := 1; i <= count; i++ {
		email, err := g.generateOne(label, i, count)
		if err != nil {
			errors = append(errors, err)
			g.formatter.Error(err.Error())
			continue
		}
		successfulEmails = append(successfulEmails, email)
	}

	summary := GenerationSummary{
		TotalRequested:  count,
		SuccessfulCount: len(successfulEmails),
		FailedCount:     len(errors),
		Errors:          errors,
	}

	g.displaySummary(summary)

	if !g.config.NoOutputFile && len(successfulEmails) > 0 {
		if err := g.writeToFile(successfulEmails); err != nil {
			g.formatter.Error(fmt.Sprintf("Failed to write to file: %s", err.Error()))
		}
	}

	return successfulEmails, nil
}

func (g *Generator) generateOne(label string, current, total int) (string, error) {
	genResp, err := g.client.GenerateEmail()
	if err != nil {
		return "", fmt.Errorf("failed to generate email %d of %d: %w", current, total, err)
	}

	if !genResp.Success {
		errorMsg := "unknown error"
		if genResp.Error != nil {
			errorMsg = fmt.Sprintf("%v", genResp.Error)
		}
		return "", fmt.Errorf("generation failed for email %d of %d: %s", current, total, errorMsg)
	}

	generatedEmail := genResp.Result.HME
	if generatedEmail == "" {
		return "", fmt.Errorf("generation succeeded but no email address returned")
	}

	reserveResp, err := g.client.ReserveEmail(generatedEmail, label, "")
	if err != nil {
		return "", fmt.Errorf("failed to reserve email %d of %d (%s): %w", current, total, generatedEmail, err)
	}

	if !reserveResp.Success {
		errorMsg := "unknown error"
		if reserveResp.Error != nil {
			errorMsg = fmt.Sprintf("%v", reserveResp.Error)
		}
		return "", fmt.Errorf("reservation failed for email %d of %d (%s): %s", current, total, generatedEmail, errorMsg)
	}

	g.formatter.Progress(current, total, "Generating emails")
	return generatedEmail, nil
}

func (g *Generator) writeToFile(emails []string) error {
	outputDir := filepath.Dir(g.config.OutputFile)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("Failed to write to file %s: permission denied", g.config.OutputFile)
			}
			return fmt.Errorf("Failed to create parent directories for %s: %w", g.config.OutputFile, err)
		}
	}

	file, err := os.OpenFile(g.config.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("Failed to write to file %s: permission denied", g.config.OutputFile)
		}
		if strings.Contains(strings.ToLower(err.Error()), "no space left") {
			return fmt.Errorf("Failed to write to file %s: disk full", g.config.OutputFile)
		}
		return fmt.Errorf("Failed to open file %s: %w", g.config.OutputFile, err)
	}

	defer func() {
		file.Sync()
		file.Close()
	}()

	lineEnding := "\n"
	if runtime.GOOS == "windows" {
		lineEnding = "\r\n"
	}

	for _, email := range emails {
		_, err := file.WriteString(email + lineEnding)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("Failed to write to file %s: permission denied", g.config.OutputFile)
			}
			if strings.Contains(strings.ToLower(err.Error()), "no space left") {
				return fmt.Errorf("Failed to write to file %s: disk full", g.config.OutputFile)
			}
			return fmt.Errorf("Failed to write email to file %s: %w", g.config.OutputFile, err)
		}
	}

	return nil
}

type GenerationSummary struct {
	TotalRequested  int
	SuccessfulCount int
	FailedCount     int
	Errors          []error
}

func (g *Generator) displaySummary(summary GenerationSummary) {
	if summary.FailedCount == 0 {
		g.formatter.Success(fmt.Sprintf("Successfully generated %d/%d emails",
			summary.SuccessfulCount, summary.TotalRequested))
	} else if summary.SuccessfulCount == 0 {
		g.formatter.Error(fmt.Sprintf("Failed to generate all emails (0/%d successful)",
			summary.TotalRequested))
	} else {
		g.formatter.Log(fmt.Sprintf("Generated %d/%d emails (%d failed)",
			summary.SuccessfulCount, summary.TotalRequested, summary.FailedCount),
			output.WARNING)
	}
}
