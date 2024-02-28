# Onex consumer chain

## Detail

The `onex-mainnet-1` chain will be launched as a consumer chain in Onomy mainnet.

- Network information: https://github.com/onomyprotocol/validator/tree/main/testnet
- Chain ID: `onex-mainnet-1`
* Spawn time: `March 4th, 2024` (Will be updated soon)
* Genesis file (without CCV): https://raw.githubusercontent.com/onomyprotocol/validator/main/testnet/onex-mainnet-1/genesis-without-ccv.json
* Genesis with CCV: Available soon
- Current version: `v1.0.3-onex`
* Binary: 
   * Version: [v1.0.3-onex](https://github.com/onomyprotocol/multiverse/releases/tag/v1.0.3-onex)
   * SHA256: `e719ac86618953dee759e111c414cded126d54b8b8ae4f9f8b21015b06c1d91c`
* Onex GitHub repository: https://github.com/onomyprotocol/multiverse
- Peers: ``
- Endpoints: 
    - RPC: ``
    - API: ``
    - gRPC: ``
- Block Explorer: ``

## IBC detail
| | onex-mainnet-1 | onomy-mainnet-1 |
|-------------|---------------------|-----------------|
|Client |`Available soon`| `Available soon`|
|Connections | `Available soon` | `Available soon` |
|Channels | `transfer`: `Available soon` <br/><br/> `consumer`: `Available soon` | `transfer`: `Available soon` <br/><br/> `consumer`: `Available soon` |

## Setup Instruction

### 1. Joining Onomy provider chain (onomy-mainnet-1) as a validator
First, validators need to run the Onomy provider chain. To set up the node and join the network, please follow the instructions in [mainnet documentation](https://github.com/onomyprotocol/validator/blob/main/mainnet/readme.md).

Here is the detail of the Onomy provider chain:
- Chain ID: `onomy-mainnet-1`
- Version: [v1.1.4](https://github.com/onomyprotocol/onomy/releases/tag/v1.1.4)
- Genesis: https://raw.githubusercontent.com/onomyprotocol/validator/main/mainnet/genesis/genesis-mainnet-1.json
- Seeds: 
```
211535f9b799bcc8d46023fa180f3359afd4c1d3@44.213.44.5:26656,00ce2f84f6b91639a7cedb2239e38ffddf9e36de@44.195.221.88:26656,cd9a47cebe8eef076a5795e1b8460a8e0b2384e5@3.210.0.126:26656,60194df601164a8b5852087d442038e392bf7470@180.131.222.74:26656,0dbe561f30862f386456734f12f431e534a3139c@34.133.228.142:26656,4737740b63d6ba9ebe93e8cc6c0e9197c426e9f4@195.189.96.106:52756
```
- Persistent Peers:
```
211535f9b799bcc8d46023fa180f3359afd4c1d3@44.213.44.5:26656,00ce2f84f6b91639a7cedb2239e38ffddf9e36de@44.195.221.88:26656,cd9a47cebe8eef076a5795e1b8460a8e0b2384e5@3.210.0.126:26656,60194df601164a8b5852087d442038e392bf7470@180.131.222.74:26656,0dbe561f30862f386456734f12f431e534a3139c@34.133.228.142:26656,4737740b63d6ba9ebe93e8cc6c0e9197c426e9f4@195.189.96.106:52756,00ce2f84f6b91639a7cedb2239e38ffddf9e36de@44.195.221.88:26656
```


### 4. Setup Onex consumer chain
The validators also need to set up the `onex-mainnet-1` consumer chain. Here are the commands to install the binary and set up the new chain.
```bash
# detail of setup will appear here
cd $HOME/go/bin
wget -O onexd https://github.com/onomyprotocol/multiverse/releases/download/v1.0.3-onex/onexd && chmod +x onexd
onexd version # v1.0.3-onex
onexd init <moniker> --chain-id onex-mainnet-1
cd $HOME/.onex/
wget -O config/genesis.json https://raw.githubusercontent.com/onomyprotocol/validator/main/testnet/onex-mainnet-1/genesis-without-ccv.json
```

The validators **MUST NOT** run the node but wait until the new genesis is published on the Onomy repository, which will be detailed in step **[5. Vote the consumer-addition proposal](#5-vote-the-consumer-addition-proposal)**.

### 5. Vote on the consumer-addition proposal
The proposal to launch `onex-mainnet-1` as a consumer chain will be submitted on the Onomy provider mainnet and the validators should participate in voting for the proposal. After the proposal is passed, the validators should wait until the `spawn_time` and replace the old genesis file with the new `genesis-with-ccv.json` file from the Onomy repository.

```bash
wget -O /$HOME/.onex/config/genesis.json https://raw.githubusercontent.com/onomyprotocol/validator/main/testnet/onex-mainnet-1/genesis.json
```

### 6. Wait for genesis and run

At the genesis time, validators can start the consumer chain by running
```bash
onexd start
```


## Launch Stages
|Step|When?|What do you need to do?|What is happening?|
|----|--------------------------------------------------|----------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
|1   |ASAP                                              |Join the Onomy mainnet `onomy-mainnet-1`  as a full node and sync to the tip of the chain.|Validator machines getting caught up on existing Composable chain's history                                                                         |
|2   | Consumer Addition proposal on provider chain | [PROVIDER] Optional: Vote for the consumer-addition proposal.  | The proposals that provide new details for the launch.                            |
|3   |The proposals passed                                 |Nothing                                                                           | The proposals passed, `spawn_time` is set. After `spawn_time` is reached, the `ccv.json` file containing `ccv` state will be provided from the provider chain.
|4   |`spawn_time` reached                                  |The `genesis-with-ccv.json` file will be provided in the testnets repo. Replace the old `genesis.json` in the `$HOME/.onex/config` directory with the new `genesis-with-ccv.json`. The new `genesis-with-ccv.json` file with ccv data will be published in [onomyprotocol/valiadtor](https://github.com/onomyprotocol/validator/tree/main/testnet/onex-mainnet-1) |
|5   |Genesis reached     | Start your node with the consumer binary | onex-mainnet-1 chain will start and become a consumer chain.                                                                                     |
|6   |3 blocks after upgrade height                     |Celebrate! :tada: 🥂                                                |<chain> blocks are now produced by the provider validator set|
