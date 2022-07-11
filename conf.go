// Conf defines and initalizes several variables used throughout the program

package main

import (
	"os"

	"fmt"

	"github.com/joho/godotenv"
)

// The bot is currently being hosted via Heroku. To run it locally, change this to false.
const herokuUsed bool = true

var (
	Token string
	Key   string

	// KM is the number of kilometers around a location to search, for use in GetRecentObservations.
	KM int

	//RIT is a location representing Rochester Institute of Technology in eBird's API.
	RIT Location
	//Braddock is a location representing Braddock Bay in eBird's API.
	Braddock Location
	//Mendon is a location representing Mendon Ponds Park in eBird's API.
	Mendon Location
)

// Location holds informations about a location in eBird's API.
type Location struct {
	// Code is a string that is used by eBird's API to identify a location.
	code string
	// Lat stores the latitude value of the location, to two decimal places.
	lat float64
	// Long stores the longitude value of the location, to two decimal places.
	long float64
	// Name stores the actual name of the location, for use in printing strings.
	name string
}

func init() {

	// Since Heroku uses its own config vars, this section will cause an error if not being run locally.
	// To fix this, we put it in a conditional statement that checks to see if herokuUsed is false.
	if !herokuUsed {
		// Load .env file
		err := godotenv.Load("config.env")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	// Token stores the bot's Discord token
	Token = os.Getenv("FLAMINGO_TOK")

	// Key stores the eBird API key
	Key = os.Getenv("EBIRD_KEY")

	// Number of kilometers to search around a location
	KM = 5

	// Create RIT Location
	RIT.code = "L976278" //eBird location code
	RIT.lat = 43.08
	RIT.long = -77.67
	RIT.name = "Rochester Institute of Technology"

	// Create Braddock Bay Location
	Braddock.code = "L772198"
	Braddock.lat = 43.30
	Braddock.long = -77.71
	Braddock.name = "Braddock Bay Park"

	// Create Mendon Ponds Location
	Mendon.code = "L139800"
	Mendon.lat = 43.02
	Mendon.long = -77.57
	Mendon.name = "Mendon Ponds Park"

}
