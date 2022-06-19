//This file defines the functions used in the messageHandler in bot.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Bird struct {
	ComName string
	HowMany int
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

	//Formatting return string
	rString := fmt.Sprintf("Number of birds seen within %d km of %v:\n", radius, loc.name)
	for i := 0; i < len(b); i++ {
		rString += fmt.Sprintf("%v: %d\n", b[i].ComName, b[i].HowMany)
	}

	return rString
}
