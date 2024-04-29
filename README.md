# ONEX Consumer Chain

## Hybrid AMM + On-chain Orderbook DEX

![GH2AFkdW4AAcOe4](https://github.com/onomyprotocol/onex/assets/76499838/83d43179-3781-4e7a-95af-f3b054bc8e74)

## Detail

The `onex-mainnet-1` chain will be launched as a consumer chain in Onomy mainnet.

- Network information: https://github.com/onomyprotocol/onex/tree/main/chain/onex-mainnet-1
- Chain ID: `onex-mainnet-1`
* Spawn time: `April 28th, 2024 17:00 UTC`
* Spawn time: `April 29th, 2024 16:00 UTC`
* Genesis file (without CCV): https://raw.githubusercontent.com/onomyprotocol/onex/main/chain/onex-mainnet-1/genesis-without-ccv.json
* Genesis with CCV: 
   * URL: https://github.com/onomyprotocol/onex/blob/main/chain/onex-mainnet-1/genesis-with-ccv.json
   * SHA256: `cec18fcc1e984be3e972e7b9a053c8017e3f172a767613966e7de69e203a7406`
* Binary: 
   * Version: [v1.1.0](https://github.com/onomyprotocol/onex/releases/tag/v1.1.0)
   * SHA256: `5961859fb363b31beaf7da751a88ff3517d34d28807b568fde4b55c51d76ab77`
* Onex GitHub repository: https://github.com/onomyprotocol/onex
- Peers: `2f96d16645fd52dba217fb477a66c7b637fbb3c7@64.71.153.56:26756,e6e0a2fef354c509f31d573305626cc2a5cc9982@64.71.153.53:26756,df54649025577784c68dfbcd01ef9c59fb62e401@180.131.222.71:26756,48fbb5128a9b044d90cfa28de5078fa7e40f8b5b@180.131.222.72:26756,f80867e8181a07b26a17e4f597b0cfb7408b1b2a@180.131.222.73:26756,eb823e14ff73127ccce3e17bd674046b290416f1@51.250.106.107:36656`
- Endpoints: 
    - RPC: `https://rpc-onex.decentrio.ventures`
    - API: `https://api-onex.decentrio.ventures`
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

### 2. Setup Onex consumer chain
The validators also need to set up the `onex-mainnet-1` consumer chain. Here are the commands to install the binary and set up the new chain.
```bash
# detail of setup will appear here
cd $HOME/go/bin
wget -O onexd https://github.com/onomyprotocol/onex/releases/download/v1.1.0/onexd && chmod +x onexd
onexd version # v1.1.0
onexd init <moniker> --chain-id onex-mainnet-1
cd $HOME/.onomy_onex
wget -O config/genesis.json https://raw.githubusercontent.com/onomyprotocol/onex/main/chain/onex-mainnet-1/genesis-with-ccv.json
```

### 3. Key assignment

One important thing the validators are encouraged, but not required to do when a consumer chain is being onboarded, is to assign a key different from the key they use on the provider chain. However, if the validators choose to have different validator keys on consumer chains, they ***MUST*** do key-assignment before the `spawn_time` and only handle on provider chain. Failing to assign key during these times can cause liveliness problems with the consumer chains. If validators have not used key assignment for the consumer node yet, please wait until Onex is live and receiving valset updates from the provider before doing so.

To get current validator key in consumer chains, the validators should run this command:
```bash
# Run this command to get the validator key, which is encrypted 
# in `ed25519` format, not `secp256k1`
onexd tendermint show-validator
> {"@type":"/cosmos.crypto.ed25519.PubKey","key":"abcdDefGh123456789="}
```

Assigning different key to the consumer: 
```bash
onomyd tx provider assign-consensus-key onex-mainnet-1 <consensus-key-above> --from <signer> 
```

To recheck the keys:
```bash
onomyd query provider validator-consumer-key onex-mainnet-1 $(onomyd tendermint show-address)
```

**Note**: If the validators choose not to use different validator keys, the validator keys used in Onomy can be simply copied to the Onex chain config.

```bash
cp .onomy/config/priv_validator_key.json .onomy_onex/config/priv_validator_key.json
```

When the chain is started, valdiators can check if it is ready to sign blocks with:

```bash
onexd query tendermint-validator-set | grep "$(onexd tendermint show-address)"
```
### 4. Wait for genesis and run

At the genesis time, validators can start the consumer chain by running
```bash
onexd start
```

> Note: if validators choose to run onex and onomy in the same machine, it is highly recommended to setup separate ports to prevent clashing. These ports are: P2P, RPC, REST, gRPC, gRPC-web

The validators can also use service to run and monitor the node. Here is the example of `/etc/systemd/system/onex.service`:
```
[Unit]
Description=Onex node
After=network.target

[Service]
ExecStart=/$HOME/go/bin/onexd start --p2p.persistent_peers="2f96d16645fd52dba217fb477a66c7b637fbb3c7@64.71.153.56:26756,e6e0a2fef354c509f31d573305626cc2a5cc9982@64.71.153.53:26756,df54649025577784c68dfbcd01ef9c59fb62e401@180.131.222.71:26756,48fbb5128a9b044d90cfa28de5078fa7e40f8b5b@180.131.222.72:26756,f80867e8181a07b26a17e4f597b0cfb7408b1b2a@180.131.222.73:26756,eb823e14ff73127ccce3e17bd674046b290416f1@51.250.106.107:36656"
Restart=always
RestartSec=3
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

After that, run these commands to enable and start the chain:
```bash
systemctl daemon-reload
systemctl enable onex.service
systemctl restart onex.service
```
and run `journalctl -fu onex -n150` to check the log. 


## Launch Stages
|Step|When?|What do you need to do?|What is happening?|
|----|--------------------------------------------------|----------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
|1   |ASAP                                              |Join the Onomy mainnet `onomy-mainnet-1`  as a full node and sync to the tip of the chain.|Validator machines getting caught up on existing Composable chain's history                                                                         |
|2   | Consumer Addition proposal on provider chain | [PROVIDER] Optional: Vote for the consumer-addition proposal.  | The proposals that provide new details for the launch.                            |
|3   |The proposals passed                                 |Nothing                                                                           | The proposals passed, `spawn_time` is set. After `spawn_time` is reached, the `ccv.json` file containing `ccv` state will be provided from the provider chain.
|4   |`spawn_time` reached                                  |The `genesis-with-ccv.json` file will be provided in the testnets repo. Replace the old `genesis.json` in the `$HOME/.onex/config` directory with the new `genesis-with-ccv.json`. The new `genesis-with-ccv.json` file with ccv data will be published in [onomyprotocol/valiadtor](https://github.com/onomyprotocol/validator/tree/main/testnet/onex-mainnet-1) |
|5   |Genesis reached     | Start your node with the consumer binary | onex-mainnet-1 chain will start and become a consumer chain.                                                                                     |
|6   |3 blocks after upgrade height                     |Celebrate! :tada: 🥂                                                |<chain> blocks are now produced by the provider validator set|
