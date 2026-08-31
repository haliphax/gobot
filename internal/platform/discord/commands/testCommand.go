package commands

import "github.com/bwmarrin/discordgo"

var Test = Command{
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
