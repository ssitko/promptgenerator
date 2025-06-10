
This document provides a guide on how to use the `cli-tool` application.  This tool leverages the Gemini AI model to generate content based on a provided prompt and actor.  It supports input from files, output to files, and optional persistence of prompts to a SQLite database.

## Prerequisites

*   **Go:** Ensure you have Go installed on your system.  You can download it from [https://go.dev/dl/](https://go.dev/dl/).
*   **Environment Variables:** The tool requires the following environment variables to be set:
    *   `GEMINI_API_KEY`: Your Gemini API key.
    *   `GEMINI_BASE_URL`: The base URL for the Gemini API.
    *   `AI_TEMPERATURE` (optional): The temperature parameter for the AI model (default: 1.0).
    *   `AI_TOP_P` (optional): The top-p parameter for the AI model (default: 0.8).
    *   `AI_MAX_TOKENS` (optional): The maximum number of tokens for the AI model (default: 4096).
    *   `AI_NUM_RESULTS` (optional): The number of results to return from the AI model (default: 10).

*   **`.env` file (optional):** You can store the environment variables in a `.env` file in the same directory as the executable. The tool will automatically load these variables.  If a `.env` file exists, the tool will attempt to load it; otherwise, it will use the system environment variables directly.

## Installation

1.  **Clone the repository (if you have the source code):**

    ```bash
    git clone <repository_url>
    cd <repository_directory>
    ```

2.  **Build the executable:**

    ```bash
    go build -o cli-tool .
    ```

## Usage

The `cli-tool` is a command-line application with the following options:


cli-tool [flags]


### Flags

*   `-i, --input string`:  Optional input file to load data and attach it to the prompt.  The content of this file will be appended to the prompt.
*   `-o, --output string`: Optional output file. If provided, the generated content will be saved to this file. If not provided, the output will be printed to the console.
*   `-f, --file string`: Optional database path. If provided and the database exists, prompts will be persisted. If provided and the database does not exist, it will be created locally (SQLite).
*   `-d, --documentation`:  If provided, the prompt will include a request to generate documentation.
*   `-e, --explanations`: If provided, the prompt will include a request to generate code explanations.
*   `-c, --comments`:   If provided, the prompt will include a request to generate extensive in-code comments.
*   `-a, --actor string`:  **Required.** The actor name to define the prompt persona. This helps the AI understand the context of the request (e.g., "Senior Golang Developer", "Technical Writer").
*   `-p, --prompt string`: **Required.** The main content of the prompt. This is the instruction you want the AI to follow.

### Examples

1.  **Basic Usage (Console Output):**

    ```bash
    ./cli-tool --actor "Senior Golang Developer" --prompt "Write a function to reverse a string."
    ```

    This will send a request to the Gemini API with the specified actor and prompt, and print the generated code to the console.  Make sure you have `GEMINI_API_KEY` and `GEMINI_BASE_URL` set as environment variables.

2.  **Using Input and Output Files:**

    ```bash
    ./cli-tool --actor "Senior Golang Developer" --prompt "Refactor this code for better readability and performance." --input input.go --output output.go
    ```

    This will read the content of `input.go`, append it to the prompt, send the request to the Gemini API, and save the generated code to `output.go`.

3.  **Saving Prompts to a Database:**

    ```bash
    ./cli-tool --actor "Senior Golang Developer" --prompt "Write a unit test for this function." --input input.go --file prompts.db
    ```

    This will read the content of `input.go`, append it to the prompt, send the request to the Gemini API, and save the prompt (along with the actor and input file content) to the `prompts.db` SQLite database. If `prompts.db` doesn't exist, it will be created.

4.  **Generating Documentation, Explanations, and Comments:**

    ```bash
    ./cli-tool --actor "Senior Golang Developer" --prompt "Implement a simple HTTP server." --documentation --explanations --comments --output server.go
    ```

    This will send a request to the Gemini API with the specified actor and prompt, and request the AI to generate documentation, explanations, and comments within the generated code.  The result will be saved to `server.go`.

5.  **Using a `.env` file:**

    Create a `.env` file in the same directory as the `cli-tool` executable with the following content (replace with your actual API key and base URL):

    ```
    GEMINI_API_KEY=YOUR_GEMINI_API_KEY
    GEMINI_BASE_URL=YOUR_GEMINI_BASE_URL
    ```

    Then, you can run the tool without explicitly setting the environment variables in your shell:

    ```bash
    ./cli-tool --actor "Senior Golang Developer" --prompt "Write a function to calculate the factorial of a number."
    ```

## Database Schema

If the `--file` flag is used, the tool will create or connect to a SQLite database with the following schema for the `prompts` table:

| Column        | Type    | Description                                                                                                                                       |
| ------------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`          | INTEGER | Primary key, auto-incrementing.                                                                                                                  |
| `prompt`      | TEXT    | The main prompt provided to the AI.                                                                                                             |
| `content`     | TEXT    | The content of the input file (if provided).                                                                                                    |
| `actor`       | TEXT    | The actor name used in the prompt.                                                                                                                |
| `comments`    | BOOLEAN | Indicates whether the prompt included a request for comments.                                                                                     |
| `documentation` | BOOLEAN | Indicates whether the prompt included a request for documentation.                                                                                |
| `explanations`  | BOOLEAN | Indicates whether the prompt included a request for explanations.                                                                                   |
| `created_at`  | DATETIME| Timestamp indicating when the prompt was created.                                                                                               |
| `updated_at`  | DATETIME| Timestamp indicating when the prompt was last updated.                                                                                               |

## Error Handling

*   The tool will exit with an error message if the required `--actor` or `--prompt` flags are not provided.
*   The tool will print an error message and exit if it fails to read the input file or create the output file.
*   The tool will print an error message and exit if it fails to connect to the database or insert the prompt.
*   The tool will print an error message and exit if it fails to load environment variables.
*   The tool will print an error message and exit if the Gemini API returns an error.

## Notes

*   Ensure that your `GEMINI_BASE_URL` ends with `?key=` if that is how your Gemini API requires the key.
*   The tool uses the `github.com/caarlos0/env/v10` library to load environment variables, which supports various configuration options.  Refer to the library's documentation for more details.
*   The `cleanResultText` function removes unnecessary "```go" and "```" markers from the generated code to improve readability.
