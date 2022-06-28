//Defines and initalizes several variables used throughout the program

package main

import (
	"os"
)

var (
	Token string
	Key   string

	KM int //Radius in km around a location to search for birds

	RIT      Location
	Braddock Location
	Mendon   Location
)

type Location struct {
	code string
	lat  float64
	long float64
	name string
}

func init() {
	// load .env file

	//Turned off for Heroku deployment
	// err := godotenv.Load("config.env")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	//Get bot Discord token
	Token = os.Getenv("FLAMINGO_TOK")
	//Get Ebird API Key
	Key = os.Getenv("EBIRD_KEY")

	KM = 5

	//Create RIT Location
	RIT.code = "L976278" //eBird location code
	RIT.lat = 43.08
	RIT.long = -77.67
	RIT.name = "Rochester Institute of Technology"

	//Create Braddock Bay Location
	Braddock.code = "L772198"
	Braddock.lat = 43.30
	Braddock.long = -77.71
	Braddock.name = "Braddock Bay Park"

	//Create Mendon Ponds Location
	Mendon.code = "L139800"
	Mendon.lat = 43.02
	Mendon.long = -77.57
	Mendon.name = "Mendon Ponds Park"

}
