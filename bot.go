package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	BotID string
)

func Start() {
	//creating new bot session
	goBot, err := discordgo.New("Bot " + Token)

	//Handling error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Making our bot a user using User function .
	u, err := goBot.User("@me")
	//Handlinf error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Storing our id from u to BotId .
	BotID = u.ID

	// Adding handler function to handle our messages using AddHandler from discordgo package. We will declare messageHandler function later.
	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	//Error handling
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//If every thing works fine we will be printing this.
	fmt.Println("Bot is running!")
}

//Definition of messageHandler function. It takes two arguments first one is discordgo.Session which is s, second one is discordgo.MessageCreate which is m.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check.
	if m.Author.ID == BotID {
		return
	}

	m.Content = strings.ToLower(m.Content) //Lowercasing message content to standardize commands

	//Calls DisplayHelp() to display a list of commands and their usage
	if m.Content == "!flamingo" {
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, DisplayHelp())
	}

	//Calls GetRecentObs and returns a list of birds nearby and how many were seen.
	//Separate commands for locations relevant to the RIT Birding Club
	if m.Content == "!get rit" {
		_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObs(RIT, KM))
	}

	if m.Content == "!get braddock" || m.Content == "!get braddock bay" || m.Content == "!get braddock bay park" {
		_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObs(Braddock, KM))
	}

	if m.Content == "!get mendon" || m.Content == "!get mendon ponds" || m.Content == "!get mendon ponds park" {
		_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObs(Mendon, KM))
	}

}
