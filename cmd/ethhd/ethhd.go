package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/juztin/ethhd"
)

func usage() {
	fmt.Println("Creates and Ethereum HD wallet using a random mnemonic, or, when set, using the '--mnemonic' value")
	flag.PrintDefaults()
	fmt.Println("\n  example")
	fmt.Println(`    % ./etherhd -password "secret" -count 15`)
	fmt.Println(`    % ./etherhd -mnemonic "flower sweet join fuel gold title add language sting rocket happy open wide actual slow hover loud chuckle liquid remain wisdom foil sheriff mixed"`)
	fmt.Println(`    % ./etherhd -password "secret" -mnemonic "flower sweet join fuel gold title add language sting rocket happy open wide actual slow hover loud chuckle liquid remain wisdom foil sheriff mixed" -count 5`)
}

// outputKeys writes HD public/private key-pairs to the console
func outputKeys(mnemonic string, count, index int, w *hdwallet.Wallet) error {
	for ; index < count; index++ {
		if index == 0 {
			fmt.Println(mnemonic)
		}
		key, pkey, err := w.KeysForIndex(index)
		if err != nil {
			return err
		}
		fmt.Println(key, pkey)
	}
	return nil
}

func writeMnemonic(mnemonic, path, addr string) error {
	p := filepath.Join(path, fmt.Sprintf("mnemonic-%s", addr))
	f, err := os.Create(p)
	if err == nil {
		_, err = io.WriteString(f, mnemonic)
	}
	return err
}

// writeKeys creates and stores the HD addresses, writing each public key to the console
func writeKeys(path, mnemonic, password string, count, index int, w *hdwallet.Wallet) error {
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	for ; index < count; index++ {
		// Create the account
		_, pkey, err := w.KeysForIndex(index)
		if err != nil {
			return fmt.Errorf("Failed to generate keys at index: %d, %v", index, err)
		}
		// Store the account with the key-store password
		k, err := crypto.HexToECDSA(pkey)
		if err != nil {
			return fmt.Errorf("Failed to generate ECDSA at index: %d, %v", index, err)
		}
		a, err := ks.ImportECDSA(k, password)
		if err != nil {
			return fmt.Errorf("Failed to export account at index: %d, %v", index, err)
		}
		addr := a.Address.Hex()
		fmt.Println(addr)
		if index == 0 {
			err = writeMnemonic(mnemonic, path, addr)
			if err != nil {
				return errors.New("Failed to write mnemonic file")
			}
		}
	}
	return nil
}

func main() {
	count := flag.Int("count", 1, "The number of HD Wallet accounts to create (default 1)")
	index := flag.Int("index", 0, "The index of the derivation path (default 0)")
	mnemonic := flag.String("mnemonic", "", "A BIP-39 mnemonic to use, or a randomly generated one when not set")
	path := flag.String("keystore", "", "Optional directory name to store the accounts")
	flag.Parse()

	if *count < 1 {
		fmt.Println("\n")
		usage()
		log.Fatalln(errors.New("A count greater than 0 is required"))
	}
	if *index < 0 {
		log.Fatalln(errors.New("An index greater than, or equal to 0 is required"))
		usage()
	}
	password, err := console.Stdin.PromptPassword("Password: ")
	if err != nil {
		log.Fatalln(err)
	}

	// Create the HD Wallet
	var w *hdwallet.Wallet
	if *mnemonic == "" && password == "" {
		w, *mnemonic, err = hdwallet.New()
	} else if *mnemonic == "" {
		w, *mnemonic, err = hdwallet.NewFromPassword(password)
	} else {
		w, err = hdwallet.NewFromMnemonicAndPassword(*mnemonic, password)
	}

	if err != nil {
		log.Fatalln(err)
	}

	// Process accounts
	*count = *index + *count
	if *path != "" {
		err = writeKeys(*path, *mnemonic, password, *count, *index, w)
	} else {
		err = outputKeys(*mnemonic, *count, *index, w)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
