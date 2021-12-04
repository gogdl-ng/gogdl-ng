package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/environment"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	credentialsFileName = "credentials.json"
	tokenFileName       = "token.json"
)

func getAuthorizedHttpClient(configDir string, config *oauth2.Config) *http.Client {
	tokenFilePath := filepath.Join(configDir, tokenFileName)

	token, err := getTokenFromFile(tokenFilePath)

	if err != nil {
		token = getTokenFromWeb(config)
		saveTokenToFile(tokenFilePath, token)
	}

	return config.Client(context.Background(), token)
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

func saveTokenToFile(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		log.Fatalf("Unable to save OAuth token: %v", err)
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string

	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)

	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return token
}

func New() *drive.Service {
	configDir, err := environment.GetConfigDir()

	if err != nil {
		log.Fatalf("Failed to get path to configuration directory: %v", err)
	}

	credentialsFile := filepath.Join(configDir, credentialsFileName)
	bytes, err := ioutil.ReadFile(credentialsFile)

	if err != nil {
		log.Fatalf("Failed to read file %v: %v", credentialsFileName, err)
	}

	config, err := google.ConfigFromJSON(bytes, drive.DriveReadonlyScope)

	if err != nil {
		log.Fatalf("Failed to parse client secret file to config: %v", err)
	}

	client := getAuthorizedHttpClient(configDir, config)

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Failed to create Google Drive service: %v", err)
	}

	return driveService
}
