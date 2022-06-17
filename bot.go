package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Variables used for command line parameters
var (
	Token string
	Key   string
	BotID string
)

func init() {
	// flag.StringVar(&Token, "t", "", "Bot Token")
	// flag.Parse()

	// load .env file
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Get bot token
	Token = os.Getenv("FLAMINGO_TOK")
	//Get Ebird API Key
	Key = os.Getenv("EBIRD_KEY")

}

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
	fmt.Println("Bot is running !")
}

//Definition of messageHandler function it takes two arguments first one is discordgo.Session which is s , second one is discordgo.MessageCreate which is m.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check.
	if m.Author.ID == BotID {
		return
	}
	//If we message ping to our bot in our discord it will return us pong .
	if m.Content == "!ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}
