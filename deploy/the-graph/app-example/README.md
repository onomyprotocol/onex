# Connect to the indexer via ssh tunnel

```bash
export SSH_USER="your-ssh-user"
export SSH_HOST="your-ssh-host"
ssh $SSH_USER@$SSH_HOST -f -N  -L 0.0.0.0:8020:0.0.0.0:8020 -L 0.0.0.0:5001:0.0.0.0:5001 -L 0.0.0.0:8000:0.0.0.0:8000
```

To stop the tunnel later use

```bash
pkill ssh
```

# Build and deploy for onex-testnet-5

```bash
yarn install
yarn run prepare:onex-testnet-5
yarn run codegen
yarn run create-local
yarn run deploy-local
```

# Open GraphQL UI

http://0.0.0.0:8000/subgraphs/name/bank-txs/graphql

# Playground example

Example:

Request:

```graphql
{
    bankTxes(first: 10, orderBy: height, orderDirection: desc) {
        id
        height
    }
}
```

Result

````json
{
  "data": {
    "bankTxes": [
      {
        "id": "0x2a9e5509c49afbf1325e7d0b543dad7c98d25a526a523e9b09402866ef76b834",
        "height": "12833"
      },
      {
        "id": "0x90f8dafe91714efc7fa6b0b7ffa66e4a10ebdd8e6d09eb5eec89e7958f02346a",
        "height": "18047"
      }
    ]
  }
}
````