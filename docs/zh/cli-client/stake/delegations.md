# iriscli stake delegations

## 描述

查询某个委托者发起的所有委托记录

## 用法

```
iriscli stake delegations [delegator-address] [flags]
```
打印帮助信息
```
iriscli stake delegations --help
```

## 示例

### 查询某个委托者发起的所有委托记录

```
iriscli stake delegations [delegator-address]
```

运行成功以后，返回的结果如下：

```json
[
  {
    "delegator_addr": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
    "validator_addr": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
    "shares": "0.2000000000",
    "height": "290"
  }
]
```
