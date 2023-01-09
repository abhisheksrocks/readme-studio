package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const API_ENDPOINT = "https://api.github.com/graphql"
const GITHUB_TOKEN_ENV_KEY = "GITHUB_TOKEN"
const GITHUB_TOKEN_ENV_KEY_HELPER_TEXT = `This is used as Github Authentication Token. It is required to access github's servers.
			
	You can create yours at:-
	https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token

	While generating your token, no permissions or scopes are required.`

const ENV_FILE = ".env"
const EXAMPLE_ENV_FILE = ".env.example"

func main() {

	// Get current location path
	location, err := os.Getwd()
	if err != nil {
		log.Fatal("Cound't get the current working directory!")
		log.Fatal(err)
		return
	}

	envFileLocation := filepath.Join(location, ENV_FILE)
	exampleEnvFileLocation := filepath.Join(location, EXAMPLE_ENV_FILE)

	tokenHandler, err := newEnvHandler(
		envFileLocation,
		exampleEnvFileLocation,
		envHandlerKey{
			GITHUB_TOKEN_ENV_KEY,
			GITHUB_TOKEN_ENV_KEY_HELPER_TEXT,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	username := "abhisheksrocks"
	reponame := "async_button"

	query := fmt.Sprintf(`
		{
			repository(name: "%s", owner: "%s") {
				name
				description
				parent {
					nameWithOwner
				}
				languages(first: 1, orderBy: {field: SIZE, direction: DESC}) {
					nodes {
						name
						color
					}
				}
				stargazerCount
				forkCount
			}
		}
	`, reponame, username)

	body := map[string]string{
		"query": query,
	}

	jsonValue, _ := json.Marshal(body)

	request, err := http.NewRequest(http.MethodPost, API_ENDPOINT, bytes.NewBuffer(jsonValue))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "bearer "+tokenHandler.keyWithValue.originalValue)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	client := &http.Client{Timeout: time.Second * 10}

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	defer response.Body.Close()
	data, _ := io.ReadAll(response.Body)
	fmt.Println(string(data))
}
