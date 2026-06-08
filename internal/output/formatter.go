package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// LogLevel defines the log message severity level
type LogLevel int

const (
	INFO LogLevel = iota
	ERROR
	SUCCESS
	WARNING
)

// Formatter provides rich console output with progress indicators and tables
type Formatter struct {
	progressBar   *pterm.ProgressbarPrinter
	lastUpdate    time.Time
	updateThrottle time.Duration
}

// NewFormatter creates and initializes a new Formatter instance
func NewFormatter() *Formatter {
	return &Formatter{
		updateThrottle: 500 * time.Millisecond,
		lastUpdate:     time.Time{},
	}
}

// Log displays a message with colored output based on the level
// Green for success, red for error, cyan for info, yellow for warning
func (f *Formatter) Log(message string, level LogLevel) {
	switch level {
	case SUCCESS:
		pterm.Success.Println(message)
	case ERROR:
		pterm.Error.Println(message)
	case INFO:
		pterm.Info.Println(message)
	case WARNING:
		pterm.Warning.Println(message)
	}
}

// Progress displays a visual progress indicator that updates every 500ms
// Shows current/total count with a message
func (f *Formatter) Progress(current, total int, message string) {
	now := time.Now()
	
	// Initialize progress bar on first call or if it's nil
	if f.progressBar == nil {
		f.progressBar, _ = pterm.DefaultProgressbar.
			WithTotal(total).
			WithTitle(message).
			Start()
		f.lastUpdate = now
	}
	
	// Throttle updates to every 500ms minimum
	if now.Sub(f.lastUpdate) >= f.updateThrottle || current == total {
		f.progressBar.UpdateTitle(fmt.Sprintf("%s [%d/%d]", message, current, total))
		f.progressBar.Current = current
		f.lastUpdate = now
	}
	
	// Stop and cleanup progress bar when complete
	if current >= total {
		f.progressBar.Stop()
		f.progressBar = nil
	}
}

// Table renders data as a formatted ASCII table with headers and rows
// Columns are automatically wrapped at max 80 characters
func (f *Formatter) Table(headers []string, rows [][]string) {
	// Prepare table data with headers
	tableData := pterm.TableData{headers}
	
	// Truncate cell content to max 80 characters per column
	for _, row := range rows {
		truncatedRow := make([]string, len(row))
		for i, cell := range row {
			if len(cell) > 80 {
				truncatedRow[i] = cell[:77] + "..."
			} else {
				truncatedRow[i] = cell
			}
		}
		tableData = append(tableData, truncatedRow)
	}
	
	// Render the table
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

// Success displays a success message in green with a checkmark prefix
func (f *Formatter) Success(message string) {
	pterm.Success.WithPrefix(pterm.Prefix{
		Text:  "✓",
		Style: pterm.NewStyle(pterm.FgGreen),
	}).Println(message)
}

// Error displays an error message in red with a cross prefix
func (f *Formatter) Error(message string) {
	pterm.Error.WithPrefix(pterm.Prefix{
		Text:  "✗",
		Style: pterm.NewStyle(pterm.FgRed),
	}).Println(message)
}

// Rule displays a visual separator (horizontal line) between output sections
func (f *Formatter) Rule() {
	// Create a horizontal line using pterm's section style
	width := 80 // Default terminal width
	pterm.DefaultSection.WithLevel(1).Println(strings.Repeat("─", width))
}
