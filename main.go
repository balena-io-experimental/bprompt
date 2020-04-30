package main

// A Small but Useful(tm) utility to manage different Balena profiles
// and upate a Bash prompt accordingly.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	// "os/user"
)

type balenaAccount struct {
	Name      string
	Url       string
	TokenName string
	Emoji     string
}

// FIXME: Hardcoding for now

var myAccounts = []balenaAccount{
	balenaAccount{
		Name:      "support",
		Url:       "balena-cloud.com",
		TokenName: "support",
		Emoji:     "🔥⚠😑",
	},
	balenaAccount{
		Name:      "personal",
		Url:       "balena-cloud.com",
		TokenName: "personal",
		Emoji:     "🐩",
	},
	balenaAccount{
		Name:      "staging",
		Url:       "balena-staging.com",
		TokenName: "staging",
		Emoji:     "🐋",
	},
	balenaAccount{
		Name:      "local",
		Url:       "my.devenv.local",
		TokenName: "local",
		Emoji:     "⛳",
	},
}

var balenaDir = "/home/hugh/.balena"
var balenaOneTrueTOken = balenaDir + "/token"
var tokenPrefix = "token."

func printAllAccounts() {
	fmt.Println("List of accounts")
	for _, acct := range myAccounts {
		fmt.Printf("%s %s \n", acct.Name, acct.Emoji)
	}
}

func switchAccount(name string) {
	var targetAcct balenaAccount
	for _, acct := range myAccounts {
		if acct.Name == name {
			targetAcct = acct
			break
		}
	}
	if targetAcct.Name == "" {
		log.Fatal("I can't find that account!")
	}
	fmt.Printf("looks like you want to switch to %s\n", targetAcct.Name)
	possibleToken := findMatchingToken(targetAcct.TokenName)
	if possibleToken == "" {
		log.Fatal("I can't find that token!")
	}
	fmt.Printf("Found that token: %s\n", possibleToken)
	updateOneTrueToken(targetAcct)
	updateBalenaRc(targetAcct)
}

func updateOneTrueToken(targetAcct balenaAccount) {
	fmt.Printf("[DEBUG] Switching to %+v\n", targetAcct)
	src := fmt.Sprintf("%s/token.%s", "/home/hugh/.balena", targetAcct.TokenName)
	fmt.Printf("[DEBUG] Source will be %s\n", src)
	target := "/home/hugh/.balena/token"
	targetStat, _ := os.Lstat(target)
	if targetStat != nil {
		if string(targetStat.Mode().String()[0]) != "L" {
			log.Fatalf("%s not a symlink, refusing to remove it! Stat: %+v", target, targetStat.Mode().String())
		}
	}
	if err := os.Remove(target); err != nil {
		log.Fatalf("Could not remove %s!")
	}
	os.Symlink(src, target)
}

func updateBalenaRc(targetAcct balenaAccount) {
	// Not great, but: we can get away for now with just a simple printf.
	urlString := []byte(fmt.Sprintf("balenaUrl: %s", targetAcct.Url))
	ioutil.WriteFile("/home/hugh/.balenarc.yml", urlString, 0755)
}

func findMatchingToken(targetToken string) string {
	allTokens := findAllTokens()
	for _, foundToken := range allTokens {
		if foundToken == fmt.Sprintf("%s%s", tokenPrefix, targetToken) {
			// Return first match
			return foundToken
		}
	}
	return ""
}

func findAllTokens() []string {
	files, err := ioutil.ReadDir(balenaDir)
	if err != nil {
		log.Fatal(err)
	}
	var tokens []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), tokenPrefix) {
			tokens = append(tokens, file.Name())
		}
	}
	return tokens
}

func findCurrentAcct() balenaAccount {
	// FIXME: Just returning a random account for now
	rand.Seed(time.Now().UnixNano())
	return myAccounts[rand.Intn(len(myAccounts))]
}

func showPromptForCurrentAccount() {
	currentAcct := findCurrentAcct()
	fmt.Printf("%s %s", currentAcct.Name, currentAcct.Emoji)
}

func main() {
	printPtr := flag.Bool("print", false, "print accounts")
	switchPtr := flag.String("switch", "", "switch accounts")
	promptPtr := flag.Bool("prompt", false, "show prompt for current account")
	flag.Parse()
	if *printPtr == true {
		printAllAccounts()
		os.Exit(0)
	}
	if *switchPtr != "" {
		switchAccount(*switchPtr)
		os.Exit(0)
	}
	if *promptPtr == true {
		showPromptForCurrentAccount()
		os.Exit(0)
	}
}
