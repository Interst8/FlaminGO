package main

import "flag"

// "encoding/json"
// "flag"
// "fmt"
// "io/ioutil"
// "net/http"
// "os"
// "os/signal"
// "strings"
// "syscall"

// "github.com/bwmarrin/discordgo"

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}
