package main

import (
	"fmt"
	"github.com/chungyuhuang/stellar-go/environment"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"net/http"
	"os"
)

type StellarStandalone struct {
	client horizonclient.Client
}

// Get the current sequence number for specific Stellar account
func (op *StellarStandalone) getSeqNumber(address string) error {
	accountRequest := horizonclient.AccountRequest{AccountID: address}
	hAccount0, err := op.client.AccountDetail(accountRequest)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("seq number: %s\n", hAccount0.Sequence)
	return nil
}

// Generate a Stellar account
func generateAccount() (*keypair.Full, error) {
	pair, err := keypair.Random()
	if err != nil {
		return nil, err
	}

	fmt.Printf("address: %s\n", pair.Address())
	fmt.Printf("seed: %s\n", pair.Seed())

	return pair, nil
}


// Fund the Stellar account with 1000000 Lumen
func (op *StellarStandalone) fundAccount() error {
	pair, err := generateAccount()
	if err != nil {
		return err
	}

	err = op.getSeqNumber(pair.Address())
	if err != nil {
		return err
	}

	createAccount := txnbuild.CreateAccount{
		Destination: pair.Address(),
		Amount:      os.Getenv(environment.EnvFoundAmountKey),
	}

	accountRequest := horizonclient.AccountRequest{AccountID: os.Getenv(environment.EnvRootAccountAddress)}
	hAccount0, err := op.client.AccountDetail(accountRequest)

	// Construct the transaction that will carry the operation
	txParams := txnbuild.TransactionParams{
		SourceAccount:        &hAccount0,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&createAccount},
		Timebounds:           txnbuild.NewTimeout(300),
		BaseFee:              100,
	}

	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return err
	}

	// Sign the transaction, and base 64 encode its XDR representation
	signedTx, err := tx.SignWithKeyString(os.Getenv(environment.EnvNetworkPhrase), os.Getenv(environment.EnvRootAccountSeed))
	if err != nil {
		return err
	}

	txeBase64, _ := signedTx.Base64()
	fmt.Println("Transaction base64: ", txeBase64)

	// Submit the transaction
	resp, err := op.client.SubmitTransactionXDR(txeBase64)
	if err != nil {
		hError := err.(*horizonclient.Error)
		fmt.Println("Error submitting transaction:", hError.Problem)
		return err
	}

	fmt.Println("\nTransaction response: ", resp)
	return nil
}

func main() {
	node := StellarStandalone{
		client:horizonclient.Client{
		HorizonURL:  os.Getenv(environment.EnvHorizonServerURL),
		HTTP: http.DefaultClient,
	}}

	// need to change the base reserve to 1 in order to pass the balance check for newly created account

	err := node.fundAccount()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return
}
