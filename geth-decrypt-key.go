package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func checkFileExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("should not be directory")
	}
	return nil
}

func main() {
	// parse args
	keyInput := flag.String("key", "", "secret key file path")
	passwordInput := flag.String("password", "", "password file path or its content")
	flag.Parse()
	if *keyInput == "" {
		flag.PrintDefaults()
		return
	}

	// load key
	if err := checkFileExists(*keyInput); err != nil {
		fmt.Printf("Failed to load key, reason: %v", err)
		os.Exit(1)
	}
	key, err := ioutil.ReadFile(*keyInput)
	if err != nil {
		fmt.Printf("Failed to load key: %v\n", err)
		os.Exit(1)
	}

	// load password
	password := *passwordInput
	if err := checkFileExists(*passwordInput); err == nil {
		data, err := ioutil.ReadFile(*passwordInput)
		if err != nil {
			fmt.Printf("failed to load password: %v\n", err)
			os.Exit(1)
		}
		password = string(data)
	}

	// decrypt key
	result, err := keystore.DecryptKey(key, password)
	if err != nil {
		fmt.Printf("Failed to decrypt key: %v\n", err)
		os.Exit(1)
	}

	// print key
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal key: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
