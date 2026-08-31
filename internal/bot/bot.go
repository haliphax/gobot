// Package bot
package bot

import (
	"flag"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/haliphax/gobot/internal/commands"
	"github.com/haliphax/gobot/internal/harness"
	"github.com/haliphax/gobot/internal/types"
)

var (
	GuildID  = flag.String("guild", "", "Guild ID else globally registers commands")
	BotToken = flag.String("token", "", "Bot access token")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	slashCommands = []types.Command{
		commands.Test,
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

func keepTyping(s *discordgo.Session, channelID string, stopTyping chan bool) {
	for {
		select {
		case <-stopTyping:
			return
		default:
			if s.ChannelTyping(channelID) != nil {
				log.Printf("[WARNING] Could not send typing indicator to channel %v", channelID)
			}
			time.Sleep(4000)
		}
	}
}

func init() {
	for _, c := range slashCommands {
		commandHandlers[c.Meta.Name] = c.Handler
	}

	// handle slash commands
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	})

	// handle chat messages
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// ignore messages from the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// typing indicator
		stopTyping := make(chan bool, 1)
		defer func() {
			stopTyping <- true
			close(stopTyping)
		}()
		go keepTyping(s, m.ChannelID, stopTyping)

		// process message
		resp, err := harness.HandleChat(m.Content)
		if err != nil {
			log.Fatal(err.Error())
		}

		// send reply
		_, err = s.ChannelMessageSendReply(m.ChannelID, resp, m.Reference())
		if err != nil {
			log.Fatal(err.Error())
		}
	})
}

func Main(stop chan bool) {
	defer func() {
		if err := s.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(slashCommands))

	for i, v := range slashCommands {
		log.Printf("Adding %v", v.Meta.Name)
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v.Meta)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Meta.Name, err)
		}

		registeredCommands[i] = cmd
	}

	// wait for interrupt
	<-stop

	// cleanup
	log.Println("Removing commands...")

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
