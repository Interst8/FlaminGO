//This file defines the functions used in the messageHandler in bot.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
)

//Definition of structures used in commands
type Bird struct {
	ComName string
	HowMany int
	LocName string
	ObsDt   string
}

type EmbedInfo struct { //Variables for a bird information embed
	Name           string
	ScientificName string
	Order          string
	Family         string
	Habitat        string
	Food           string
	Nesting        string
	Behavior       string
	Description    string
	Facts          []string
	URL            string
	ImageURL       string
}

func DisplayHelp() *discordgo.MessageEmbed { //Returns a DiscordGo embed message listing FlaminGo's commands and usage
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
			//!bird
			{
				Name:   "!bird (Full Bird Name)",
				Value:  "Displays info for the specified bird. Uses information and names from AllAboutBirds.org.",
				Inline: false,
			},
		},
		Title: "FlaminGo Command Help",
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

func scrapeEmbedInfo(formatted_name string) EmbedInfo { //Returns a message embed containing information about the bird pulled from AllAboutBirds.org

	var embed EmbedInfo

	//Resolving URL
	embed.URL = fmt.Sprintf("https://www.allaboutbirds.org/guide/%s", formatted_name)

	c := colly.NewCollector(
		colly.AllowedDomains("www.allaboutbirds.org"),
	)

	//Since putting in a nonexistent bird returns a search page and doesn't give an error, this acts as a form of
	//error checking by identifying an element unique to the search page.
	c.OnHTML("h1[class='page-title']", func(e *colly.HTMLElement) {
		embed.Name = "Bird not found!"
	})

	//Getting information from Species Info box
	c.OnHTML(".callout[aria-label='Species Info']", func(e *colly.HTMLElement) {
		embed.Name = e.ChildText(".species-name") //Species Name
		embed.ScientificName = e.ChildText("em")  //Scientific Name

		e.ForEach("li", func(_ int, ch *colly.HTMLElement) { //Grabbing order and family information
			strings := strings.Split(ch.Text, " ")
			switch strings[0] {
			case "ORDER:":
				embed.Order = strings[1]
			case "FAMILY:":
				embed.Family = strings[1]
			}
		})

		next := "" //A variable that is used in the following ForEach loop to determine which embed fields to fill
		e.ForEach("span", func(_ int, ch *colly.HTMLElement) {
			switch next { //Filling in the embed variables based on the status of 'next'
			case "habitat":
				embed.Habitat = ch.Text
				next = ""
			case "food":
				embed.Food = ch.Text
				next = ""
			case "nesting":
				embed.Nesting = ch.Text
				next = ""
			case "behavior":
				embed.Behavior = ch.Text
				next = ""
			}

			switch ch.Text { //Switch used to set 'next'
			case "Habitat":
				next = "habitat"
			case "Food":
				next = "food"
			case "Nesting":
				next = "nesting"
			case "Behavior":
				next = "behavior"
			}
		})
	})

	c.OnHTML(".speciesInfoCard", func(e *colly.HTMLElement) { //Getting species description
		e.ForEach("div", func(_ int, ch *colly.HTMLElement) {
			if ch.ChildText("h2") == "Basic Description" {
				embed.Description = ch.ChildText("p")
			}
		})
	})

	c.OnHTML("li[class='is-active']", func(e *colly.HTMLElement) { //Adding bird facts to slice
		e.ForEach("li", func(_ int, ch *colly.HTMLElement) {
			embed.Facts = append(embed.Facts, ch.Text)
		})
	})

	c.OnHTML(".hero-menu", func(e *colly.HTMLElement) { //Getting image URL
		e.ForEachWithBreak("img", func(_ int, ch *colly.HTMLElement) bool {
			if embed.ImageURL == "" {
				//Janky, but it works until I figure out a solution for Attr("src") not working
				interchange := ch.Attr("data-interchange")
				images := strings.Split(interchange, "[")
				image := strings.Split(images[3], ",")
				embed.ImageURL = image[0]
				return false
			} else {
				return false
			}
		})
	})

	c.Visit(embed.URL)

	// fmt.Println(embed.Name)
	// fmt.Println(embed.ScientificName)
	// fmt.Println(embed.Order)
	// fmt.Println(embed.Family)
	// fmt.Println(embed.Habitat)
	// fmt.Println(embed.Food)
	// fmt.Println(embed.Nesting)
	// fmt.Println(embed.Behavior)
	// fmt.Println(embed.Description)
	// fmt.Println(embed.Facts)
	// fmt.Println(embed.URL)
	// fmt.Println(embed.ImageURL)
	// fmt.Println("Done!")
	return embed
}

func DisplayBird(formatted_name string) *discordgo.MessageEmbed {
	embed := scrapeEmbedInfo(formatted_name)

	if embed.Name == "Bird not found!" {
		return &discordgo.MessageEmbed{
			Color:       16711833,
			Title:       "Bird not found!",
			Description: "Make sure you spelled it right and have the name properly punctuated. Also make sure you have the full name (e.g. \"American Robin\" instead of just \"Robin\").",
		}
	} else {
		//from: https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds
		return &discordgo.MessageEmbed{
			Color:       16711833,
			Description: embed.ScientificName,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Order",
					Value:  embed.Order,
					Inline: false,
				},
				{
					Name:   "Family",
					Value:  embed.Family,
					Inline: false,
				},
				{
					Name:   "Habitat",
					Value:  embed.Habitat,
					Inline: false,
				},
				{
					Name:   "Food",
					Value:  embed.Food,
					Inline: false,
				},
				{
					Name:   "Nesting",
					Value:  embed.Nesting,
					Inline: false,
				},
				{
					Name:   "Behavior",
					Value:  embed.Behavior,
					Inline: false,
				},
				{
					Name:   "Description",
					Value:  embed.Description,
					Inline: false,
				},
				{
					Name:   "Cool Fact",
					Value:  embed.Facts[rand.Intn(len(embed.Facts))],
					Inline: false,
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: embed.ImageURL,
			},
			URL:   embed.URL,
			Title: embed.Name,
		}
	}

}
