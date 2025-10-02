package newspaper

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/tenjin/analysis"
)

// Newspaper manages daily newspaper generation and caching
type Newspaper struct {
	generator *Generator

	// Current newspaper cache
	currentEdition  *Edition
	editionDate     time.Time
	editionMutex    sync.RWMutex

	// Generation state
	isGenerating    bool
	generatingMutex sync.Mutex
}

// Edition represents a complete newspaper with multiple stories
type Edition struct {
	Date    time.Time
	Stories []Story
}

// NewNewspaper creates a new newspaper manager
func NewNewspaper() (*Newspaper, error) {
	generator, err := NewGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to create generator: %w", err)
	}

	return &Newspaper{
		generator: generator,
	}, nil
}

// GenerateEdition creates a new newspaper edition from current metrics
func (n *Newspaper) GenerateEdition(ctx context.Context, metrics *analysis.Metrics) error {
	// Check if already generating
	n.generatingMutex.Lock()
	if n.isGenerating {
		n.generatingMutex.Unlock()
		return fmt.Errorf("edition already being generated")
	}
	n.isGenerating = true
	n.generatingMutex.Unlock()

	defer func() {
		n.generatingMutex.Lock()
		n.isGenerating = false
		n.generatingMutex.Unlock()
	}()

	control.Log("Generating newspaper edition...")
	startTime := time.Now()

	// Collect story data from metrics
	storyData := CollectStoryData(metrics)

	// Select which stories to generate
	storyTypes := SelectStoriesToGenerate(storyData)

	// Generate each story
	stories := []Story{}
	for i, storyType := range storyTypes {
		control.Log(fmt.Sprintf("Generating story %d/%d: %s", i+1, len(storyTypes), storyType))

		var data map[string]interface{}
		switch storyType {
		case StoryTypePerformance:
			data = storyData.Performance
		case StoryTypeSentiment:
			data = storyData.Sentiment
		case StoryTypeRecord:
			if len(storyData.Records) > 0 {
				data = storyData.Records[0] // Use first record
			}
		case StoryTypeIncident:
			if len(storyData.Incidents) > 0 {
				data = storyData.Incidents[0] // Use first incident
			}
		}

		if data == nil {
			continue
		}

		story, err := n.generator.GenerateStory(ctx, storyType, data)
		if err != nil {
			control.Log(fmt.Sprintf("Failed to generate %s story: %v", storyType, err))
			continue
		}

		story.Timestamp = time.Now()
		stories = append(stories, story)
	}

	// Create new edition
	edition := &Edition{
		Date:    time.Now(),
		Stories: stories,
	}

	// Cache the edition
	n.editionMutex.Lock()
	n.currentEdition = edition
	n.editionDate = edition.Date
	n.editionMutex.Unlock()

	duration := time.Since(startTime)
	control.Log(fmt.Sprintf("Newspaper edition generated with %d stories in %.2fs", len(stories), duration.Seconds()))

	return nil
}

// GetCurrentEdition returns the current cached edition (thread-safe)
func (n *Newspaper) GetCurrentEdition() *Edition {
	n.editionMutex.RLock()
	defer n.editionMutex.RUnlock()

	return n.currentEdition
}

// NeedsNewEdition checks if a new edition should be generated
func (n *Newspaper) NeedsNewEdition() bool {
	n.editionMutex.RLock()
	defer n.editionMutex.RUnlock()

	// No edition exists
	if n.currentEdition == nil {
		return true
	}

	// Check if day has changed
	now := time.Now()
	editionDay := n.editionDate.YearDay()
	currentDay := now.YearDay()

	return editionDay != currentDay
}

// IsGenerating returns true if an edition is currently being generated
func (n *Newspaper) IsGenerating() bool {
	n.generatingMutex.Lock()
	defer n.generatingMutex.Unlock()

	return n.isGenerating
}

// GetFormattedEdition returns a human-readable string of the current edition
func (n *Newspaper) GetFormattedEdition() string {
	edition := n.GetCurrentEdition()
	if edition == nil {
		return "No newspaper edition available yet."
	}

	output := "═══════════════════════════════════════════════════\n"
	output += fmt.Sprintf("    METRO DAILY NEWS - %s\n", edition.Date.Format("January 2, 2006"))
	output += "═══════════════════════════════════════════════════\n\n"

	for i, story := range edition.Stories {
		output += fmt.Sprintf("📰 %s\n", story.Headline)
		output += fmt.Sprintf("%s\n", story.Article)

		if i < len(edition.Stories)-1 {
			output += "\n───────────────────────────────────────────────────\n\n"
		}
	}

	output += "\n═══════════════════════════════════════════════════\n"

	return output
}
