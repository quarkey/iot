package webhooks

import (
	"testing"
)

func TestNewDiscord(t *testing.T) {
	wh, err := ParseConfig("../../config/discord.json")
	if err != nil {
		t.Error(err)
	}
	//TODO: add testcase
	wh.Discord.Sendf("TestNewDiscord() Message from go test")
}
