package main

// A Small but Useful(tm) utility to manage different Balena profiles
// and upate a Bash prompt accordingly.

import (
	"flag"
	"fmt"
	// "os"
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

func main() {
	fmt.Println("Hello, world!")
	printPtr := flag.Bool("print", false, "print accounts")
	flag.Parse()
	if *printPtr == true {
		fmt.Println("List of accounts")
		for _, acct := range myAccounts {
			fmt.Printf("%s %s \n", acct.Name, acct.Emoji)
		}
	}
}
