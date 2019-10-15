package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/juztin/etherhd"
)

func usage() {
	fmt.Println("Generates addresses, stopping on one matching the given phrase")
	flag.PrintDefaults()
	fmt.Println("\n  example")
	fmt.Println(`    % ethervan ffff`)
}

// outputKey writes HD public/private key-pairs to the console
func outputKey(mnemonic string, w *hdwallet.Wallet) error {
	key, pkey, err := w.KeysForIndex(0)
	if err != nil {
		return err
	}
	_, err = fmt.Printf("%s\n%s\n%s\n", mnemonic, key, pkey)
	return err
}

func writeMnemonic(mnemonic, path, addr string) error {
	p := filepath.Join(path, fmt.Sprintf("mnemonic-%s", addr))
	f, err := os.Create(p)
	if err == nil {
		_, err = io.WriteString(f, mnemonic)
	}
	return err
}

// writeKey creates and stores the HD addresses, writing each public key to the console
func writeKey(path, mnemonic, password string, w *hdwallet.Wallet) error {
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	// Create the account
	_, pkey, err := w.KeysForIndex(0)
	if err != nil {
		return fmt.Errorf("Failed to generate key: %v", err)
	}
	// Store the account with the key-store password
	k, err := crypto.HexToECDSA(pkey)
	if err != nil {
		return fmt.Errorf("Failed to generate ECDSA: %v", err)
	}
	a, err := ks.ImportECDSA(k, password)
	if err != nil {
		return fmt.Errorf("Failed to export account: %v", err)
	}
	addr := a.Address.Hex()
	err = writeMnemonic(mnemonic, path, addr)
	if err == nil {
		_, err = fmt.Println(addr)
	}
	return err
}

func main() {
	f := flag.NewFlagSet("", flag.ExitOnError)
	path := f.String("keystore", "", "Optional directory name to store the accounts")
	if len(os.Args) < 2 {
		usage()
		log.Fatalln(errors.New("Requires partial address match arg"))
	}
	f.Parse(os.Args[2:])

	var w *hdwallet.Wallet
	var mnemonic string
	prefix := strings.ToLower("0x" + os.Args[1])
	password, err := console.Stdin.PromptPassword("Password: ")
	if err != nil {
		log.Fatalln(err)
	}

	// Create the wallet
	for i := 0; ; i++ {
		if password != "" {
			w, mnemonic, err = hdwallet.NewFromPassword(password)
		} else {
			w, mnemonic, err = hdwallet.New()
		}
		if err != nil {
			log.Fatalln(err)
		}
		k, _, err := w.KeysForIndex(0)
		if err != nil {
			log.Fatalln(err)
		} else if strings.HasPrefix(strings.ToLower(k), prefix) {
			//i--
			fmt.Println()
			break
		}
		fmt.Fprintf(os.Stderr, "\033[2K\r%d generated...", i)
	}

	if *path != "" {
		err = writeKey(*path, mnemonic, password, w)
	} else {
		err = outputKey(mnemonic, w)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
