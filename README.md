# 项目说明
接口自动化项目，用于流程化测试（主流程或者冒烟测试），接口定时巡检等
开发语言 golang

# 目录说明
```
/conf   //配置文件解析
/exec  // 执行流程文件
/flow_test  // 流程测试文件
/logger  // 日志记录包
/logs    // 测试日志记录（内含接口相应日志）
/templates  // 测试用例编写目录（主要）
/tool  // 工具包
```

# 构建执行文件
```
// windows 
GOOS=windows GOARCH=amd64 go build -o test_api.exe .
// linux 
GOOS=linux go build -o test_api_linux .
// mac
GOOS=darwin go build -o test_api_darwin .
```

# 运行
```
// windows
./test_api.exe ./flow.test.yaml  // flow.test.yaml 是要执行的流程测试文件，内含需要配置的测试用例文件
// linux 
./test_api_linux ./flow.test.yaml  // flow.test.yaml 是要执行的流程测试文件，内含需要配置的测试用例文件
// mac
./test_api_darwin ./flow.test.yaml  // flow.test.yaml 是要执行的流程测试文件，内含需要配置的测试用例文件
```

# 文件说明
### 共用配置文件
```
# 固定参数 数据格式只能是数字和字符串
api_domain: http://test.local.com  // 接口请求域名
timeout: 5  // 接口请求超时设置 （秒）
feishu_robot: https://open.feishu.cn/open-apis/bot/v2/hook/xxxx-xxx-xxxx  // 飞书机器人webhook地址

# 自定义参数 数据格式只能是数字和字符串
token: a112324
```
### 执行文件说明
```
include: 
  - ./templates/user/login.test.yaml  // 测试用例文件有前后次序的一定要按顺序
  - ./templates/user/menuList.test.yaml
```

### 测试用例文件
```
api_domain: http://test.local.com // 接口域名，可选，优先级最高，没有配置用配置文件里面的配置
api: /User/AdminUser/operationLogin  // 接口路径
method: POST                // 请求方式 目前仅支持GET  POST
headers:
  Authorization: "Bearer $cache.token"
params:         // 请求参数
  email: test@admin.com
  password: 123456

expect:     // 预期值设置
  assertEquals:  // 断言值相等
    code: 0
    message: "登录成功"
  assertNotEquals:  // 断言值不相等
    data.token: ""
    data.user_id: 0

//  AssertContains  断言值包含
//  AssertMatches  断言正则表达式
//  AssertLength  断言长度
//  AssertType  断言类型
//  AssertTrue  断言真
//  AssertFalse  断言假
//  AssertIn  断言在数组内

cache:   // 缓存数据 以便后面接口使用
  token: data.token
  user_id: data.user_id
```

缓存数据/全局变量数据使用方式
```
Authorization: "Bearer $cache.token"  // 使用cache变量 $cache.+变量名
Authorization: "Bearer $token"  // 使用全局变量 config.yaml内配置 $+变量
```

比较中含有接口获取的值
```
dataValue(data.list.0.id)   // dataValue(path)  path是值的层级
```

# 支持的方法
### 方法列表
```
RandInt(min,max)        // min,max 数字整型， 随机一个介于min,max数字
Time()                         // 当前时间戳
TimeMillis()                 // 当前数据毫秒
TimeMicros()               // 当前时间微秒
TimeNanos()                // 当前时间纳秒数
Date()                         // 当前日期 格式 2026-01-26
DateTime()                 // 当前时间 格式 2026-01-26 19:48:23
RandString(flag,length)  // 随机字符串 flag标识 N代表含大写字母 n代表含小写字母 1代表含数字 S代表含特殊字符   length 长度  调用方式fn.RandString("Nn1",10)
```

### 使用方式
```
headers:
  Content-Type: application/json
  Timestamp: fn.DateTime()
params:
  email: test@admin.com
  password: 123456
  timestamp: fn.Time()
  nonce: fn.RandString("Nn1",10)@aaaa
```

### 计算表达式
```
expr(5*7/8)
```