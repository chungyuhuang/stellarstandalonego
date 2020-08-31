# stellarstandalonego
This simple program is to create and fund a Stellar account in a Stellar Standalone node.
## Environment variables
In order to establishing the connection to the Standalone node, user need to provide the Stellar Horzion server URL and the network phrase which should be already configured for the Standalone node.
```
export HORIZON_SERVER_URL="your horizon server url"
export NETWORK_PHRASE="your network phrase for the standalone node"
```
To fund the newly created Stellar account, user need to setup the root account address and seed.
```
export ROOT_ACCOUNT_ADDRESS="root account address"
export ROOT_ACCOUNT_SEED="root account seed"
```

Root account address can be found in the Postgress database used by the Stellar core.
Use `psql -d core -U stellar` to login to the database. Then type in `SELECT * FROM accounts;` to view the root account.

Use `stellar-core new-db` to find the seed for the root account.

Finally, set the amount you want to fund for each account.
```
export FUND_AMOUNT="amount you want to fund"
```
