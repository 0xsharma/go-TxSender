package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var RPC_SERVER string
var SK string
var ToAddress common.Address

var Nonce uint64 = 0

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	RPC_SERVER = os.Getenv("RPC_SERVER")
	if len(RPC_SERVER) == 0 {
		fmt.Println("Invalid RPC_SERVER flag")
		return
	}
	SK = os.Getenv("SK")
	if len(SK) == 0 {
		fmt.Println("Invalid SK flag")
		return
	}
	ToAddress = common.HexToAddress(os.Getenv("TO_ADDRESS"))
	if len(SK) == 0 {
		fmt.Println("Invalid SK flag")
		return
	}

	ctx := context.Background()
	cl, err := ethclient.Dial(RPC_SERVER)
	if err != nil {
		log.Println("Error in dial connection: ", err)
	}

	chainID, err := cl.ChainID(ctx)
	if err != nil {
		log.Println("Error in fetching chainID: ", err)
	}

	sk := crypto.ToECDSAUnsafe(common.FromHex(SK))
	ksOpts, err := bind.NewKeyedTransactorWithChainID(sk, chainID)
	if err != nil {
		log.Println("Error in getting ksOpts: ", err)
	}
	add := crypto.PubkeyToAddress(sk.PublicKey)

	nonce, err := cl.PendingNonceAt(ctx, add)
	if err != nil {
		log.Fatalln("Error in getting pendingNonce: ", nonce)
	} else {
		Nonce = nonce
	}

	runTransaction(ctx, cl, ToAddress, chainID, add, ksOpts, nonce, 100)
}

func runTransaction(ctx context.Context, Clients *ethclient.Client, recipient common.Address, chainID *big.Int,
	senderAddress common.Address, opts *bind.TransactOpts, nonce uint64, value int64) {

	fmt.Println("Running transaction : ", nonce)
	var data []byte
	gasLimit := uint64(21000)

	gasPrice := big.NewInt(10000000000)

	val := big.NewInt(value)

	tx := types.NewTransaction(nonce, recipient, val, gasLimit, gasPrice, data)

	signedTx, err := opts.Signer(senderAddress, tx)
	// signedTx, err4 := types.SignTx(tx, types.NewEIP155Signer(chainID), txCfg.Acc.PrivateKey)
	if err != nil {
		log.Fatal("Error in signing tx: ", err)
	}

	log.Println(signedTx.Hash())

	err = Clients.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal("Error in sending tx: ", err)
	}
	// Nonce++
}
