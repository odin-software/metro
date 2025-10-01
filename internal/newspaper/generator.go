package newspaper

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

// Generator handles LLM-based story generation using Ollama
type Generator struct {
	client *api.Client
	model  string
}

// NewGenerator creates a new story generator
func NewGenerator() (*Generator, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	return &Generator{
		client: client,
		model:  "llama3.2:1b",
	}, nil
}

// GenerateStory creates a playful newspaper article from a story prompt
func (g *Generator) GenerateStory(ctx context.Context, storyType StoryType, data map[string]interface{}) (Story, error) {
	prompt := g.buildPrompt(storyType, data)

	var result strings.Builder
	req := &api.GenerateRequest{
		Model:  g.model,
		Prompt: prompt,
	}

	err := g.client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		result.WriteString(resp.Response)
		return nil
	})

	if err != nil {
		return Story{}, fmt.Errorf("failed to generate story: %w", err)
	}

	// Parse the response (headline + article)
	text := result.String()
	headline, article := g.parseResponse(text)

	return Story{
		Type:     storyType,
		Headline: headline,
		Article:  article,
	}, nil
}

// buildPrompt creates the LLM prompt based on story type and data
func (g *Generator) buildPrompt(storyType StoryType, data map[string]interface{}) string {
	switch storyType {
	case StoryTypePerformance:
		return fmt.Sprintf(
			`Write a playful newspaper article about a metro system's daily performance.

Score: %v (Grade: %s)
Passenger Satisfaction: %.1f%%
Service Efficiency: %.1f%%

Format:
HEADLINE: (catchy, one line)
ARTICLE: (2-3 sentences, playful tone, like a real newspaper but fun)`,
			data["score"], data["grade"], data["satisfaction"], data["efficiency"],
		)

	case StoryTypeRecord:
		return fmt.Sprintf(
			`Write a playful newspaper article about a record-breaking event in a metro system.

Event: %s
Value: %v

Format:
HEADLINE: (catchy, one line)
ARTICLE: (2-3 sentences, celebratory tone)`,
			data["event"], data["value"],
		)

	case StoryTypeIncident:
		return fmt.Sprintf(
			`Write a playful newspaper article about an incident in a metro system.

Incident: %s
Impact: %v passengers affected

Format:
HEADLINE: (catchy, one line)
ARTICLE: (2-3 sentences, concerned but not too serious tone)`,
			data["incident"], data["impact"],
		)

	case StoryTypeSentiment:
		return fmt.Sprintf(
			`Write a playful newspaper article about passenger sentiment in a metro system.

Average Sentiment: %.1f/100
Trend: %s
Top Station: %s (%v waiting)

Format:
HEADLINE: (catchy, one line)
ARTICLE: (2-3 sentences, empathetic tone)`,
			data["sentiment"], data["trend"], data["station"], data["waiting"],
		)

	default:
		return "Write a short playful newspaper article about a metro system."
	}
}

// parseResponse extracts headline and article from LLM output
func (g *Generator) parseResponse(text string) (headline, article string) {
	lines := strings.Split(strings.TrimSpace(text), "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Look for HEADLINE: prefix
		if strings.HasPrefix(strings.ToUpper(line), "HEADLINE:") {
			headline = strings.TrimSpace(strings.TrimPrefix(strings.ToUpper(line), "HEADLINE:"))
			headline = strings.TrimSpace(strings.TrimPrefix(line, "HEADLINE:"))
			headline = strings.TrimSpace(strings.TrimPrefix(headline, "headline:"))
			continue
		}

		// Look for ARTICLE: prefix
		if strings.HasPrefix(strings.ToUpper(line), "ARTICLE:") {
			// Collect remaining lines as article
			articleLines := []string{strings.TrimSpace(strings.TrimPrefix(line, "ARTICLE:"))}
			articleLines[0] = strings.TrimSpace(strings.TrimPrefix(articleLines[0], "article:"))

			for j := i + 1; j < len(lines); j++ {
				if strings.TrimSpace(lines[j]) != "" {
					articleLines = append(articleLines, strings.TrimSpace(lines[j]))
				}
			}
			article = strings.Join(articleLines, " ")
			break
		}

		// Fallback: first non-empty line is headline, rest is article
		if headline == "" && line != "" {
			headline = line
		} else if headline != "" && line != "" && article == "" {
			article = line
		} else if article != "" && line != "" {
			article += " " + line
		}
	}

	// Cleanup
	headline = strings.Trim(headline, `"'`)
	article = strings.TrimSpace(article)

	return headline, article
}
