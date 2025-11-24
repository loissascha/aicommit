package ai

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

func getClient(ctx context.Context) (*genai.Client, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// returns: header, message, error (if any)
func GenerateCommitMessage(diffs string, showTokens bool) (string, string, error) {
	ctx := context.Background()
	client, err := getClient(ctx)
	if err != nil {
		return "", "", err
	}

	generateTool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name:        "generate_commit_message",
			Description: "Generate commit message for git.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"header": {
						Type:        genai.TypeString,
						Description: "the generated commit header",
					},
					"message": {
						Type:        genai.TypeString,
						Description: "the generated commit message",
					},
				},
				Required: []string{"message", "header"},
			},
		}},
	}

	prompt := genai.Text(
		fmt.Sprintf("Use the generate_commit_message tool to generate a professional and descriptive commit header and commit message for this git commit. Should be short and descriptive of all the diffs that have happened within this commit. Here are the diffs: %s", diffs),
	)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		prompt,
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{generateTool},
		},
	)
	if err != nil {
		return "", "", err
	}

	usageData := result.UsageMetadata
	if showTokens {
		fmt.Println("Tokens used:", usageData.TotalTokenCount)
	}

	header := ""
	message := ""

	for _, part := range result.Candidates[0].Content.Parts {
		res := part.FunctionCall
		if res != nil {
			// fmt.Println(res.Args)
			header = res.Args["header"].(string)
			message = res.Args["message"].(string)
		} else {
			fmt.Println("res is nil")
		}
	}

	if strings.TrimSpace(header) == "" {
		return "", "", fmt.Errorf("There was an issue with the AI response.")
	}

	return header, message, nil
}
