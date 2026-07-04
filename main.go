package main

import (
	"flag"
	"fmt"

	"context"
)

func main() {

	ctx := context.Background()
	url := flag.String("url", "none", "a string")
	typer := flag.String("type", "none", "a string")

	flag.Parse()

	if *url == "none" || *typer == "none" {
		fmt.Println("Please provide url")
		//return
	}


	AuthSpotify(ctx)
	fmt.Println(*url, *typer)
}
