# iriscli bank

## 描述

Bank模块允许你管理你本地账户的资产。 

## 使用方式

```
 iriscli bank [command]
```

 

## 相关命令

| 命令      | 描述                   |
| --------- | ---------------------- |
| [coin-type](coin-type.md) | 查询货币类型           |
| [token-stats](token-stats.md) | 查询token统计数据           |
| [account](account.md)   | 查询账户余额           |
| [send](send.md)      | 创建和签名一个转账请求 |
| [sign](sign.md)      | 签名离线传输文件       |

## 标志

| 命令，速记 | 默认值 | 描述         | 是否必须 |
| ---------- | ------ | ------------ | -------- |
| -h, --help |        | Bank模块帮助 | 否       |

## 全局标志

| 命令，速记            | 默认值         | 描述                                | 是否必须 |
| --------------------- | -------------- | ----------------------------------- | -------- |
| -e, --encoding string | hex            | 字符串二进制编码 (hex \|b64 \|btc ) | 否       |
| --home string         | /root/.iriscli | 配置和数据存储目录                  | 否       |
| -o, --output string   | text           | 输出格式 (text \|json)              | 否       |
| --trace               |                | 出错时打印完整栈信息                | 否       |