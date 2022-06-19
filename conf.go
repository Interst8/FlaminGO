//Defines and initalizes several variables used throughout the program

package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Token string
	Key   string

	LCode = "L976278" //EBird Location code for Rochester Institute of Technology. Hardcoded for RIT Birding club, I plan to make this overridable by command
	KM    = 10        //Default size to search around location, 5 kilometers.
	Lat   = 43.08     //Latitude of RIT
	Long  = -77.67    //Longitude of RIT
	RIT   Location
)

type Location struct {
	code string
	lat  float64
	long float64
	name string
}

func init() {
	// load .env file
	err := godotenv.Load("config.env")
	if err != nil {
		fmt.Println(err.Error())
		return
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
