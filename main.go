package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func commonRequestHeaders(readEnv *ReadEnv) []RequestHeader {
	return []RequestHeader{
		makeDefaultContentTypeHeader(),
		makeAuthorizationHeader(readEnv.KeyVal.GetCacheValue()),
	}
}

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

	readEnv, err := NewReadEnv(envFileLocation, exampleEnvFileLocation, keyData, new(DefReadEnvEnvironment))
	if err != nil {
		switch err.Error() {
		case string(ReadEnvErrorExampleFileNotFound):
			log.Fatalln("\n\tCouldn't find \"" + readEnv.ExampleFilePath + "\"" +
				"\n\n\tTIP: It is not mandatory to have an example file. So you can skip this." +
				"\n\tBut it is a good idea to always provide one for ease of use\n")

		case string(ReadEnvErrorFileNotFound):
			log.Fatalln("\n\tCouldn't load \"" + readEnv.FilePath + "\"" +
				"\n\n\tTip: You don't necessarily have to pass this value, if the value" +
				"\n\tis already present in system environment variables\n")

		case string(ReadEnvErrorValueNotFound):
			defPrint := "\n\tCouldn't read value for key \"" + readEnv.KeyVal.KeyData.Key + "\" from environment variables"
			if notEmpty(readEnv.KeyVal.KeyData.UsedFor) {
				defPrint += "\n\tHere's something that may explain its use:" +
					"\n\n\t" + readEnv.KeyVal.KeyData.UsedFor
			}
			if notEmpty(readEnv.ExampleFilePath) {
				defPrint += "\n\tUse \"" + readEnv.ExampleFilePath + "\" file for reference"
			}
			log.Fatalln(defPrint)

		default:
			log.Fatalln(err)
		}
	}

	username := "abhisheksrocks"
	reponame := "async_button"

	var queryResult GithubResultModel[GithubRepositoryCardModel]
	query := queryResult.Data.makeQuery(reponame, username)

	returnedError := makeRequest(APIEndpoint, query, commonRequestHeaders(readEnv),
		new(http.Client), &queryResult)

	if returnedError != nil {
		res, _ := json.MarshalIndent(returnedError, "", "    ")
		log.Println(string(res))
	} else {
		res, _ := json.MarshalIndent(queryResult, "", "    ")
		log.Println(string(res))
	}
}
