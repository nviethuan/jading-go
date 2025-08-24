package utils

import (
	"testing"
	"time"

	"github.com/slack-go/slack"
)

func TestSlackClient_SendInfo_Integration(t *testing.T) {
	// Skip test if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Initialize Slack client
	slackClient := &SlackClient{}
	client := slackClient.NewSlackClient()

	// Verify client is properly initialized
	if client.client == nil {
		t.Fatal("Slack client should not be nil")
	}
	if client.infoChannel == "" {
		t.Fatal("Info channel should not be empty")
	}

	t.Run("SendInfo_WithTextBlock", func(t *testing.T) {
		// Create a simple text block message
		textBlock := slack.NewTextBlockObject("mrkdwn", "Test message from integration test", false, false)
		sectionBlock := slack.NewSectionBlock(textBlock, nil, nil)

		// Send the message
		client.SendInfo(sectionBlock)

		// Give some time for the async operation to complete
		time.Sleep(2 * time.Second)

		// Note: In a real integration test, you might want to verify the message was actually sent
		// by checking the Slack channel or using Slack's API to retrieve recent messages
		// For now, we just verify the method doesn't panic and completes execution
		// Test passes if no panic occurs
	})

	t.Run("SendInfo_WithComplexBlock", func(t *testing.T) {
		// Create a more complex message with multiple blocks
		headerText := slack.NewTextBlockObject("plain_text", "Integration Test Header", false, false)
		headerBlock := slack.NewHeaderBlock(headerText)

		bodyText := slack.NewTextBlockObject("mrkdwn", "*Bold text* and `code` with _italic_", false, false)
		bodyBlock := slack.NewSectionBlock(bodyText, nil, nil)

		dividerBlock := slack.NewDividerBlock()

		fields := []*slack.TextBlockObject{
			slack.NewTextBlockObject("mrkdwn", "*Field 1:*\nValue 1", false, false),
			slack.NewTextBlockObject("mrkdwn", "*Field 2:*\nValue 2", false, false),
		}
		fieldsBlock := slack.NewSectionBlock(nil, fields, nil)

		// Combine all blocks
		blocks := []slack.Block{headerBlock, bodyBlock, dividerBlock, fieldsBlock}

		// Send each block individually to test SendInfo
		for _, block := range blocks {
			client.SendInfo(block)
			time.Sleep(500 * time.Millisecond) // Small delay between messages
		}

		// Give time for all async operations to complete
		time.Sleep(3 * time.Second)

		// Test passes if no panic occurs
	})

	t.Run("SendInfo_WithContextBlock", func(t *testing.T) {
		// Test with context block
		contextText := slack.NewTextBlockObject("mrkdwn", "Test context message", false, false)
		contextBlock := slack.NewContextBlock("", contextText)

		client.SendInfo(contextBlock)

		time.Sleep(2 * time.Second)

		// Test passes if no panic occurs
	})
}

func TestSlackClient_NewSlackClient_Integration(t *testing.T) {
	// Skip test if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("NewSlackClient_Initialization", func(t *testing.T) {
		slackClient := &SlackClient{}
		client := slackClient.NewSlackClient()

		// Verify client initialization
		if client == nil {
			t.Fatal("Client should not be nil")
		}
		if client.client == nil {
			t.Fatal("Slack client should not be nil")
		}
		if client.infoChannel == "" {
			t.Fatal("Info channel should be set")
		}
		if client.warningChannel == "" {
			t.Fatal("Warning channel should be set")
		}
		if client.errorChannel == "" {
			t.Fatal("Error channel should be set")
		}

		// Test that calling NewSlackClient again returns the same instance
		client2 := client.NewSlackClient()
		if client != client2 {
			t.Fatal("Should return the same client instance")
		}
	})
}
