// Package types
package types

import "github.com/bwmarrin/discordgo"

type Command struct {
	Meta    *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
