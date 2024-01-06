# crypto

Simple encryption algorithm, for utf-8.

一个简单的加密工具，支持unicode；

## 核心算法: 凯撒算法-Caesar cipher

### 加密:

明文字符(字符串) -->  转化为unicode码 --> 编码偏移加密 --> 编码为字符 --> base64编码 --> 打乱顺序(前后双指针方式) = 密文

### 解密:

密文字符 --> 恢复顺序 --> 得到base64编码 --> 解码为字符串 --> 转化为unicode码 --> 编码偏移解密 --> 解析unicode  = 明文

### 当前已支持的语言:

+ php
+ golang
+ javascript
+ typescript
