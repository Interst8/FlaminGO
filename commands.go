//This file defines the functions used in the messageHandler in bot.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	// "time"

	"github.com/bwmarrin/discordgo"
)

//Definition of structures used in commands
type Bird struct {
	ComName string
	HowMany int
	LocName string
	ObsDt   string
}

func DisplayHelp() *discordgo.MessageEmbed { //Returns a DiscordGo embed message listing FlaminGO's commands and usage
	//from: https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds
	return &discordgo.MessageEmbed{
		Color: 16711833,
		Fields: []*discordgo.MessageEmbedField{
			//!flamingo
			{
				Name:   "!flamingo",
				Value:  "Displays this list of commands",
				Inline: false,
			},
			//!get
			{
				Name:   "!get (RIT/Mendon/Braddock)",
				Value:  "Returns a list of birds seen within 5km of the specified location in the past 2 weeks",
				Inline: false,
			},
			//!rare
			{
				Name:   "!rare",
				Value:  "Returns a list of notable bird sightings (rare, out of season, etc.) within 15km of RIT",
				Inline: false,
			},
		},
		// Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title: "FlaminGO Command Help",
	}
}

func GetRecentObs(loc Location, radius int) string { //Gets a list of nearby observations in the specified radius (km) from a location
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

	//Sorting list of birds alphabetically
	sort.Slice(b, func(i, j int) bool {
		return b[i].ComName < b[j].ComName
	})

	//Formatting return string
	rString := fmt.Sprintf("**Verified eBird sightings within %d km of %v in the past 2 weeks:**\n", radius, loc.name)
	for i := 0; i < len(b); i++ {
		if b[i].HowMany > 0 {
			rString += fmt.Sprintf("%v: %d\n", b[i].ComName, b[i].HowMany)
		}
	}

	return rString
}

func GetRareObs(loc Location, radius int) string { //Gets a list of nearby notable bird sightings
	url := fmt.Sprintf("https://api.ebird.org/v2/data/obs/geo/recent/notable?lat=%v&lng=%v&dist=%d&sort=species&hotspot=true", loc.lat, loc.long, radius)
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

	//Sorting list of birds alphabetically
	sort.Slice(b, func(i, j int) bool {
		return b[i].ComName < b[j].ComName
	})

	//Formatting return string
	rString := ""
	if len(b) == 0 {
		rString = "**No notable eBird sightings found.**"
	} else {
		rString = fmt.Sprintf("**Notable eBird sightings within %d km of %v in the past 2 weeks:**\n", radius, loc.name)
		for i := 0; i < len(b); i++ {
			if b[i].HowMany > 0 {
				strings := strings.Split(b[i].ObsDt, " ")
				b[i].ObsDt = strings[0] //Removing the hours/minutes from observation
				rString += fmt.Sprintf("%v: %d [%s: %s]\n ", b[i].ComName, b[i].HowMany, b[i].LocName, b[i].ObsDt)
			}
		}
	}
	return rString
}
