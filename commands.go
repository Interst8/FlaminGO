// Commands defines the commands available in FlaminGo

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gocolly/colly"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Defining some structures used by FlaminGo commands

// BirdSighting holds information related to a specific bird sighting in GetRecentObservations and GetRareObservations commands.
type BirdSighting struct {
	ComName string
	HowMany int
	LocName string
	ObsDt   string
}

// EmbedInfo holds information retrieved from AllAboutBirds for a bird info embed, for use in the DisplayBird() command
type EmbedInfo struct {
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

// Defining Commands

// DisplayHelp() returns a DiscordGo embed message listing FlaminGo's commands and usage
func DisplayHelp() *discordgo.MessageEmbed {
	//from: https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds
	return &discordgo.MessageEmbed{
		Color: 16711833, // Pink
		Fields: []*discordgo.MessageEmbedField{
			// !flamingo
			{
				Name:   "!flamingo",
				Value:  "Displays this list of commands",
				Inline: false,
			},
			// !get
			{
				Name:   "!get (RIT/Mendon/Braddock) {reversed}",
				Value:  "Returns a list of birds seen within 5km of the specified location in the past 2 weeks. Optionally, include 'reversed' to reverse the alphabetical order.",
				Inline: false,
			},
			// !rare
			{
				Name:   "!rare (RIT/Mendon/Braddock) {reversed}",
				Value:  "Returns a list of notable bird sightings (rare, out of season, etc.) within 15km of the specified location. Optionally, include 'reversed' to reverse the alphabetical order.",
				Inline: false,
			},
			// !bird
			{
				Name:   "!bird (Full Bird Name)",
				Value:  "Displays info for the specified bird. Uses information and names from AllAboutBirds.org.",
				Inline: false,
			},
			// !generate
			{
				Name:   "!generate {0-3}",
				Value:  "Randomly generates a bird name using a list of every bird species. Optionally, include 0-3 to specify the number of adjectives. Credit to Aidan Mahar for the lists and original idea!",
				Inline: false,
			},
		},
		Title: "FlaminGo Command Help",
	}
}

// GetRecentObservations returns a list of nearby observations in the specified radius (km) from the specified location.
func GetRecentObservations(loc Location, radius int, reverseSort bool) string {
	// Creating URL
	url := fmt.Sprintf("https://api.ebird.org/v2/data/obs/geo/recent?lat=%v&lng=%v&sort=species&dist=%d", loc.lat, loc.long, radius)
	method := "GET"

	//Creating HTTP request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	// Error handling
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	// Adding API token to header
	req.Header.Add("X-eBirdApiToken", Key)

	// Sends the request
	res, err := client.Do(req)
	//Error handling
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer res.Body.Close()

	// Body is a JSON response
	body, err := ioutil.ReadAll(res.Body)
	// Error handling
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	// Creates an array of BirdSighting and Unmarshals the contents of body into this array
	var b []BirdSighting
	err = json.Unmarshal(body, &b)
	// Error handling
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	// Sorting list of birds alphabetically, depending on whether reverseSort is true or false
	if reverseSort {
		sort.Slice(b, func(i, j int) bool {
			return b[i].ComName > b[j].ComName
		})
	} else {
		sort.Slice(b, func(i, j int) bool {
			return b[i].ComName < b[j].ComName
		})
	}

	// Formatting return string
	rString := fmt.Sprintf("**Verified eBird sightings within %d km of %v in the past 2 weeks:**\n", radius, loc.name)
	for i := 0; i < len(b); i++ {
		if b[i].HowMany > 0 {
			rString += fmt.Sprintf("%v: %d\n", b[i].ComName, b[i].HowMany)
		}
	}

	//Trimming return string to Discord max message length, minus 5 characters to show max length has been reached
	rString = truncateText(rString, 1995)

	return rString
}

// GetRareObservations returns a list of nearby notable observations in the specified radius (km) from the specified location.
// A notable observation may be a rare bird or a bird out of season.
func GetRareObservations(loc Location, radius int, reverseSort bool) string {
	// Creating URL
	url := fmt.Sprintf("https://api.ebird.org/v2/data/obs/geo/recent/notable?lat=%v&lng=%v&dist=%d&sort=species&hotspot=true", loc.lat, loc.long, radius)
	method := "GET"

	//Creating HTTP request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	// Error handling
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	// Adding API token to header
	req.Header.Add("X-eBirdApiToken", Key)

	// Sends the request
	res, err := client.Do(req)
	// Error handling
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer res.Body.Close()

	// Body is a JSON response
	body, err := ioutil.ReadAll(res.Body)
	// Error handling
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	// Creates an array of BirdSighting and Unmarshals the contents of body into this array
	var b []BirdSighting
	err = json.Unmarshal(body, &b)
	// Error handling
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	// In order to remove birds found at the same location/date, we create a map to combine these entries
	dupeMap := make(map[string]int)

	// Formatting return string
	rString := ""
	if len(b) == 0 {
		rString = "**No notable eBird sightings found.**"
	} else {
		rString = fmt.Sprintf("**Notable eBird sightings within %d km of %v in the past 2 weeks:**\n", radius, loc.name)
		for i := 0; i < len(b); i++ {
			if b[i].HowMany > 0 {
				strings := strings.Split(b[i].ObsDt, " ")
				b[i].ObsDt = strings[0] //Removing the hours/minutes from observation

				// Combining observations with same date and location
				key := (b[i].ComName + "|" + b[i].LocName + "|" + b[i].ObsDt + "|")
				if value, ok := dupeMap[key]; ok {
					dupeMap[key] = (value + b[i].HowMany)
				} else {
					dupeMap[key] = b[i].HowMany
				}
			}
		}

		// Reformatting each line after duplicate removal and adding it to an array
		var a []string
		for k, v := range dupeMap {
			t := strings.Split(k, "|")
			a = append(a, fmt.Sprintf("%v: %d [%s: %s]\n ", t[0], v, t[1], t[2]))
		}
		// Sorting list of birds alphabetically, depending on whether reverseSort is true or false
		if reverseSort {
			sort.Slice(a, func(i, j int) bool {
				return a[i] > a[j]
			})
		} else {
			sort.Slice(a, func(i, j int) bool {
				return a[i] < a[j]
			})
		}

		// Adding to return string from array
		for i := 0; i < len(a); i++ {
			rString += a[i]
		}

		//Trimming return string to Discord max message length, minus 5 characters to show max length has been reached
		rString = truncateText(rString, 1995)
	}

	return rString
}

// scrapeEmbedInfo attempts to gather information about a bird from AllAboutBirds.org.
// formattedName is a string created by messageHandler() that is given to DisplayBird() to be added to the end of the URL
func scrapeEmbedInfo(formattedName string) EmbedInfo {

	var embed EmbedInfo

	// Resolving URL
	embed.URL = fmt.Sprintf("https://www.allaboutbirds.org/guide/%s", formattedName)

	c := colly.NewCollector(
		colly.AllowedDomains("www.allaboutbirds.org"),
	)

	// Since putting in a nonexistent bird returns a search page and doesn't give an error, this acts as a form of
	// error checking by identifying an element unique to the AllAboutBirds search page.
	c.OnHTML("h1[class='page-title']", func(e *colly.HTMLElement) {
		embed.Name = "Bird not found!"
	})

	// Getting information from Species Info box
	c.OnHTML(".callout[aria-label='Species Info']", func(e *colly.HTMLElement) {
		embed.Name = e.ChildText(".species-name") // Species Name
		embed.ScientificName = e.ChildText("em")  // Scientific Name

		// Grabbing order and family information
		e.ForEach("li", func(_ int, ch *colly.HTMLElement) {
			strings := strings.Split(ch.Text, " ")
			switch strings[0] {
			case "ORDER:":
				embed.Order = strings[1]
			case "FAMILY:":
				embed.Family = strings[1]
			}
		})

		// next is a variable that is used in the following ForEach loop to determine which embed fields to fill
		next := ""
		e.ForEach("span", func(_ int, ch *colly.HTMLElement) {
			// Filling in the embed variables based on the status of 'next'
			switch next {
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

			// Switch used to set 'next'
			switch ch.Text {
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

	// Getting species description
	c.OnHTML(".speciesInfoCard", func(e *colly.HTMLElement) {
		e.ForEach("div", func(_ int, ch *colly.HTMLElement) {
			if ch.ChildText("h2") == "Basic Description" {
				embed.Description = ch.ChildText("p")
			}
		})
	})

	// Adding bird facts to slice
	c.OnHTML("li[class='is-active']", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, ch *colly.HTMLElement) {
			embed.Facts = append(embed.Facts, ch.Text)
		})
	})

	// Getting image URL
	c.OnHTML(".hero-menu", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("img", func(_ int, ch *colly.HTMLElement) bool {
			if embed.ImageURL == "" {
				// Janky, but it works until I figure out a solution for Attr("src") not working
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

	// Visits the URL, beginning the search for applicable HTML elements.
	c.Visit(embed.URL)

	return embed
}

// DisplayBird() creates and returns a Discord embed containing information about the bird.
// formattedName is a URL compatible string that is fed to scrapeEmbedInfo
func DisplayBird(formattedName string) *discordgo.MessageEmbed {
	embed := scrapeEmbedInfo(formattedName)

	// If the URL does not return a bird, the bot will return this error embed.
	if embed.Name == "Bird not found!" {
		return &discordgo.MessageEmbed{
			Color:       16711833, // Pink
			Title:       "Bird not found!",
			Description: "Make sure you spelled it right and have the name properly punctuated. Also make sure you have the full name (e.g. \"American Robin\" instead of just \"Robin\"). Birds outside of North America are unavailable.",
		}
	} else {
		// from: https://github.com/bwmarrin/discordgo/wiki/FAQ#sending-embeds
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

// GenerateBird returns a randomly generated bird name from adjectives and a noun. User can specify 0-3 adjectives, otherwise it is randomly chosen.
func GenerateBird(adjectives int) string {
	rand.Seed(time.Now().UnixNano())

	if adjectives >= 4 || adjectives < 0 {
		//Randomly generate number of adjectives
		r := rand.Intn(100)
		if r <= 9 {
			adjectives = 3
		} else if r <= 55 {
			adjectives = 2
		} else if r <= 99 {
			adjectives = 1
		} else {
			adjectives = 0
		}
	}

	//Generating noun
	bird := Nouns[rand.Intn(len(Nouns))]

	//Adding adjectives
	for i := 0; i < adjectives; i++ {
		bird = Adjectives[rand.Intn(len(Adjectives))] + " " + bird
	}

	//Capitalizing words
	bird = cases.Title(language.Und).String(bird)

	return bird
}

// truncateText trims the given string to the nearest newline character and adds an ellipse if above max length
// from: https://stackoverflow.com/a/59955447
func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	s = s[:strings.LastIndex(s[:max], "\n")]
	s += "\n..."
	return s
}
