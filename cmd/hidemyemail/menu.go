package hidemyemail

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yourusername/hidemyemail-generator/internal/api"
	"github.com/yourusername/hidemyemail-generator/internal/config"
	"github.com/yourusername/hidemyemail-generator/internal/generator"
	"github.com/yourusername/hidemyemail-generator/internal/lister"
	"github.com/yourusername/hidemyemail-generator/internal/output"
)

func showMenu(ctx context.Context) error {
	reader := bufio.NewReader(os.Stdin)
	formatter := output.NewFormatter()

	// Load config once
	cfg, err := config.LoadConfig("cookies.txt")
	if err != nil {
		formatter.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		return err
	}

	apiClient, err := api.NewClient(cfg.Cookies)
	if err != nil {
		formatter.Error(fmt.Sprintf("Failed to initialize API client: %v", err))
		return err
	}

	for {
		// Check if context was cancelled (Ctrl+C)
		select {
		case <-ctx.Done():
			fmt.Println("\nExiting...")
			return nil
		default:
		}
		
		printMenuOptions()
		
		fmt.Print("Choose option: ")
		choice, err := reader.ReadString('\n')
		if err != nil {
			// Handle EOF or other read errors (e.g., Ctrl+C during input)
			fmt.Println("\nExiting...")
			return nil
		}
		choice = strings.TrimSpace(choice)
		
		switch choice {
		case "1":
			if err := handleGenerate(ctx, reader, apiClient, formatter, cfg); err != nil {
				formatter.Error(fmt.Sprintf("Error: %v", err))
			}
		case "2":
			if err := handleList(ctx, apiClient, formatter); err != nil {
				formatter.Error(fmt.Sprintf("Error: %v", err))
			}
		case "3":
			fmt.Println("\nGoodbye!")
			return nil
		case "":
			// Empty input, just show menu again
			continue
		default:
			formatter.Error("Invalid choice. Please enter 1, 2, or 3")
		}
		
		fmt.Println()
	}
}

func printMenuOptions() {
	fmt.Println("╔═══════════════════════════════════════╗")
	fmt.Println("║              Menu                     ║")
	fmt.Println("╠═══════════════════════════════════════╣")
	fmt.Println("║  1. Generate emails                   ║")
	fmt.Println("║  2. List all emails                   ║")
	fmt.Println("║  3. Exit                              ║")
	fmt.Println("╚═══════════════════════════════════════╝")
	fmt.Println()
}

func handleGenerate(ctx context.Context, reader *bufio.Reader, apiClient *api.Client, formatter *output.Formatter, cfg *config.Config) error {
	fmt.Print("\nEnter label: ")
	label, _ := reader.ReadString('\n')
	label = strings.TrimSpace(label)
	
	if label == "" {
		return fmt.Errorf("label cannot be empty")
	}
	
	fmt.Print("Enter count (1-100): ")
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return fmt.Errorf("invalid count: %v", err)
	}
	
	if count < 1 || count > 100 {
		return fmt.Errorf("count must be between 1 and 100")
	}
	
	// Ask about auto-numbering
	autoNumber := false
	if count > 1 {
		fmt.Print("Add number to label? (y/n): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))
		autoNumber = choice == "y" || choice == "yes"
	}
	
	// Create output directory
	outputDir := "generated"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", outputDir, err)
	}
	
	// Set output file in generated/ directory
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	cfg.OutputFile = fmt.Sprintf("%s/emails_%s_%s.txt", outputDir, label, timestamp)
	cfg.NoOutputFile = false
	
	fmt.Println()
	formatter.Log(fmt.Sprintf("Generating %d email(s) with label '%s'...", count, label), output.INFO)
	
	gen := generator.NewGenerator(apiClient, formatter, cfg)
	allEmails := make([]string, 0, count)
	
	// Generate emails sequentially (like Python version)
	for i := 0; i < count; i++ {
		currentLabel := label
		if autoNumber && count > 1 {
			currentLabel = fmt.Sprintf("%s%d", label, i+1)
		}
		
		emails, err := gen.Generate(currentLabel, 1)
		if err != nil {
			formatter.Error(fmt.Sprintf("Error generating #%d: %v", i+1, err))
			continue
		}
		
		if len(emails) > 0 {
			allEmails = append(allEmails, emails...)
			formatter.Success(fmt.Sprintf("[%d/%d] %s (label: %s)", i+1, count, emails[0], currentLabel))
		}
		
		// Small delay between requests (like Python async behavior)
		if i < count-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}
	
	if len(allEmails) > 0 {
		fmt.Println()
		formatter.Success(fmt.Sprintf("Successfully generated %d email(s)", len(allEmails)))
		fmt.Println()
		for _, email := range allEmails {
			fmt.Println("  " + email)
		}
		fmt.Println()
		formatter.Success(fmt.Sprintf("Saved to %s", cfg.OutputFile))
	}
	
	return nil
}

func handleList(ctx context.Context, apiClient *api.Client, formatter *output.Formatter) error {
	listerInstance := lister.NewLister(apiClient, formatter)
	return listerInstance.List("", false, false) // Show all emails (both active and inactive)
}
