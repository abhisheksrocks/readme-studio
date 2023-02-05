package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type ReadEnvEnvironment interface {
	Load(filenames ...string) error
	FileExist(filename string) bool
	Getenv(key string) string
}

type DefReadEnvEnvironment struct{}

func (*DefReadEnvEnvironment) Load(filenames ...string) error {
	return godotenv.Load(filenames...)
}

func (*DefReadEnvEnvironment) FileExist(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (*DefReadEnvEnvironment) Getenv(key string) string {
	return os.Getenv(key)
}

type EnvKey struct {
	Key     string
	UsedFor string
}

type envKeyValue struct {
	KeyData     EnvKey
	val         string
	environment ReadEnvEnvironment
}

type ReadEnvError string

const (
	ReadEnvErrorUndefined           ReadEnvError = "unknown error"
	ReadEnvErrorInvalidKey          ReadEnvError = "invalid key"
	ReadEnvErrorValueNotFound       ReadEnvError = "value not found"
	ReadEnvErrorFileNotFound        ReadEnvError = "file not found"
	ReadEnvErrorExampleFileNotFound ReadEnvError = "example file not found"
)

func (e envKeyValue) GetCacheValue() string {
	return e.val
}

func (e *envKeyValue) GetValue() (string, error) {
	if empty(e.KeyData.Key) {
		return "", errors.New(string(ReadEnvErrorInvalidKey))
	}
	v := e.environment.Getenv(e.KeyData.Key)
	if empty(v) {
		return "", errors.New(string(ReadEnvErrorValueNotFound))
	}
	e.val = v
	return v, nil
}

type ReadEnv struct {
	FilePath        string
	ExampleFilePath string
	KeyVal          *envKeyValue
}

func NewReadEnv(envPath string, exampleEnvPath string, key EnvKey, environment ReadEnvEnvironment) (*ReadEnv, error) {

	if notEmpty(exampleEnvPath) {
		if !environment.FileExist(exampleEnvPath) {
			return nil, errors.New(string(ReadEnvErrorExampleFileNotFound))
		}
	}

	if notEmpty(envPath) {
		err := environment.Load(envPath)
		if err != nil {
			return nil, errors.New(string(ReadEnvErrorFileNotFound))
		}
	}

	keyValueData := envKeyValue{
		KeyData:     key,
		val:         "",
		environment: environment,
	}

	_, err := keyValueData.GetValue()
	if err != nil {
		return nil, err
	}

	return &ReadEnv{
		FilePath:        envPath,
		ExampleFilePath: exampleEnvPath,
		KeyVal:          &keyValueData,
	}, err
}
