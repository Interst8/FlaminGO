package main

import (
	"fmt"

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
	//Calls GetRecentObs and returns a list of birds and how many were seen.
	if m.Content == "!get" {
		_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObs(RIT, KM))
	}
}
