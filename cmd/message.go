package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Interact with Slack messages",
	Long:  `Commands for interacting with Slack messages.`,
}

var messageGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a message by timestamp and channel ID",
	Long:  `Retrieve a Slack message using its timestamp and channel ID.`,
	RunE:  runMessageGet,
}

var (
	channelID string
	timestamp string
)

func init() {
	messageCmd.AddCommand(messageGetCmd)
	messageGetCmd.Flags().StringVarP(&channelID, "channel", "c", "", "Channel ID (required)")
	messageGetCmd.Flags().StringVarP(&timestamp, "timestamp", "t", "", "Message timestamp (required)")
	messageGetCmd.MarkFlagRequired("channel")
	messageGetCmd.MarkFlagRequired("timestamp")
}

func runMessageGet(cmd *cobra.Command, args []string) error {
	// Get Slack token from environment
	token, err := getSlackToken()
	if err != nil {
		return err
	}

	// Create Slack client
	api := slack.New(token)

	// Get conversation history with the specific timestamp
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Latest:    timestamp,
		Inclusive: true,
		Limit:     1,
	}

	history, err := api.GetConversationHistory(params)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	// Check if we got any messages
	if len(history.Messages) == 0 {
		return fmt.Errorf("no message found with timestamp %s in channel %s", timestamp, channelID)
	}

	// Get the first message and verify it matches the requested timestamp
	message := history.Messages[0]
	if message.Timestamp != timestamp {
		return fmt.Errorf("no message found with exact timestamp %s in channel %s", timestamp, channelID)
	}

	// Convert message to JSON and print to stdout
	jsonData, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}
