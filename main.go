package main

func main() {
	Start()
	<-make(chan struct{})
}
