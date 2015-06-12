---
id: readme
title: Introduction
category: Getting Started
---

Horizon is the client facing API server for the Stellar ecosystem.  See [an overview of the Stellar ecosystem](#) for more details.

Horizon provides three significant portions of functionality:  The Transactions API, the History API, and the Trading API.

## Transactions API

The Transactions API exists to help you make transactions against the Stellar network.  It provides ways to help you create valid transactions, such as providing an account's sequence number or latest known balances.

In addition to the read endpoints, the Transactions API also provides the endpoint to submit transactions.

### Important Endpoints

- [Post transaction]({{< relref "endpoint/transactions_create.md" >}})
- [Account details]({{< relref "endpoint/accounts_one.md" >}})
- [Calculate payment path](#)

## History API

The History API provides endpoints for retrieving data about what has happened in the past on the Stellar network.  It provides (or will provide) endpoints that let you:

- Retrieve transaction details
- Load transactions that effect a given account
- Load payment history for an account
- Load trade history for a given order book

### Important Endpoints

- [Transactions for account]({{< relref "endpoint/transactions_for_account.md" >}})
- [Transaction fetails]({{< relref "endpoint/transactions_one.md" >}})
- [All ledgers]({{< relref "endpoint/ledgers_all.md" >}})

## Trading API

The Trading API provides endpoints for retrieving data about the distributed
currency exchange within stellar.  It provides data regarding open offers to
exchange currency (often called an order book) and also provides data about
trades that were executed within the exchange.

### Important Endpoints

- [Orderbook details](#)
- [Trades for orderbook](#)

