# iriscli bank sign

## Description

Sign transactions in generated offline file. The file created with the --generate-only flag.

## Usage:

```
iriscli bank sign <file> [flags]
```


## Global Flags

| Name,shorthand        | Default        | Description                                 | Required | Type   |
| --------------------- | -------------- | ------------------------------------------- | -------- | ------ |
| -e, --encoding string | hex            | String   Binary encoding (hex \|b64 \|btc ) | False    | String |
| --home string         | /root/.iriscli | Directory for config and data               | False    | String |
| -o, --output string   | text           | Output format (text \|json)                 | False    | String |
| --trace               |                | Print out full stack trace on errors        | False    |        |

## Examples

### Sign a send file 

First you must use **iriscli bank send**  command with flag **--generate-only** to generate a send recorder. Just like this.

```shell  
iriscli bank send --to=faa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429yx  --from=test  --fee=0.004iris --chain-id=test-irishub--amount=10iris --generate-only > tx.json
```


And then save the output in file `tx.json`.

Then you can sign the offline file.

```
iriscli bank sign tx.json --name=test  --offline=true --print-sigs=false --append=true
```

After that, you will get the detail info for the sign. Like the follow output you will see the signature 

**ci+5QuYUVcsARBQWyPGDgmTKYu/SRj6TpCGvrC7AE3REMVdqFGFK3hzlgIphzOocGmOIa/wicXGlMK2G89tPJg==**

```
iriscli bank sign /root/output/output/node0/test_send_10iris.txt --name=test  --offline=false --print-sigs=false --append=true
```

```json
{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"faa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"faa1893x4l2rdshytfzvfpduecpswz7qtpstpr9x4h","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"40000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"Auouudrg0P86v2kq2lykdr97AJYGHyD6BJXAQtjR1gzd"},"signature":"sJewd6lKjma49rAiGVfdT+V0YYerKNx6ZksdumVCvuItqGm24bEN9msh7IJ12Sil1lYjqQjdAcjVCX/77FKlIQ==","account_number":"0","sequence":"3"}],"memo":"test"}}```
```
After signing a transaction, it could be broadcast to the network with [broadcastc command](./broadcast.md)