package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type envHandlerKey struct {
	keyIdentifier string
	helperText    string
}

type envHandlerKeyError struct {
	key                envHandlerKey
	envExampleFilePath string
}

type envHandlerExampleFileReadError struct {
	envExampleFilePath string
}

func (e envHandlerExampleFileReadError) Error() string {
	return `
	Couldn't read ` + e.envExampleFilePath + `
	
	TIP: It is not mandatory to have an example file. So you can skip or send ""
	But it is a good idea to always provide one for ease of use`
}

type envHandlerFileLoadError struct {
	envFilePath string
}

func (e envHandlerFileLoadError) Error() string {
	return `
	Couldn't load ` + e.envFilePath + `

	Tip: You don't necessarily have to pass this value, if the value
	is already present in system environment variables`
}

func (e envHandlerKey) Error() string {
	return `
	Couldn't read value for key "` + e.keyIdentifier + `" from environment variables
	Here's something that may explain its use:	

	` + e.helperText
}

func (e envHandlerKeyError) Error() string {
	hasExampleFile := len(e.envExampleFilePath) > 0

	toReturn := e.key.Error()

	if hasExampleFile {
		toReturn += "\n\n\tHave a look at environment file \"" + e.envExampleFilePath + "\" for setting up the variables.\n"
	}

	return toReturn
}

type envHandler struct {
	envFilePath        string
	exampleEnvFilePath string
	keyWithValue       envHandlerKeyWithValue
}

func (e envHandler) checkEnvFile() error {
	if len(e.envFilePath) != 0 {
		if err := godotenv.Load(e.envFilePath); err != nil {
			return err
		}
	}
	return nil
}

type envHandlerKeyWithValue struct {
	key           envHandlerKey
	originalValue string
}

func (key envHandlerKey) getLatestValue() (string, error) {
	value := os.Getenv(key.keyIdentifier)
	if len(value) == 0 {
		return "", key
	}
	return value, nil
}

func newEnvHandler(envFilePath string, exampleEnvFilePath string, key envHandlerKey) (envHandler, error) {

	toReturn := envHandler{}

	toReturn.exampleEnvFilePath = exampleEnvFilePath

	if len(exampleEnvFilePath) != 0 {
		if _, err := os.Stat(exampleEnvFilePath); errors.Is(err, os.ErrNotExist) {
			return toReturn, envHandlerExampleFileReadError{envExampleFilePath: exampleEnvFilePath}
		}
	}

	toReturn.envFilePath = envFilePath

	if err := toReturn.checkEnvFile(); err != nil {
		return toReturn, envHandlerFileLoadError{envFilePath: envFilePath}
	}

	value, err := key.getLatestValue()
	if err != nil {
		return toReturn, envHandlerKeyError{key: key, envExampleFilePath: exampleEnvFilePath}
	}
	toReturn.keyWithValue = envHandlerKeyWithValue{
		key:           key,
		originalValue: value,
	}

	return toReturn, nil
}
