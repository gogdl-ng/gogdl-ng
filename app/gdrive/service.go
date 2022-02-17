package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	credentialsFileName = "credentials.json"
	tokenFileName       = "token.json"
)

var service *drive.Service

func New() error {
	configurationDirectory, err := env.GetConfigurationFolder()

	if err != nil {
		return err
	}

	config, err := readOAuthConfigFromFile(configurationDirectory)

	if err != nil {
		return err
	}

	client, err := getAuthorizedHttpClient(configurationDirectory, config)

	if err != nil {
		return err
	}

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		return err
	}

	service = driveService

	return nil
}

func readOAuthConfigFromFile(configurationDirectory string) (*oauth2.Config, error) {
	credentialsFile := filepath.Join(configurationDirectory, credentialsFileName)
	bytes, err := ioutil.ReadFile(credentialsFile)

	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(bytes, drive.DriveReadonlyScope)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func getAuthorizedHttpClient(configurationDirectory string, config *oauth2.Config) (*http.Client, error) {
	tokenFilePath := filepath.Join(configurationDirectory, tokenFileName)

	token, err := getTokenFromFile(tokenFilePath)

	if err != nil {
		token, err = getTokenFromWeb(config)

		if err != nil {
			return nil, err
		}

		saveTokenToFile(tokenFilePath, token)
	}

	return config.Client(context.Background(), token), nil
}

func getTokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

func saveTokenToFile(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return err
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string

	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	token, err := config.Exchange(context.TODO(), authCode)

	if err != nil {
		return nil, err
	}

	return token, nil
}
