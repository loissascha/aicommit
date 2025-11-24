# AICommit

AICommit is a command-line tool that uses Google's Gemini API to automatically generate commit messages for your staged changes. It analyzes the output of `git diff --staged` and creates a descriptive commit message, helping you write better commits faster.

## Features

- **Automatic Commit Messages:** Generates commit messages from your staged changes.
- **Powered by Gemini:** Uses Google's Gemini API for high-quality commit messages.
- **Easy to Use:** Simply run `aicommit` in your terminal to generate and commit your changes.
- **Token Usage:** Optionally display the number of tokens used for the API call.

## Configuration

 **Add your API key:**

Create an environment variable with your gemini api key

```
GEMINI_API_KEY=your_api_key
```

You can get your API key from [Google AI Studio](https://aistudio.google.com/app/apikey).

## Usage

1.  **Stage your changes:**

    ```bash
    git add .
    ```

2.  **Run AICommit:**

    ```bash
    aicommit
    ```

    The tool will generate a commit message and commit your changes.


## Development

### Building from source

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/loissascha/aicommit.git
    cd aicommit
    ```

2.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Run the application:**

    ```bash
    go run cmd/aicommit/
    ```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
