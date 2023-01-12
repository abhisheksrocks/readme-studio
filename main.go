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

const APIEndpoint = "https://api.github.com/graphql"

const (
	EnvFile        = ".env"
	ExampleEnvFile = ".env.example"
)

const (
	GithubTokenEnvKey           = "GITHUB_TOKEN"
	GithubTokenEnvKeyHelperText = "This is used as Github Authentication Token. It is required to access github's servers." +
		"\n\n\tYou can create yours at:-" +
		"\n\thttps://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token" +
		"\n\n\tWhile generating your token, no permissions or scopes are required.\n"
)

func main() {

	// Get current location path
	location, err := os.Getwd()
	if err != nil {
		log.Fatal("Cound't get the current working directory!")
		log.Fatal(err)
		return
	}

	envFileLocation := filepath.Join(location, EnvFile)
	exampleEnvFileLocation := filepath.Join(location, ExampleEnvFile)

	keyData := EnvKey{
		Key:     GithubTokenEnvKey,
		UsedFor: GithubTokenEnvKeyHelperText,
	}

	readEnv, err := NewReadEnv(envFileLocation, exampleEnvFileLocation, keyData, &DefReadEnvEnvironment{})
	if err != nil {
		switch err.Error() {
		case ReadEnvErrorExampleFileNotFound:
			log.Fatalln("\n\tCouldn't find \"" + readEnv.ExampleFilePath + "\"" +
				"\n\n\tTIP: It is not mandatory to have an example file. So you can skip this." +
				"\n\tBut it is a good idea to always provide one for ease of use\n")

		case ReadEnvErrorFileNotFound:
			log.Fatalln("\n\tCouldn't load \"" + readEnv.FilePath + "\"" +
				"\n\n\tTip: You don't necessarily have to pass this value, if the value" +
				"\n\tis already present in system environment variables\n")

		case ReadEnvErrorValueNotFound:
			defPrint := "\n\tCouldn't read value for key \"" + readEnv.KeyVal.KeyData.Key + "\" from environment variables" +
				"\n\tHere's something that may explain its use:" +
				"\n\n\t" + readEnv.KeyVal.KeyData.UsedFor

			if notEmpty(readEnv.ExampleFilePath) {
				log.Fatalln(defPrint + "\n\tUse \"" + readEnv.ExampleFilePath + "\" file for reference")
			} else {
				log.Fatalln(defPrint)
			}

		default:
			log.Fatalln(err)
		}
	}

	username := "abhisheksrocks"
	reponame := "async_button"

	query := fmt.Sprintf(`
		{
			repository(name: "%s", owner: "%s") {
				name
				isArchived
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

	request, err := http.NewRequest(http.MethodPost, APIEndpoint, bytes.NewBuffer(jsonValue))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "bearer "+readEnv.KeyVal.GetCacheValue())
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
