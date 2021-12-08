package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/urfave/cli"
)

// Config represents the config.json required to run the sample
type Config struct {
	ClientID  string   `json:"client_id"`
	Authority string   `json:"authority"`
	Scopes    []string `json:"scopes"`
}

// CreateConfig creates the Config struct from a json file.
func CreateConfig(fileName string) *Config {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func acquireTokenInteractive() (string, error) {
	config := CreateConfig("config_interactive.json")
	app, err := public.New(config.ClientID, public.WithAuthority(config.Authority), public.WithCache(&TokenCache{"cache.json"}))
	if err != nil {
		return "", err
	}
	var userAccount public.Account
	accounts := app.Accounts()
	for _, account := range accounts {
		userAccount = account
	}
	for _, scope := range config.Scopes {
		fmt.Println(scope)
	}
	result, err := app.AcquireTokenSilent(context.Background(), config.Scopes, public.WithSilentAccount(userAccount))
	if err != nil {
		result, err = app.AcquireTokenInteractive(context.Background(), config.Scopes)
		if err != nil {
			return "", err
		}
	}
	return result.AccessToken, nil
}

func callGraph(aT string) error {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	req.Header.Set("Authorization", "Bearer "+aT)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	dst := &bytes.Buffer{}
	if err := json.Indent(dst, data, "", "  "); err != nil {
		return err
	}
	fmt.Println(dst.String())
	return nil
}

func acquireTokenAndCallGraph() error {
	accessToken, err := acquireTokenInteractive()
	if err != nil {
		return err
	}
	fmt.Println("Acquired Token, calling graph")
	err = callGraph(accessToken)
	if err != nil {
		return err
	}
	return nil
}

func commands(cliApp *cli.App) {
	cliApp.Commands = []cli.Command{
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "Login using interactive auth",
			Action: func(c *cli.Context) error {
				err := acquireTokenAndCallGraph()
				if err != nil {
					return err
				}
				return nil
			},
		},
	}
}

func info(cliApp *cli.App) {
	cliApp.Name = "CLI Client"
	cliApp.Usage = "An example CLI for logging in using Interactive Authentication"
	cliApp.Author = "Abhidnya"
	cliApp.Version = "1.0.0"
}

func main() {
	cliApp := cli.NewApp()
	info(cliApp)
	commands(cliApp)
	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
