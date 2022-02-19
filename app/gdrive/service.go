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
	"github.com/LegendaryB/gogdl-ng/app/logging"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	credentialsFileName = "credentials.json"
	tokenFileName       = "token.json"
)

var logger = logging.NewLogger()

var service *drive.Service

func New() error {
	configurationFolder := env.ConfigurationFolder

	config, err := readOAuthConfigFromFile(configurationFolder)

	if err != nil {
		logger.Errorf("failed to read oauth configuration from file. %w", err)
		return err
	}

	client, err := getAuthorizedHttpClient(configurationFolder, config)

	if err != nil {
		logger.Errorf("failed to retrieve authorized http client. %w", err)
		return err
	}

	driveService, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		logger.Errorf("failed to instantiate drive service. %w", err)
		return err
	}

	service = driveService

	return nil
}

func readOAuthConfigFromFile(configurationDirectory string) (*oauth2.Config, error) {
	credentialsFile := filepath.Join(configurationDirectory, credentialsFileName)
	bytes, err := ioutil.ReadFile(credentialsFile)

	if err != nil {
		logger.Errorf("failed to read credentials file. %w", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(bytes, drive.DriveReadonlyScope)

	if err != nil {
		logger.Errorf("failed to read configuration from credentials file. %w", err)
		return nil, err
	}

	return config, nil
}

func getAuthorizedHttpClient(configurationDirectory string, config *oauth2.Config) (*http.Client, error) {
	tokenFilePath := filepath.Join(configurationDirectory, tokenFileName)

	token, err := getTokenFromFile(tokenFilePath)

	if err != nil {
		logger.Errorf("failed to retrieve token from file. trying to retrieve it via web. %w", err)
		token, err = getTokenFromWeb(config)

		if err != nil {
			logger.Errorf("failed to retrieve token via web. %w", err)
			return nil, err
		}

		saveTokenToFile(tokenFilePath, token)
	}

	return config.Client(context.Background(), token), nil
}

func getTokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)

	if err != nil {
		logger.Errorf("failed to open token file. %w", err)
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
		logger.Errorf("failed to open token file. %w", err)
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
		logger.Errorf("failed to parse authorization code. %w", err)
		return nil, err
	}

	token, err := config.Exchange(context.TODO(), authCode)

	if err != nil {
		logger.Errorf("failed to convert authorization code to token. %w", err)
		return nil, err
	}

	return token, nil
}
