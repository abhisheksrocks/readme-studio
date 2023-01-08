package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const API_ENDPOINT = "https://api.github.com/graphql"
const GITHUB_TOKEN_ENV_KEY = "GITHUB_TOKEN"

const ENV_FILE = ".env"
const EXAMPLE_ENV_FILE = ".env.example"

func main() {
	if err := godotenv.Load(ENV_FILE); err != nil {
		fmt.Printf("Coudn't find [%s] file, create this file to continue\n", ENV_FILE)
		fmt.Printf("This file holds the github access token required to make API calls\n")
		fmt.Printf("Use [%s] file for reference\n\n", EXAMPLE_ENV_FILE)
		fmt.Printf("To create Personal access token (classic), use the following URL:-\n")
		fmt.Printf("https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token\n\n")
		fmt.Printf("Not permission/scopes are required\n")
		return
	}

	githubToken := os.Getenv(GITHUB_TOKEN_ENV_KEY)

	if len(githubToken) == 0 {
		fmt.Printf("Coudn't read value for \"%s\" in [%s] file\n", GITHUB_TOKEN_ENV_KEY, ENV_FILE)
		fmt.Printf("Use [%s] file for reference\n\n", EXAMPLE_ENV_FILE)
		fmt.Printf("To create Personal access token (classic), use the following URL:-\n")
		fmt.Printf("https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token\n\n")
		fmt.Printf("Not permission/scopes are required\n")
		return
	}

	username := "abhisheksrocks"
	reponame := "ProTasks"

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
	request.Header.Add("Authorization", "bearer "+githubToken)
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
