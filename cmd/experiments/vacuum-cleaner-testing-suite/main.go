package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	logLevel string
	baseURL  string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "vacuum-testing-suite",
		Short: "Comprehensive testing suite for VacuumMart e-commerce platform",
		Long: `A comprehensive testing framework for the vacuum cleaner retail store including:
- API testing for product catalog and cart endpoints
- End-to-end shopping flow testing
- Payment integration testing
- Performance and load testing
- Security testing (XSS, injection attacks)
- Mobile responsiveness testing
- User acceptance testing`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setupLogging()
		},
	}

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (trace, debug, info, warn, error)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "http://localhost:8080", "Base URL of the vacuum store API")

	// Add subcommands
	rootCmd.AddCommand(
		newAPITestCmd(),
		newE2ETestCmd(),
		newPerformanceTestCmd(),
		newSecurityTestCmd(),
		newMobileTestCmd(),
		newUATTestCmd(),
		newRunAllTestsCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Command execution failed")
		os.Exit(1)
	}
}

func setupLogging() {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func newAPITestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Run API tests for product catalog and cart endpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("base_url", baseURL).Msg("Running API tests")
			tester := NewAPITester(baseURL)
			return tester.RunTests(context.Background())
		},
	}
}

func newE2ETestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "e2e",
		Short: "Run end-to-end shopping flow tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("base_url", baseURL).Msg("Running E2E tests")
			tester := NewE2ETester(baseURL)
			return tester.RunTests(context.Background())
		},
	}
}

func newPerformanceTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "performance",
		Short: "Run performance and load tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("base_url", baseURL).Msg("Running performance tests")
			tester := NewPerformanceTester(baseURL)
			return tester.RunTests(context.Background())
		},
	}
}

func newSecurityTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "security",
		Short: "Run security tests (XSS, injection attacks)",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("base_url", baseURL).Msg("Running security tests")
			tester := NewSecurityTester(baseURL)
			return tester.RunTests(context.Background())
		},
	}
}

func newMobileTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mobile",
		Short: "Run mobile responsiveness tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("base_url", baseURL).Msg("Running mobile tests")
			tester := NewMobileTester(baseURL)
			return tester.RunTests(context.Background())
		},
	}
}

func newUATTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uat",
		Short: "Run user acceptance tests based on business requirements",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("base_url", baseURL).Msg("Running UAT tests")
			tester := NewUATTester(baseURL)
			return tester.RunTests(context.Background())
		},
	}
}

func newRunAllTestsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Run all test suites",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("base_url", baseURL).Msg("Running all test suites")
			
			ctx := context.Background()
			
			// Run tests in order
			testSuites := []struct {
				name   string
				tester TestRunner
			}{
				{"API Tests", NewAPITester(baseURL)},
				{"E2E Tests", NewE2ETester(baseURL)},
				{"Performance Tests", NewPerformanceTester(baseURL)},
				{"Security Tests", NewSecurityTester(baseURL)},
				{"Mobile Tests", NewMobileTester(baseURL)},
				{"UAT Tests", NewUATTester(baseURL)},
			}
			
			for _, suite := range testSuites {
				log.Info().Str("suite", suite.name).Msg("Starting test suite")
				if err := suite.tester.RunTests(ctx); err != nil {
					log.Error().Err(err).Str("suite", suite.name).Msg("Test suite failed")
					return fmt.Errorf("test suite %s failed: %w", suite.name, err)
				}
				log.Info().Str("suite", suite.name).Msg("Test suite completed successfully")
			}
			
			log.Info().Msg("All test suites completed successfully")
			return nil
		},
	}
}

// TestRunner interface for all test suites
type TestRunner interface {
	RunTests(ctx context.Context) error
}

// Factory functions for creating test runners
func NewAPITester(baseURL string) TestRunner {
	return &APITester{
		baseURL:    baseURL,
		httpClient: createHTTPClient(),
		results:    make([]TestResult, 0),
	}
}

func NewE2ETester(baseURL string) TestRunner {
	return &E2ETester{
		baseURL:    baseURL,
		httpClient: createHTTPClient(),
		scenarios:  make([]TestScenario, 0),
	}
}

func NewPerformanceTester(baseURL string) TestRunner {
	return &PerformanceTester{
		baseURL:    baseURL,
		httpClient: createHTTPClient(),
		results:    make([]LoadTestResult, 0),
	}
}

func NewSecurityTester(baseURL string) TestRunner {
	return &SecurityTester{
		baseURL:    baseURL,
		httpClient: createHTTPClient(),
		results:    make([]SecurityTestResult, 0),
	}
}

func NewMobileTester(baseURL string) TestRunner {
	return &MobileTester{
		baseURL:    baseURL,
		httpClient: createHTTPClient(),
		results:    make([]MobileTestResult, 0),
	}
}

func NewUATTester(baseURL string) TestRunner {
	return &UATTester{
		baseURL:    baseURL,
		httpClient: createHTTPClient(),
		scenarios:  make([]UATScenario, 0),
	}
}

func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}
