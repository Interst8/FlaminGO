//Bot contains the start function that actually runs the bot, and the messageHandler function to handle commands

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	// BotID keeps track of the bot's user ID to make sure it doesn't respond to its own messages.
	BotID string
	// adjectives and nouns are slices used when loading and using the bird generator
	Adjectives []string
	Nouns      []string
)

func Start() {
	// Creating new bot session
	goBot, err := discordgo.New("Bot " + Token)
	// Error handling
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Making our bot a user using User function.
	u, err := goBot.User("@me")
	// Error handling
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Storing our ID from u to BotID.
	BotID = u.ID

	// Adding messageHandler function to handle our messages using AddHandler from discordgo package.
	goBot.AddHandler(messageHandler)
	err = goBot.Open()
	// Error handling
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Updates FlaminGo's Discord status to display the help command, plus a cute little flamingo.
	goBot.UpdateGameStatus(0, "!flamingo 🦩")

	// Prints a string to confirm that the bot has successfully started.
	fmt.Println("Bot is running!")

	//Loading bird generator arrays
	loadGenerator("./birdgen.csv")
}

// messageHandler is called whenever a Discord message is created, and will identify if the message is a FlaminGo command.
// If the message is for FlaminGo, it will call the corresponding command function.
// s is a discordgo.Session
// m is a discordgo.MessageCreate
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Checking to see if the message author is the bot
	if m.Author.ID == BotID {
		return
	}

	m.Content = strings.ToLower(m.Content) //Lowercasing message content to standardize commands

	messageTokens := strings.Split(m.Content, " ")

	// !flamingo Calls DisplayHelp() command
	if messageTokens[0] == "!flamingo" {
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, DisplayHelp())
	}

	// !get Calls GetRecentObservations() command
	// Separate options for locations relevant to the RIT Birding Club
	if messageTokens[0] == "!get" {
		// Adding a third message token if the user did not input an optional argument, so that the bot does not crash
		if len(messageTokens) < 3 {
			messageTokens = append(messageTokens, "")
		}

		switch messageTokens[1] {
		case "rit":
			switch messageTokens[2] {
			case "reversed":
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObservations(RIT, KM, true))
			default:
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObservations(RIT, KM, false))
			}
		case "braddock":
			switch messageTokens[2] {
			case "reversed":
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObservations(Braddock, KM, true))
			default:
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObservations(Braddock, KM, false))
			}
		case "mendon":
			switch messageTokens[2] {
			case "reversed":
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObservations(Mendon, KM, true))
			default:
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRecentObservations(Mendon, KM, false))
			}
		default:
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: '%s' is not a valid option for !get", messageTokens[1]))
		}
	}

	// !rare calls GetRareObservations() command
	// Separate options for locations relevant to the RIT Birding Club
	// KM value is tripled to grant a larger search radius, due to the low amount of rare sightings.
	if messageTokens[0] == "!rare" {
		// Adding a third message token if the user did not input an optional argument, so that the bot does not crash
		if len(messageTokens) < 3 {
			messageTokens = append(messageTokens, "")
		}

		switch messageTokens[1] {
		case "rit":
			switch messageTokens[2] {
			case "reversed":
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRareObservations(RIT, KM*3, true))
			default:
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRareObservations(RIT, KM*3, false))
			}
		case "braddock":
			switch messageTokens[2] {
			case "reversed":
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRareObservations(Braddock, KM*3, true))
			default:
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRareObservations(Braddock, KM*3, false))
			}
		case "mendon":
			switch messageTokens[2] {
			case "reversed":
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRareObservations(Mendon, KM*3, true))
			default:
				_, _ = s.ChannelMessageSend(m.ChannelID, GetRareObservations(Mendon, KM*3, false))
			}
		default:
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: '%s' is not a valid option for !rare", messageTokens[1]))
		}
	}

	// !bird calls DisplayBird command
	if messageTokens[0] == "!bird" {
		// Constructing formatted bird name for use in URL
		formattedName := ""
		for i := 1; i < len(messageTokens); i++ {
			formattedName += cases.Title(language.Und).String(messageTokens[i])
			if i < (len(messageTokens) - 1) {
				formattedName += "_"
			}
		}

		// Since certain birds (e.g. Swainson's Thrush) contain an apostrophe in their name that is not used in the URL,
		// the bot will return no bird found. To avoid this, we call ReplaceAll on the URL name string to remove apostrophes.
		formattedName = strings.ReplaceAll(formattedName, "'", "")

		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, DisplayBird(formattedName))
	}

	// !generate Calls GenerateBird() command
	// User can specify
	if messageTokens[0] == "!generate" {
		// Adding a second message token if the user did not input an optional argument, so that the bot does not crash
		if len(messageTokens) < 2 {
			messageTokens = append(messageTokens, "-1")
		}
		//Converting optional argument to integer
		i, err := strconv.Atoi(messageTokens[1])
		// Error handling
		if err != nil {
			fmt.Println(err)
		}

		_, _ = s.ChannelMessageSend(m.ChannelID, GenerateBird(i))

	}

}

//loadGenerator loads the given .csv file path into the bird generator
func loadGenerator(file string) {
	//Opening reader with .csv file
	r, err := os.Open(file)
	//Error checking
	if err != nil {
		fmt.Println(err.Error())
		return

	}
	defer r.Close()

	//Creating csv reader
	reader := csv.NewReader(r)

	records, err := reader.ReadAll()
	//Error checking
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Adding fields to respective arrays
	for i := 0; i < len(records); i++ {
		Adjectives = append(Adjectives, records[i][0])
		if len(records[i][1]) > 0 {
			Nouns = append(Nouns, records[i][1])
		}
	}
}
