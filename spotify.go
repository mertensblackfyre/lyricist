package main

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify/v2/auth"
	"log"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"
)

func AuthSpotify(ctx context.Context) {
	
	config := &clientcredentials.Config{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	msg, page, err := client.FeaturedPlaylists(ctx)
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	fmt.Println(msg)
	for _, playlist := range page.Playlists {
		fmt.Println("  ", playlist.Name)
	}

}
