package ethhd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/juztin/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

const DERIVATION_PATH = "m/44'/60'/0'/0/%d"

// Wallet is a wrapper around `hdwallet`
type Wallet struct {
	*hdwallet.Wallet
}

// AccountForIndex creates a public/private key pair at the given index of the derivation path, and returns the account
func (w *Wallet) AccountForIndex(index int) (accounts.Account, error) {
	p := fmt.Sprintf(DERIVATION_PATH, index)
	path := hdwallet.MustParseDerivationPath(p)
	return w.Derive(path, true)
}

// KeysForIndex creates a public/private key pair at the given index of the derivation path, and returns the public, private keys
func (w *Wallet) KeysForIndex(index int) (string, string, error) {
	var key, pkey string
	account, err := w.AccountForIndex(index)
	if err == nil {
		// Get the hex of both the public and private keys
		key = account.Address.Hex()
		pkey, err = w.PrivateKeyHex(account)
	}
	return key, pkey, err
}

// NewMnemonic generates a new HD Wallet mnemonic phrase
func NewMnemonic() (string, error) {
	var mnemonic string
	entropy, err := bip39.NewEntropy(256)
	if err == nil {
		mnemonic, err = bip39.NewMnemonic(entropy)
	}
	return mnemonic, err
}

// NewFromMnemonicAndPassword generates a new HD Wallet using the given password for BIP-39 seed, and mnemonic phrase
func NewFromMnemonicAndPassword(mnemonic, password string) (*Wallet, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return nil, err
	}
	w, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return nil, err
	}
	return &Wallet{w}, err
}

// NewFromMnemonic generates a new HD Wallet using the given mnemonic phrase
func NewFromMnemonic(mnemonic string) (*Wallet, error) {
	w, err := hdwallet.NewFromMnemonic(mnemonic)
	return &Wallet{w}, err
}

// NewFromPassword generates a new HD Wallet using the given password for BIP-39 seed, and randomly generated mnemonic phrase
func NewFromPassword(password string) (*Wallet, string, error) {
	mnemonic, err := NewMnemonic()
	if err != nil {
		return nil, mnemonic, err
	}
	w, err := NewFromMnemonicAndPassword(mnemonic, password)
	return w, mnemonic, err
}

// New generates a new HD Wallet using the given password for BIP-39 seed, and randomly generated mnemonic phrase
func New() (*Wallet, string, error) {
	mnemonic, err := NewMnemonic()
	if err != nil {
		return nil, mnemonic, err
	}
	w, err := NewFromMnemonic(mnemonic)
	return w, mnemonic, err
}
