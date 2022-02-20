package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/config"
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

type DriveService struct {
	logger logging.Logger
	conf   *config.TransferConfiguration
	drive  *drive.Service
}

func NewDriveService(conf *config.TransferConfiguration, logger logging.Logger) (*DriveService, error) {
	configurationFolder := env.ConfigurationFolder

	var service = &DriveService{logger: logger}

	config, err := service.readOAuthConfigFromFile(configurationFolder)

	if err != nil {
		service.logger.Errorf("failed to read oauth configuration from file. %v", err)
		return nil, err
	}

	client, err := service.getAuthorizedHttpClient(configurationFolder, config)

	if err != nil {
		service.logger.Errorf("failed to retrieve authorized http client. %v", err)
		return nil, err
	}

	drive, err := drive.NewService(context.Background(), option.WithHTTPClient(client))

	if err != nil {
		service.logger.Errorf("failed to instantiate drive service. %v", err)
		return nil, err
	}

	service.drive = drive

	return service, nil
}

func (service *DriveService) readOAuthConfigFromFile(configurationDirectory string) (*oauth2.Config, error) {
	credentialsFile := filepath.Join(configurationDirectory, credentialsFileName)
	bytes, err := ioutil.ReadFile(credentialsFile)

	if err != nil {
		service.logger.Errorf("failed to read credentials file. %v", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(bytes, drive.DriveReadonlyScope)

	if err != nil {
		service.logger.Errorf("failed to read configuration from credentials file. %v", err)
		return nil, err
	}

	return config, nil
}

func (service *DriveService) getAuthorizedHttpClient(configurationDirectory string, config *oauth2.Config) (*http.Client, error) {
	tokenFilePath := filepath.Join(configurationDirectory, tokenFileName)

	token, err := service.getTokenFromFile(tokenFilePath)

	if err != nil {
		service.logger.Errorf("failed to retrieve token from file. trying to retrieve it via web. %v", err)
		token, err = service.getTokenFromWeb(config)

		if err != nil {
			service.logger.Errorf("failed to retrieve token via web. %v", err)
			return nil, err
		}

		service.saveTokenToFile(tokenFilePath, token)
	}

	return config.Client(context.Background(), token), nil
}

func (service *DriveService) getTokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)

	if err != nil {
		service.logger.Errorf("failed to open token file. %v", err)
		return nil, err
	}

	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

func (service *DriveService) saveTokenToFile(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		service.logger.Errorf("failed to open token file. %v", err)
		return err
	}

	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}

func (service *DriveService) getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string

	if _, err := fmt.Scan(&authCode); err != nil {
		service.logger.Errorf("failed to parse authorization code. %v", err)
		return nil, err
	}

	service.logger.Infof("auth code was: %s", authCode)

	token, err := config.Exchange(context.TODO(), authCode)

	if err != nil {
		service.logger.Errorf("failed to convert authorization code to token. %v", err)
		return nil, err
	}

	return token, nil
}
