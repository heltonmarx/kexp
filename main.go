package main

import (
	"fmt"
	"log"
)

func main() {
	np, err := NowPlaying(kexpURL)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("#NowPlaying on KEXP :\n")
	for _, play := range np {
		fmt.Printf("%s\n", play)
	}
}
