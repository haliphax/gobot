// Package commands - Test Command
package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/haliphax/gobot/internal/types"
)

var Test = types.Command{
	Meta: &discordgo.ApplicationCommand{
		Name:        "test-command",
		Description: "Test command",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Testing!",
			},
		})
	},
}
