package utilsCreateContractBytesCode

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nodauf/web3Toolbox/utilsWeb3"
	"log"
	"math/big"
)

func CreateContract(client *ethclient.Client, privateKeyStr string, data []byte) *types.Transaction {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatalf("Error while converting the private key: %s", err.Error())
	}

	auth, err := utilsWeb3.NewTransactor(client, privateKeyStr)
	if err != nil {
		log.Fatalf("Error while generating the transactor: %s", err.Error())
	}
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		Data: data,
	})
	if err != nil {
		log.Fatalf("Error when estimating the gas: %s", err.Error())
	}
	var txData types.LegacyTx
	txData.Data = data
	txData.Nonce = auth.Nonce.Uint64()
	txData.Value = big.NewInt(0)
	txData.Gas = gasLimit
	txData.GasPrice = auth.GasPrice
	tx := types.NewTx(&txData)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("ChainID: " + err.Error())
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Error while signing the TX: %s", err.Error())
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Error while sending the TX: %s", err.Error())
	}
	return signedTx

	//fmt.Println(hexutil.Encode(signedTx.Data()))
}
