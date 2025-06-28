package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/nitrix4ly/triff/core"
)

var db *core.Database

func main() {
	db = core.NewDatabase() // Assume you have this constructor

	token := os.Getenv("token")
	if token == "" {
		fmt.Println("Missing DISCORD_TOKEN environment variable")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	dg.AddMessageCreateHandler(messageHandler)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.")
	select {} // block forever
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	content := m.Content
	cmd, args, err := utils.ParseCommand(content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid command")
		return
	}

	switch cmd {
	case "GET":
		if len(args) != 1 {
			s.ChannelMessageSend(m.ChannelID, "Usage: GET <key>")
			return
		}
		val, err := db.Get(args[0])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Key not found")
		} else {
			s.ChannelMessageSend(m.ChannelID, val)
		}
	case "SET":
		if len(args) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: SET <key> <value>")
			return
		}
		db.Set(args[0], args[1])
		s.ChannelMessageSend(m.ChannelID, "OK")
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command")
	}
}
