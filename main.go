//FlaminGo Discord bot
//Main contains all the files required for FlaminGo's functionality.
//Author: Caleb Munn

package main

//Main calls the Start() function defined in the bot file, and
func main() {
	Start()
	<-make(chan struct{})

}
