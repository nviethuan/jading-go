package utils

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
)

type SlackClient struct {
	client         *slack.Client
	infoChannel    string
	warningChannel string
	errorChannel   string
}

func (s *SlackClient) NewSlackClient() *SlackClient {
	if s.client == nil {
		s.client = slack.New(os.Getenv("SLACK_TOKEN"))
		s.infoChannel = os.Getenv("SLACK_INFO_CHANNEL")
		s.warningChannel = os.Getenv("SLACK_WARN_CHANNEL")
		s.errorChannel = os.Getenv("SLACK_ERROR_CHANNEL")
	}
	return s
}

func (s *SlackClient) SendInfo(title string, threadID string, message ...slack.Block) chan string {
	tsChan := make(chan string)

	go func() {
		params := []slack.MsgOption{slack.MsgOptionText(title, false), slack.MsgOptionBlocks(message...)}

		if threadID != "" {
			params = append(params, slack.MsgOptionTS(threadID))
		}
		_, thread_ts, err := s.client.PostMessage(s.infoChannel, params...)
		if err != nil {
			fmt.Println("[Slack] Error: ", err)
		}

		tsChan <- thread_ts
	}()

	return tsChan
}

func (s *SlackClient) SendWarning(title string, threadID string, message ...slack.Block) chan string {
	tsChan := make(chan string)
	go func() {
		params := []slack.MsgOption{slack.MsgOptionText(title, false), slack.MsgOptionBlocks(message...)}

		if threadID != "" {
			params = append(params, slack.MsgOptionTS(threadID))
		}

		_, ts, err := s.client.PostMessage(s.warningChannel, params...)
		if err != nil {
			fmt.Println("[Slack] Error: ", err)
		}

		tsChan <- ts
	}()

	return tsChan
}

func (s *SlackClient) SendError(title string, threadID string, message ...slack.Block) chan string {
	tsChan := make(chan string)
	go func() {
		params := []slack.MsgOption{slack.MsgOptionText(title, false), slack.MsgOptionBlocks(message...)}

		if threadID != "" {
			params = append(params, slack.MsgOptionTS(threadID))
		}

		_, ts, err := s.client.PostMessage(s.errorChannel, params...)
		if err != nil {
			fmt.Println("[Slack] Error: ", err)
		}

		tsChan <- ts
	}()

	return tsChan
}
