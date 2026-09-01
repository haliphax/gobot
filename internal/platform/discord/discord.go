// Package discord provides the Discord message platform
package discord

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/afero"

	"github.com/haliphax/gobot/internal/harness"
	"github.com/haliphax/gobot/internal/platform"
	"github.com/haliphax/gobot/internal/platform/discord/commands"
)

const TypingIndicatorInterval = 10000

type DiscordConfig struct {
	Token string
}

type Config struct {
	Discord DiscordConfig
}

type Discord struct {
	Session *discordgo.Session
}

var (
	GuildID         = flag.String("guild", "", "Guild ID")
	slashCommands   = []commands.Command{commands.Test}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

func persistTypingIndicator(s *discordgo.Session, channelID string, stopTyping chan bool) {
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

		log.Printf("💥 %v used %v\n", user, cmd)
		h(s, i)
	}
}

func handleChatMessage(c *harness.ModelProviderClient) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// ignore messages from the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		log.Printf("⏳ Generating response to message from %v @ %v", m.Author.Username, m.GuildID)

		// typing indicator
		stopTyping := make(chan bool, 1)
		defer func() {
			stopTyping <- true
			close(stopTyping)
		}()
		go persistTypingIndicator(s, m.ChannelID, stopTyping)

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

// Load loads the specified configuration file
func LoadConfiguration(fs afero.Fs, configFilename string) *DiscordConfig {
	file, err := fs.Open(configFilename)
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var baseConf Config
	_, err = toml.Decode(string(content), &baseConf)
	if err != nil {
		panic(err)
	}

	return &baseConf.Discord
}

// New provides an initialized Discord client instance by reference
func New(fs afero.Fs, configFilename string, c harness.ModelProviderClient) *Discord {
	conf := LoadConfiguration(fs, configFilename)

	s, err := discordgo.New("Bot " + conf.Token)
	// check parameters
	if err != nil {
		log.Fatalf("❌ ERROR: Invalid bot parameters: %v", err)
	}

	// track command handlers
	for _, com := range slashCommands {
		commandHandlers[com.Meta.Name] = com.Handler
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("🛜 Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	s.AddHandler(handleSlashCommand)
	s.AddHandler(handleChatMessage(&c))

	return &Discord{s}
}

var capabilities = platform.MessagePlatformCapabilities{
	CanEditMessage:    true,
	CanReplyToMessage: true,
	CanSendMessage:    true,
	CanThreadMessages: true,
	MaxMessageLength:  4000,
	MinMessageLength:  1,
}

func (p *Discord) Capabilities() *platform.MessagePlatformCapabilities {
	return &capabilities
}

// Start maintains a persistent Discord connection and (de)registers commands
func (p *Discord) Start(stop chan bool) {
	err := p.Session.Open()
	if err != nil {
		log.Fatalf("❌ ERROR: Cannot open the session: %v", err)
	}

	// clean up Discord connection on return
	defer func() {
		if err := p.Session.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	log.Println("⚡ Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(slashCommands))

	for i, v := range slashCommands {
		log.Printf("➕ Adding %v", v.Meta.Name)
		cmd, err := p.Session.ApplicationCommandCreate(p.Session.State.User.ID, *GuildID, v.Meta)
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
	log.Println("🗑️ Removing commands...")

	for _, v := range registeredCommands {
		err := p.Session.ApplicationCommandDelete(p.Session.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Printf("❌ ERROR: Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
