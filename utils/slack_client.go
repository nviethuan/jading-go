package utils

import (
	"os"

	"github.com/slack-go/slack"
)

type SlackClient struct {
	client *slack.Client
	infoChannel string
	warningChannel string
	errorChannel string
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

func (s *SlackClient) SendInfo(message slack.Block) {
	go s.client.PostMessage(s.infoChannel, slack.MsgOptionBlocks(message))
}

func (s *SlackClient) SendWarning(message slack.Block) {
	go s.client.PostMessage(s.warningChannel, slack.MsgOptionBlocks(message))
}

func (s *SlackClient) SendError(message slack.Block) {
	go s.client.PostMessage(s.errorChannel, slack.MsgOptionBlocks(message))
}
