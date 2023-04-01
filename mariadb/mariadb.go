package mariadb

import (
	"database/sql"
	"fmt"
	"math/big"

	_ "github.com/go-sql-driver/mysql"

	"github.com/ethereum/go-ethereum/common"
)

// LogTransfer ..
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

// LogApproval ..
type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

func (l *LogTransfer) Insert(log LogTransfer) {

	db, err := sql.Open("mysql", "a1:111@tcp(127.0.0.1:3306)/testdb")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO transfers (from_address, to_address, tokens) VALUES (?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(log.From.Hex(), log.To.Hex(), log.Tokens.String())

	if err != nil {
		fmt.Println(err.Error())
	}

}

func (l *LogApproval) Insert(log LogApproval) {
	db, err := sql.Open("mysql", "a1:111@tcp(127.0.0.1:3306)/testdb")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO approvals (token_owner, spender, tokens) VALUES (?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(log.TokenOwner.Hex(), log.Spender.Hex(), log.Tokens.String())

	if err != nil {
		fmt.Println(err.Error())
	}
}
