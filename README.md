# Go Monitor 
A simple Go program to monitor a smart contract events on the Ethereum blockchain.
It uses the [go-ethereum]() library to connect to the blockchain and [go-ethereum/ethclient]() to interact with the smart contract.

It updates a local mariadb database with the events data with concurrent go routines. 
It also concurrently exposes a Prometheus endpoint to monitor the application at "localhost:2112/metrics".

## Database setup
The database is a MariaDB database.
The database is created with the following commands:
```sql
CREATE TABLE transfers (
  from_address varchar(42) NOT NULL,
  to_address varchar(42) NOT NULL,
  tokens decimal(65,0) NOT NULL
);

CREATE TABLE approvals(
  token_owner varchar(42) NOT NULL,
  spender varchar(42) NOT NULL,
  tokens decimal(65,0) NOT NULL
);
```
and is hosted on a local mariadb server at the testdb database.

## Configuration

