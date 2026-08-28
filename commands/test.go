// Package commands - Test Command
package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/haliphax/gobot/types"
)

var Test = types.Command{
	Meta: &discordgo.ApplicationCommand{
		Name:        "test-command",
		Description: "Test command",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("%v used test-command\n", i.Member.User)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Testing!",
			},
		})
	},
}
