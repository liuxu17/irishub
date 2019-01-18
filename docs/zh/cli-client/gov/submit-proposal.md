# iriscli gov submit-proposal

## 描述

提交区块链治理提议以及发起提议所涉及的初始保证金，其中提议的类型包括Text/ParameterChange/SoftwareUpgrade这三种类型。

## 使用方式

```
iriscli gov submit-proposal [flags]
```
打印帮助信息:

```
iriscli gov submit-proposal --help
```
## 标志

| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | [string] 提议的保证金（至少30%minDeposit）                                                                                                                         |          |
| --description    |                            | [string] 提议的描述                                                                                                           | Yes      |
| --key            |                            | 参数的键名称                                                                                                                        |          |
| --op             |                            | [string] 对参数的操作                                                                                                             |          |
| --param          |                            | [string] 提议参数,例如: [{key:key,value:value,op:update}]                                                                                 |          |
| --path           |                            | [string] param.json文件路径                                                                                                                      |          |
| --title          |                            | [string] 提议标题                                                                                                                           | Yes      |
| --type           |                            | [string] 提议类型,例如:ParameterChange/SoftwareUpgrade/SoftwareHalt/TxTaxUsage                                                                            | Yes      |
| --version           |            0                | [uint64] 新协议的版本信息                                                                           |       |
| --software           |           " "                 | [string] 新协议的软件地址                                                                       |       |
| --switch-height           |       0                     | [string] 新版本协议升级的高度                                                     |       |
| --threshold        | "0.85"   |  [string] 软件升级的阈值                                              |               |

## 例子

提议者必须抵押至少30%的`MinDeposit`，详情见 [Gov](../../feature/governance.md)

### 提交一个'ParameterChange'类型的提议

```shell
iriscli gov submit-proposal --chain-id=test --title="update MinDeposit proposal" --param='mint/Inflation=0.050' --type=ParameterChange --description="a new parameter change proposal" --from=node0 --fee=0.01iris
```

提交之后，您完成了提交新的“ParameterChange”提议。
更改参数的详细信息（通过查询参数获取参数，修改它，然后在“操作”上添加“更新”，使用方案中的更多详细信息）和其他类型的提议字段与文本提议类似。
注意：在这个例子中, --path 和 --param 不能同时为空。

### 提交一个'SoftwareUpgrade'类型的提议

```shell
iriscli gov submit-proposal --chain-id=test --title="irishub0.7.0 upgrade proposal" --type=SoftwareUpgrade --description="a new software upgrade proposal" --from=node0 --fee=0.01iris --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80 --threshold=0.9
```

在这种场景下，提议的 --title、--type 和--description参数必不可少，另外你也应该保留好提议ID，这是检索所提交提议的唯一方法。


如何查询提议详情？

请点击下述链接：

[query-proposal](query-proposal.md)

[query-proposals](query-proposals.md)
