package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	promptgenerator "github.com/szymonsitko/promptgenerator/helpers"
	"github.com/szymonsitko/promptgenerator/helpers/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	APIKey        string  `env:"GEMINI_API_KEY,required"`
	GeminiBaseUrl string  `env:"GEMINI_BASE_URL,required"`
	Temperature   float64 `env:"AI_TEMPERATURE" envDefault:"1.0"`
	TopP          float64 `env:"AI_TOP_P" envDefault:"0.8"`
	MaxTokens     int     `env:"AI_MAX_TOKENS" envDefault:"4096"`
	NumResults    int     `env:"AI_NUM_RESULTS" envDefault:"10"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v.  Using environment variables directly.", err)
	}

	var inputFile, outputFile, dbFile, actor, prompt string
	var documentation, explanations, comments bool

	var rootCmd = &cobra.Command{
		Use:   "cli-tool",
		Short: "A CLI tool that processes optional input and output files with actor and prompt required option(s)",
		Run: func(cmd *cobra.Command, args []string) {
			// Load configuration from environment variables
			cfg := Config{}
			if err := env.Parse(&cfg); err != nil {
				log.Fatalf("Failed to load environment variables: %v", err)
			}

			if actor == "" || prompt == "" {
				fmt.Println("Error: Both --actor and --prompt parameters are required.")
				err = cmd.Help()
				if err != nil {
					log.Fatalf("Error: Cannot get command helper %s", err)
				}
				os.Exit(1)
			}

			// Check if the database file exists, create one if not
			if dbFile != "" {
				if _, err := os.Stat(dbFile); os.IsNotExist(err) {
					file, err := os.Create(dbFile)
					if err != nil {
						log.Fatalf("Failed to create database file: %s", err)
					}
					file.Close()
					fmt.Println("Database file created.")
				}
			}

			var content string
			if inputFile != "" {
				data, err := os.ReadFile(inputFile)
				if err != nil {
					log.Fatalf("Error reading input file: %v", err)
				}
				content = string(data)
			}
			if !strings.HasSuffix(cfg.GeminiBaseUrl, "?key=") {
				cfg.GeminiBaseUrl += "?key="
			}

			// Use the loaded configuration here
			aiGenerator := promptgenerator.NewPromptHandler(
				cfg.GeminiBaseUrl,
				cfg.APIKey,
				cfg.Temperature,
				cfg.TopP,
				cfg.MaxTokens,
				cfg.NumResults,
			)

			// Attach content to a prompt
			if content != "" {
				prompt = fmt.Sprintf("%s.\nContent: \n\n%s", prompt, content)
			}

			// Actual prompt
			result, err := aiGenerator.GenerateContent(actor, prompt, promptgenerator.Modes{
				Documentation: documentation,
				Explanations:  explanations,
				Comments:      comments,
			})
			if err != nil {
				log.Fatal(err)
			}

			// Store prompt into database if provided
			if dbFile != "" {
				// Open SQLite database
				db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
				if err != nil {
					log.Fatalf("Failed to connect to database: %s", err)
					return
				}

				// Perform migration
				err = db.AutoMigrate(&database.Prompt{})
				if err != nil {
					fmt.Println("Migration failed:", err)
					return
				}

				repo := database.NewPromptRepository(db)

				// Insert a sample prompt
				newPrompt := database.Prompt{
					Prompt:        prompt,
					Content:       content,
					Actor:         actor,
					Comments:      comments,
					Documentation: documentation,
					Explanations:  explanations,
				}
				err = repo.CreatePrompt(newPrompt)
				if err != nil {
					log.Fatalf("Failed to insert prompt: %s", err)
				}
			}

			// Clean up result text
			result = cleanResultText(result)

			if outputFile != "" {
				writeToFile(outputFile, result)
			} else {
				fmt.Println(result)
			}
		},
	}

	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Optional input file to load data and attach it to the prompt")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Optional output file, if provided, prompt result will be saved to given path")
	rootCmd.Flags().StringVarP(&dbFile, "file", "f", "", "Optional database path, if provided & database exists, prompts will be persisted. If provided & db does not exist, it will be created locally (sqlite)")
	rootCmd.Flags().BoolVarP(&documentation, "documentation", "d", false, "Optional documentation param, if provided, request to generate documentation will be included in prompt")
	rootCmd.Flags().BoolVarP(&explanations, "explanations", "e", false, "Optional explanatiions param, if provided, request to generate code explanations will be included in prompt")
	rootCmd.Flags().BoolVarP(&comments, "comments", "c", false, "Optional comments param, if provided, request to generate extensive in-code comments will be included in prompt")
	rootCmd.Flags().StringVarP(&actor, "actor", "a", "", "Required prompt actor name to define prompt persona")
	rootCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "Required prompt content")

	if err := rootCmd.MarkFlagRequired("actor"); err != nil {
		panic(err)
	}
	if err := rootCmd.MarkFlagRequired("prompt"); err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

func writeToFile(filename, content string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		log.Fatalf("Error writing to output file: %v", err)
	}
	writer.Flush()
}

func cleanResultText(text string) string {
	re := regexp.MustCompile("(?m)^```go$|^```$|^\\s*$")
	return re.ReplaceAllString(text, "")
}
