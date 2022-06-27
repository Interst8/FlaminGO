package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

	goBot.UpdateGameStatus(0, "!flamingo 🦩")

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

	messageTokens := strings.Split(m.Content, " ")

	//Calls DisplayHelp() to display a list of commands and their usage
	if messageTokens[0] == "!flamingo" {
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, DisplayHelp())
	}

	//Calls GetRecentObs and returns a list of birds nearby and how many were seen
	//Separate commands for locations relevant to the RIT Birding Club
	if messageTokens[0] == "!get" {
		switch messageTokens[1] {
		case "rit":
			_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObs(RIT, KM))
		case "braddock":
			_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObs(Braddock, KM))
		case "mendon":
			_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObs(Mendon, KM))
		default:
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: '%s' is not a valid option for !get", messageTokens[1]))
		}
	}

	//Calls GetRareOns and returns a list of notable bird sightings within 15km of RIT, along with their location and date of observation
	if messageTokens[0] == "!rare" {
		_, _ = s.ChannelMessageSend(m.ChannelID, GetRareObs(RIT, 15))
	}

	//Scrapes info from AllAboutBirds from the given bird name, and displays it in a Discord embed.
	if messageTokens[0] == "!bird" {
		//Constructing formatted bird name
		formatted_name := ""
		for i := 1; i < len(messageTokens); i++ {
			formatted_name += cases.Title(language.Und).String(messageTokens[i])
			if i < (len(messageTokens) - 1) {
				formatted_name += "_"
			}
		}
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, DisplayBird(formatted_name))
	}

}
