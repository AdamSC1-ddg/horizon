---
id: operations_single
title: Operation Details
category: Endpoints
---

The operation details endpoint provides information on a single [operation][resources_operation]. The operation ID provided in the `id` argument specifies which operation to load.

## Request

```
GET /operations/{id}
```

### Arguments

|  name  |  notes  | description | example |
| ------ | ------- | ----------- | ------- |
| `id` | required, number | An operation ID. | 77309415424 |

### curl Example Request

```sh
curl https://horizon-testnet.stellar.org/operations/77309415424
```

### JavaScript Example Request

```js
var StellarLib = require('js-stellar-sdk');
var server = new StellarLib.Server({hostname:'horizon-testnet.stellar.org', secure:true, port:443});

server.operations('77309415424')
  .then(function (operationsResult) {
    console.log(operationsResult)
  })
  .catch(function (err) {
    console.log(err)
  })


```

## Response

This endpoint responds with a single Operation.  See [operation resource][] for reference.

### Example Response

```json
{
  "_links": {
    "effects": {
      "href": "/operations/77309415424/effects/{?cursor,limit,order}",
      "templated": true
    },
    "precedes": {
      "href": "/operations?cursor=77309415424&order=asc"
    },
    "self": {
      "href": "/operations/77309415424"
    },
    "succeeds": {
      "href": "/operations?cursor=77309415424&order=desc"
    },
    "transactions": {
      "href": "/transactions/77309415424"
    }
  },
  "account": "GBIA4FH6TV64KSPDAJCNUQSM7PFL4ILGUVJDPCLUOPJ7ONMKBBVUQHRO",
  "funder": "GCEZWKCA5VLDNRLN3RPRJMRZOX3Z6G5CHCGSNFHEYVXM3XOJMDS674JZ",
  "id": 77309415424,
  "paging_token": "77309415424",
  "starting_balance": 1e+14,
  "type": 0,
  "type_s": "create_account"
}
```

## Problems

- The [standard problems][].
- [not_found][problems/not_found]: A `not_found` problem will be returned if there is no operation whose ID matches the `id` argument.

[operation resource]: ./resource/operation.md
[resources_operation]: ./resources/operation.md
[transaction]: ./resource/transaction.md
[problems/not_found]: ../problem/not_found.md
[standard problems]: ../guide/problems.md#Standard_Problems
