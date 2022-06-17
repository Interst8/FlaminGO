package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

var LCode = "L976278" //EBird Location code for Rochester Institute of Technology. Hardcoded for RIT Birding club, I plan to make this overridable by command
var KM = 10           //Default size to search around location, 5 kilometers.
var Lat = 43.08       //Latitude of RIT
var Long = -77.67     //Longitude of RIT
var RIT Location

type Location struct {
	code string
	lat  float64
	long float64
	name string
}

type Bird struct {
	ComName string
	HowMany int
}

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

	//Create RIT Location
	RIT.code = LCode
	RIT.lat = Lat
	RIT.long = Long
	RIT.name = "Rochester Institute of Technology"

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

//Definition of messageHandler function it takes two arguments first one is discordgo.Session which is s, second one is discordgo.MessageCreate which is m.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check.
	if m.Author.ID == BotID {
		return
	}
	//If we message ping to our bot in our discord it will return us pong .
	if m.Content == "!get" {
		_, _ = s.ChannelMessageSend(m.ChannelID, getRecentObs(RIT, KM))
	}
}

func getRecentObs(loc Location, radius int) string { //Gets a list of nearby observations in the specified radius (km) from a location
	url := fmt.Sprintf("https://api.ebird.org/v2/data/obs/geo/recent?lat=%v&lng=%v&sort=species&dist=%d", loc.lat, loc.long, radius)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	req.Header.Add("X-eBirdApiToken", Key)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body) //body is a JSON response
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	var b []Bird
	err = json.Unmarshal(body, &b)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	//Formatting return string
	rString := fmt.Sprintf("Number of birds seen within %d km of %v:\n", radius, loc.name)
	for i := 0; i < len(b); i++ {
		rString += fmt.Sprintf("%v: %d\n", b[i].ComName, b[i].HowMany)
	}

	return rString
}
