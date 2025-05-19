package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func tokenCacheFile() string {
	usr, err := user.Current()

	if err != nil {
		panic("unable to get current user: " + err.Error())
	}

	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700) //nolint:errcheck
	cacheFile := filepath.Join(tokenCacheDir, url.QueryEscape("plsort_secrets.json"))

	return cacheFile
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)

	return t, err
}

func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	code, err := authorize(config)

	if err != nil {
		return nil, fmt.Errorf("unable to authorize: %w", err)
	}

	fmt.Println(code)
	tok, err := config.Exchange(ctx, code)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token: %w", err)
	}

	return tok, err
}

func saveToken(file string, token *oauth2.Token) error {
	log.Printf("Saving credential file to %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token) //nolint:errcheck

	return nil
}

func getHTTPClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	cacheFile := tokenCacheFile()
	tok, err := tokenFromFile(cacheFile)

	if err != nil {
		tok, err = getTokenFromWeb(ctx, config)

		if err != nil {
			return nil, err
		}

		err = saveToken(cacheFile, tok)
	}

	return config.Client(ctx, tok), err
}

func NewYoutube(ctx context.Context, credsPath string) (*youtube.Service, error) {
	creds, err := os.ReadFile(credsPath)

	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(creds, youtube.YoutubeScope)

	if err != nil {
		return nil, err
	}

	hc, err := getHTTPClient(ctx, config)

	if err != nil {
		return nil, err
	}

	return youtube.NewService(ctx, option.WithHTTPClient(hc))
}
