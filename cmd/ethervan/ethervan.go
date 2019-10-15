package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
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
func outputKey(w *hdwallet.Wallet) error {
	key, pkey, err := w.KeysForIndex(0)
	if err != nil {
		return err
	}
	_, err = fmt.Println(key, pkey)
	return err
}

// writeKey creates and stores the HD addresses, writing each public key to the console
func writeKey(ksPath, ksPassword string, w *hdwallet.Wallet) error {
	ks := keystore.NewKeyStore(ksPath, keystore.StandardScryptN, keystore.StandardScryptP)
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
	a, err := ks.ImportECDSA(k, ksPassword)
	if err != nil {
		return fmt.Errorf("Failed to export account: %v", err)
	}
	_, err = fmt.Println(a.Address.Hex())
	return err
}

func main() {
	keyDir := flag.String("keystoredir", "", "Optional directory name to store the accounts")
	keyPassword := flag.String("keystorepassword", "", "Password to use for keystore encryption")
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println("Requires partial address match arg\n")
		usage()
		os.Exit(1)
	}
	if *keyDir != "" && *keyPassword == "" {
		fmt.Println("A keystorepassword is required when a keystoredir is supplied\n")
		usage()
		os.Exit(1)
	}

	var w *hdwallet.Wallet
	var mnemonic string
	var err error
	prefix := strings.ToLower("0x" + os.Args[1])
	// Create the wallet
	for i := 0; ; i++ {
		w, mnemonic, err = hdwallet.New()
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

	fmt.Println(mnemonic)

	if *keyDir != "" {
		err = writeKey(*keyDir, *keyPassword, w)
	} else {
		err = outputKey(w)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
