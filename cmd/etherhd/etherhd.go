package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/juztin/etherhd"
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
func outputKeys(count, index int, w *hdwallet.Wallet) error {
	for ; index < count; index++ {
		key, pkey, err := w.KeysForIndex(index)
		if err != nil {
			return err
		}
		fmt.Println(key, pkey)
	}
	return nil
}

// writeKeys creates and stores the HD addresses, writing each public key to the console
func writeKeys(ksPath, ksPassword string, count, index int, w *hdwallet.Wallet) error {
	ks := keystore.NewKeyStore(ksPath, keystore.StandardScryptN, keystore.StandardScryptP)
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
		a, err := ks.ImportECDSA(k, ksPassword)
		if err != nil {
			return fmt.Errorf("Failed to export account at index: %d, %v", index, err)
		}
		fmt.Println(a.Address.Hex())
	}
	return nil
}

func main() {
	count := flag.Int("count", 1, "The number of HD Wallet accounts to create (default 1)")
	index := flag.Int("index", 0, "The index of the derivation path (default 0)")
	password := flag.String("password", "", "The password used for the BIP-39 seed")
	mnemonic := flag.String("mnemonic", "", "A BIP-39 mnemonic to use, or a randomly generated one when not set")
	keyDir := flag.String("keystoredir", "", "Optional directory name to store the accounts")
	keyPassword := flag.String("keystorepassword", "", "Password to use for keystore encryption")
	flag.Parse()
	if *count < 1 {
		fmt.Println("A count greater than 0 is required\n")
		usage()
		os.Exit(1)
	}
	if *index < 0 {
		fmt.Println(*index)
		fmt.Println("An index greater than, or equal to 0 is required\n")
		usage()
		os.Exit(1)
	}
	if *keyDir != "" && *keyPassword == "" {
		fmt.Println("A keystorepassword is required when a keystoredir is supplied\n")
		usage()
		os.Exit(1)
	}

	// Create the HD Wallet
	var w *hdwallet.Wallet
	var err error
	if *mnemonic == "" && *password == "" {
		w, *mnemonic, err = hdwallet.New()
		fmt.Println(*mnemonic)
	} else if *mnemonic == "" {
		w, *mnemonic, err = hdwallet.NewFromPassword(*password)
		fmt.Println(*mnemonic)
	} else {
		w, err = hdwallet.NewFromMnemonicAndPassword(*mnemonic, *password)
	}

	if err != nil {
		log.Fatalln(err)
	}

	// Process accounts
	*count = *index + *count
	if *keyDir != "" {
		err = writeKeys(*keyDir, *keyPassword, *count, *index, w)
	} else {
		err = outputKeys(*count, *index, w)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
