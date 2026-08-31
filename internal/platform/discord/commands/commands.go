// Package commands provides slash commands for the Discord interface
package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Meta    *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
