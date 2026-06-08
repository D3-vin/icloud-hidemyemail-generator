package hidemyemail

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/yourusername/hidemyemail-generator/internal/api"
	"github.com/yourusername/hidemyemail-generator/internal/config"
	"github.com/yourusername/hidemyemail-generator/internal/generator"
	"github.com/yourusername/hidemyemail-generator/internal/lister"
	"github.com/yourusername/hidemyemail-generator/internal/output"
)

const (
	cyan    = "\033[36m"
	yellow  = "\033[33m"
	magenta = "\033[35m"
	green   = "\033[32m"
	reset   = "\033[0m"
	bold    = "\033[1m"
)

var (
	cookieFile     string
	generateLabel  string
	generateCount  int
	generateOutput string
	noOutputFile   bool
	listLabelQuery string
	listActive     bool
	listInactive   bool
)

func Execute() {
	printBanner()
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\nExiting...")
		cancel()
		os.Exit(0)
	}()

	if err := newRootCmd(ctx).Execute(); err != nil {
		os.Exit(1)
	}
}

func printBanner() {
	fmt.Println()
	fmt.Println(cyan + "  ██╗ ██████╗██╗      ██████╗ ██╗   ██╗██████╗      ██╗  ██╗███╗   ███╗███████╗" + reset)
	fmt.Println(cyan + "  ██║██╔════╝██║     ██╔═══██╗██║   ██║██╔══██╗     ██║  ██║████╗ ████║██╔════╝" + reset)
	fmt.Println(cyan + "  ██║██║     ██║     ██║   ██║██║   ██║██║  ██║     ███████║██╔████╔██║█████╗  " + reset)
	fmt.Println(cyan + "  ██║██║     ██║     ██║   ██║██║   ██║██║  ██║     ██╔══██║██║╚██╔╝██║██╔══╝  " + reset)
	fmt.Println(cyan + "  ██║╚██████╗███████╗╚██████╔╝╚██████╔╝██████╔╝     ██║  ██║██║ ╚═╝ ██║███████╗" + reset)
	fmt.Println(cyan + "  ╚═╝ ╚═════╝╚══════╝ ╚═════╝  ╚═════╝ ╚═════╝      ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝" + reset)
	fmt.Println()
	fmt.Println(yellow + "  📧 Telegram:" + reset + " https://t.me/D3_vin")
	fmt.Println(magenta + "  👤 Author:" + reset + " @D3vin_dev")
	fmt.Println(green + "  🔗 GitHub:" + reset + " https://github.com/D3-vin/icloud-hidemyemail-generator")
	fmt.Println(cyan + "  📦 Version:" + reset + " 1.0.0")
	fmt.Println()
	fmt.Println(bold + "    ━━━ iCloud HideMyEmail Generator ━━━" + reset)
	fmt.Println()
}

func newRootCmd(ctx context.Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hidemyemail",
		Short: "Generate and manage iCloud Hide My Email addresses",
		Long:  `HideMyEmail Generator - A Go CLI application for generating and managing iCloud Hide My Email addresses with TLS fingerprinting.`,
		Run: func(cmd *cobra.Command, args []string) {
			// If no subcommand provided, show interactive menu
			if len(args) == 0 {
				if err := showMenu(ctx); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
			} else {
				cmd.Help()
			}
		},
	}

	rootCmd.AddCommand(newGenerateCmd(ctx))
	rootCmd.AddCommand(newListCmd(ctx))
	return rootCmd
}

func newGenerateCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate new Hide My Email addresses",
		Long:  `Generate one or more new Hide My Email addresses with the specified label.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate(ctx)
		},
	}

	cmd.Flags().StringVarP(&generateLabel, "label", "l", "", "Label for the generated email addresses (required)")
	cmd.Flags().IntVarP(&generateCount, "count", "c", 0, "Number of email addresses to generate (required, 1-100)")
	cmd.Flags().StringVar(&cookieFile, "cookie-file", "cookies.txt", "Path to cookie file")
	cmd.Flags().StringVarP(&generateOutput, "output", "o", "emails.txt", "Output file for generated emails")
	cmd.Flags().BoolVar(&noOutputFile, "no-output-file", false, "Skip writing to output file")

	cmd.MarkFlagRequired("label")
	cmd.MarkFlagRequired("count")
	return cmd
}

func newListCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List existing Hide My Email addresses",
		Long:  `List all existing Hide My Email addresses with optional filtering.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(ctx)
		},
	}

	cmd.Flags().StringVar(&listLabelQuery, "label-query", "", "Regular expression to filter by label")
	cmd.Flags().BoolVar(&listActive, "active", false, "Show only active emails")
	cmd.Flags().BoolVar(&listInactive, "inactive", false, "Show only inactive emails")
	cmd.Flags().StringVar(&cookieFile, "cookie-file", "cookies.txt", "Path to cookie file")
	return cmd
}

func runGenerate(ctx context.Context) error {
	if generateCount < 1 {
		return fmt.Errorf("Error: --count must be a positive integer")
	}
	if generateCount > 100 {
		return fmt.Errorf("Error: --count must be between 1 and 100")
	}

	formatter := output.NewFormatter()

	cfg, err := config.LoadConfig(cookieFile)
	if err != nil {
		formatter.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		return err
	}

	cfg.OutputFile = generateOutput
	cfg.NoOutputFile = noOutputFile

	apiClient, err := api.NewClient(cfg.Cookies)
	if err != nil {
		formatter.Error(fmt.Sprintf("Failed to initialize API client: %v", err))
		return err
	}

	gen := generator.NewGenerator(apiClient, formatter, cfg)

	emails, err := gen.Generate(generateLabel, generateCount)
	if err != nil {
		formatter.Error(fmt.Sprintf("Failed to generate emails: %v", err))
		return err
	}

	if len(emails) > 0 {
		formatter.Success(fmt.Sprintf("Successfully generated %d email(s)", len(emails)))
		for _, email := range emails {
			fmt.Println(email)
		}
	}

	return nil
}

func runList(ctx context.Context) error {
	if listActive && listInactive {
		return fmt.Errorf("Error: --active and --inactive flags are mutually exclusive")
	}

	formatter := output.NewFormatter()

	cfg, err := config.LoadConfig(cookieFile)
	if err != nil {
		formatter.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		return err
	}

	apiClient, err := api.NewClient(cfg.Cookies)
	if err != nil {
		formatter.Error(fmt.Sprintf("Failed to initialize API client: %v", err))
		return err
	}

	listerInstance := lister.NewLister(apiClient, formatter)
	return listerInstance.List(listLabelQuery, listActive, listInactive)
}
