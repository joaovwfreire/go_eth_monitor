package main

import (
	"context"
	"fmt"
	"gomonitor/constants"
	"gomonitor/mariadb"
	"gomonitor/metrics"
	"gomonitor/token"

	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	alchemyMainnetKey := os.Getenv("ALCHEMY_MAINNET_KEY")

	WSUrl := "wss://eth-mainnet.g.alchemy.com/v2/" + alchemyMainnetKey

	go metrics.PrometheusStart()

	client, err := ethclient.Dial(WSUrl)

	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(constants.ERC20Address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	contractAbi, err := abi.JSON(strings.NewReader(token.TokenABI))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
			fmt.Printf("Log Index: %d\n", vLog.Index)

			switch vLog.Topics[0].Hex() {
			case logTransferSigHash.Hex():
				fmt.Println("Log Name: Transfer")

				var transferEvent mariadb.LogTransfer
				err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}

				transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

				transferEvent.Insert(transferEvent)

				/*
					fmt.Printf("From: %s\n", transferEvent.From.Hex())
					fmt.Printf("To: %s\n", transferEvent.To.Hex())
					fmt.Printf("Tokens: %s\n", transferEvent.Tokens.String())
				*/
			case logApprovalSigHash.Hex():
				fmt.Printf("Log Name: Approval\n")

				var approvalEvent mariadb.LogApproval

				err := contractAbi.UnpackIntoInterface(&approvalEvent, "Approval", vLog.Data)
				if err != nil {
					log.Fatal(err)
				}

				approvalEvent.TokenOwner = common.HexToAddress(vLog.Topics[1].Hex())
				approvalEvent.Spender = common.HexToAddress(vLog.Topics[2].Hex())

				approvalEvent.Insert(approvalEvent)

				fmt.Printf("Token Owner: %s\n", approvalEvent.TokenOwner.Hex())
				fmt.Printf("Spender: %s\n", approvalEvent.Spender.Hex())
				fmt.Printf("Tokens: %s\n", approvalEvent.Tokens.String())
			}

			fmt.Printf("\n\n")
		}
	}

}
