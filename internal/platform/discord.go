// Package platform - Discord bot
package platform

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/haliphax/gobot/internal/commands"
	"github.com/haliphax/gobot/internal/harness"
	"github.com/haliphax/gobot/internal/types"
)

const TypingIndicatorInterval = 10000

var (
	GuildID         = flag.String("guild", "", "Guild ID")
	BotToken        = flag.String("token", "", "Bot access token")
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

// continually update typing indicator
func keepTyping(s *discordgo.Session, channelID string, stopTyping chan bool) {
	for {
		select {
		case <-stopTyping:
			return
		default:
			if s.ChannelTyping(channelID) != nil {
				log.Printf("⚠️ WARNING: Could not send typing indicator to channel %v", channelID)
			}
			time.Sleep(TypingIndicatorInterval)
		}
	}
}

func handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := i.ApplicationCommandData().Name

	if h, ok := commandHandlers[cmd]; ok {
		var user *discordgo.User

		if i.Member != nil {
			user = i.Member.User
		} else {
			user = i.User
		}

		log.Printf("%v used %v\n", user, cmd)
		h(s, i)
	}
}

func handleChatMessage(c *harness.ModelProviderClient) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// ignore messages from the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		log.Printf("Generating response to message from %v @ %v", m.Author.Username, m.GuildID)

		// typing indicator
		stopTyping := make(chan bool, 1)
		defer func() {
			stopTyping <- true
			close(stopTyping)
		}()
		go keepTyping(s, m.ChannelID, stopTyping)

		// process message
		resp, err := (*c).ProcessUserMessage(m.Content)
		if err != nil {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("❌ ERROR: %v", err.Error()), m.Reference())
		} else {
			// send reply
			_, err = s.ChannelMessageSendReply(m.ChannelID, resp, m.Reference())
			if err != nil {
				log.Printf("❌ ERROR: %v", err.Error())
			}
		}
	}
}

// Discord bot
func Discord(c harness.ModelProviderClient, stop chan bool) {
	var (
		slashCommands = []types.Command{commands.Test}
		s             *discordgo.Session
	)

	// clean up Discord connection on return
	defer func() {
		if err := s.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	// check parameters
	s, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("❌ ERROR: Invalid bot parameters: %v", err)
	}

	// track command handlers
	for _, com := range slashCommands {
		commandHandlers[com.Meta.Name] = com.Handler
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	s.AddHandler(handleSlashCommand)
	s.AddHandler(handleChatMessage(&c))

	err = s.Open()
	if err != nil {
		log.Fatalf("❌ ERROR: Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(slashCommands))

	for i, v := range slashCommands {
		log.Printf("Adding %v", v.Meta.Name)
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v.Meta)
		if err != nil {
			log.Panicf("❌ ERROR: Cannot create '%v' command: %v", v.Meta.Name, err)
		}

		registeredCommands[i] = cmd
	}

	// wait for interrupt
	<-stop

	// notify parent on termination
	defer func() { stop <- true }()

	// cleanup
	log.Println("Removing commands...")

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
