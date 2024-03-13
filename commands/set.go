package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/models"
)

var Registry map[string]interface{}

var options = map[string]interface{}{
	"PX": expire,
}

func keyInArray(key string, arr []string) (bool, int) {
	for i, item := range arr {
		if item == key || key == strings.ToUpper(item) {
			return true, i
		}
	}
	return false, -1
}

func handleOptions(parameters []string) error {
	for key, value := range options {
		found, index := keyInArray(key, parameters)
		if found {

			optionValue := parameters[index+1]
			if optionValue == "" {
				return errors.New("Invalid value for option")
			}

			err := value.(func(string, string) error)(parameters[0], optionValue)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Set(parameters []string) (string, error) {

	if len(parameters) < 2 {
		return "", models.NewNotEnoughArgsError("SET")
	}

	key := parameters[0]
	value := parameters[1]

	if Registry == nil {
		Registry = make(map[string]interface{})
	}

	err := handleOptions(parameters)

	if err != nil {
		return err.Error(), err
	}

	Registry[key] = value

	return "+OK", nil
}

func expire(key string, value string) error {

	ms, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	go (func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		fmt.Println("Deleting key: ", key)
		delete(Registry, key)
	})()

	return nil
}
