package main

// A Small but Useful(tm) utility to manage different Balena profiles
// and upate a Bash prompt accordingly.

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type balenaAccount struct {
	Name      string
	Url       string
	TokenName string
	Emoji     string
}

// FIXME: Hardcoding for now

var balenaRcPath = "/home/hugh/.balenarc.yml"

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

type balenaRc struct {
	Url string `yaml:"balenaUrl"`
}

var balenaDir = "/home/hugh/.balena"
var balenaOneTrueToken = balenaDir + "/token"
var balenaTokenBackupDir = balenaDir + "/backup"
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
	runBalenaWhoami()
}

func runBalenaWhoami() {
	fmt.Printf("[DEBUG] balena whoami:\n")
	out, err := exec.Command("balena", "whoami").Output()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
	fmt.Printf("%s\n", out)
}

func updateOneTrueToken(targetAcct balenaAccount) {
	fmt.Printf("[DEBUG] Switching to %+v\n", targetAcct)
	src := fmt.Sprintf("%s/token.%s", balenaDir, targetAcct.TokenName)
	fmt.Printf("[DEBUG] Source will be %s\n", src)
	target := balenaOneTrueToken
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
	allTokens := findAllTokens(balenaDir)
	for _, foundToken := range allTokens {
		if foundToken == fmt.Sprintf("%s%s", tokenPrefix, targetToken) {
			// Return first match
			return foundToken
		}
	}
	return ""
}

func findAllTokens(targetDir string) []string {
	files, err := ioutil.ReadDir(targetDir)
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
	currentBalenaRc := getCurrentBalenaRc()
	currentToken := getCurrentTokenSymlinkTarget()
	for _, acct := range myAccounts {
		if acct.Url == currentBalenaRc.Url && currentToken == acct.TokenName {
			return acct
		}
	}
	return balenaAccount{}
}

func getCurrentTokenSymlinkTarget() string {
	targetStat, _ := os.Lstat(balenaOneTrueToken)
	if targetStat != nil {
		if string(targetStat.Mode().String()[0]) != "L" {
			log.Fatalf("%s not a symlink, I quit!", balenaOneTrueToken)
		}
	}
	tokenSymlink, _ := os.Readlink(balenaOneTrueToken)
	currentToken := filepath.Base(tokenSymlink)
	currentToken = strings.SplitN(currentToken, ".", 2)[1]
	return currentToken
}

func getCurrentBalenaRc() *balenaRc {
	file, err := os.Open(balenaRcPath)
	if err != nil {
		log.Fatalf("Can't open balenaRcPath: %s\n", err.Error())
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	currentBalenaRc := balenaRc{}
	err = yaml.Unmarshal(byteValue, &currentBalenaRc)
	if err != nil {
		log.Fatalf("Can't parse balenaRcPath: %s\n", err.Error())
	}
	return &currentBalenaRc
}

func showPromptForCurrentAccount() {
	currentAcct := findCurrentAcct()
	fmt.Printf("%s %s", currentAcct.Name, currentAcct.Emoji)
}

func showCurrentState() {
	currentRc := getCurrentBalenaRc()
	fmt.Printf("balenarc: %+v\n", currentRc)
	currentToken := getCurrentTokenSymlinkTarget()
	fmt.Printf("token symlink target: %+v\n", currentToken)
}

func restoreTokensFromBackup() {
	backupTokens := findAllTokens(balenaTokenBackupDir)
	for _, token := range backupTokens {

		backupTokenFullPath := balenaTokenBackupDir + "/" + token
		from, err := os.Open(backupTokenFullPath)
		if err != nil {
			fmt.Printf("[ERROR] Cannot copy %s, skipping: %s", backupTokenFullPath, err.Error())
			continue
		}
		defer from.Close()
		restoredTokenFullPath := balenaDir + "/" + token
		to, err := os.OpenFile(restoredTokenFullPath, os.O_RDWR|os.O_CREATE, 0640)
		if err != nil {
			fmt.Printf("[ERROR] Cannot open %s, skipping: %s", backupTokenFullPath, err.Error())
			continue
		}
		fmt.Printf("[DEBUG] Restoring %s to %s\n", backupTokenFullPath, restoredTokenFullPath)
		defer to.Close()
		_, err = io.Copy(to, from)
		if err != nil {
			fmt.Printf("[ERROR] Problem restoring: %s\n", err.Error())
		}
	}
}

func main() {
	printPtr := flag.Bool("print", false, "print accounts")
	switchPtr := flag.String("switch", "", "switch accounts")
	promptPtr := flag.Bool("prompt", false, "show prompt for current account")
	showPtr := flag.Bool("show", false, "show state of balenaUrl and token")
	restorePtr := flag.Bool("restore", false, "restore token files from backup (ie, from ~/.balena/backup to ~/.balena)")
	flag.Parse()
	if *printPtr == true {
		printAllAccounts()
		os.Exit(0)
	}
	if *switchPtr != "" {
		switchAccount(*switchPtr)
		os.Exit(0)
	}
	if *showPtr == true {
		showCurrentState()
		os.Exit(0)
	}
	if *promptPtr == true {
		showPromptForCurrentAccount()
		os.Exit(0)
	}
	if *restorePtr == true {
		restoreTokensFromBackup()
		os.Exit(0)
	}
}
