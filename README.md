# 🎯 项目收获与难点总结

> ✨ 这里主要写一下我写这个项目的主要收获，包括项目的经典技术实现和实现项目时遇到的困难。

---

## 📐 一、项目架构与核心技术

### 1️⃣ 首先是项目的架构设计：

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
┌─────────────────────────────────┐
│ Controller 层 (控制器层) │ ← 处理 HTTP 请求/响应，参数校验
├─────────────────────────────────┤
│ Logic 层 (业务逻辑层) │ ← 核心业务逻辑编排
├─────────────────────────────────┤
│ DAO 层 (数据访问层) │ ← MySQL + Redis 数据操作
├─────────────────────────────────┤
│ Model 层 (数据模型层) │ ← 数据结构定义
└─────────────────────────────────┘
</pre>

这种分层理念提高了代码的复用行和开发效率。

### 2️⃣ 重点是 Gin 中间件链式处理思路：

利用中间件可以控制在访问一些路由之前先执行其他函数。  
比如一些路由需要登录之后才能访问，就可以在这些路由组之前通过 `Use` 请求者身份校验函数来实现：

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
// 应用 JWT 中间件到这个组
v1.Use(middlewares.JWTAuthMiddleware())
// 这个大括号内的所有路由都需要 JWT 认证
{
v1.GET("/community", controller.CommunityHandler) // 🔒 需要登录
v1.GET("/community/:id", controller.CommunityDetailHandler) // 🔒 需要登录
v1.POST("/post", controller.CreatePostHandler) // 🔒 需要登录
v1.GET("/post/:id", controller.GetPostDetailHandler) // 🔒 需要登录
v1.GET("/posts", controller.GetPostListHandler) // 🔒 需要登录
v1.GET("/posts2", controller.GetPostListHandler2) // 🔒 需要登录
v1.POST("/vote", controller.PostVoteController) // 🔒 需要登录
}
</pre>

再比如我们可以通过使用全局中间件实现记录日志、捕获 panic（防止服务器崩溃）这种各个业务都需要的函数：

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
r := gin.New()  // 注意：这里用 New() 而不是 Default()
r.Use(logger.GinLogger(), logger.GinRecovery(true))
</pre>

另外中间件内可以通过 `Next()` 实现先执行下一层的代码，`Abort()` 实现直接中断不再执行。

---

## ⚙️ 二、关键技术实现

### 1️⃣ JWT（JSON Web Token）无状态认证

JWT 是一种基于 Token 的认证机制，它将用户信息以加密的 JSON 格式存储在客户端，服务端不需要保存任何会话状态。  
（以前的技术如 cookie 需要传递一个 cookie 给服务端用来保存用户信息实现登录）

**工作流程**

a. 登录阶段  
   用户提交用户名和密码到服务器  
   服务器验证用户身份  
   验证通过后，使用密钥生成 JWT Token  
   将 Token 返回给客户端  

b. Token 结构  
   JWT 由三部分组成，用 `.` 分隔：  
   Header（头部）：指定算法和 Token 类型  
   Payload（载荷）：包含用户信息和元数据（如用户 ID、过期时间等）  
   Signature（签名）：对前两部分的签名，防止篡改  

c. 请求认证  
   客户端在后续请求中携带 Token（通常在 Authorization Header 中）  
   格式：`Authorization: Bearer <token>`  

d. 服务端验证  
   服务端接收到 Token 后，使用相同的密钥进行验证  
   检查签名是否有效  
   检查 Token 是否过期  
   从 Payload 中提取用户信息  
   无需查询数据库或缓存  

**注意事项**  
- Token 撤销困难：一旦签发，在过期前无法主动撤销（除非使用黑名单机制）  
- 安全性：Payload 只是 Base64 编码，不是加密，不要存放敏感信息  
- Token 大小：比 Session ID 大，会增加网络传输开销  
- 密钥管理：签名密钥必须妥善保管  

**补充**  
- 通过 token 比对实现单机登录：登录时将 token 存入 redis，中间件验证 Redis 中的 token 是否匹配，不匹配说明有其他设备登录。  
- 通过设置 token 过期时间可以实现需要重新登录的时间。

### 2️⃣ Redis 实现对帖子投票和查询

利用 redis 里的有序集合可以快速实现按时间或者按分数排序。  
利用分数 `// 每票 432 分 = 60 * 60 * 24 / (7 * 24 * 60 * 60) * 1000000` 实现帖子热度概念  
（赞成票（1 分）→ 增加 432 分；反对票（-1 分）→ 减少 432 分；改变投票 → 只计算差值；取消投票 → 减去全部分数）

**重点思想：**

因为同时间内可能有大量人投票，反复多次访问修改数据库网络 RTT 高，可以使用 pipeline 优化：

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
pipeline := client.Pipeline()
for _, id := range ids {
    pipeline.ZCount(key, "1", "1")
}
cmders, err := pipeline.Exec()  // 一次网络往返
</pre>

即先把请求存入 pipeline 再一次发送给数据库

**进一步优化：** 因为一次投票包括对数据库的两次改动：一次是修改帖子分数，一次是修改投票用户投票记录；  
这两个事件不分先后同时出现，很容易想到可以绑定成一个事务保证原子性：

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
pipeline := client.TxPipeline()
pipeline.ZAdd(...)
pipeline.SAdd(...)
_, err := pipeline.Exec()  // 要么都成功，要么都失败
</pre>

### 3️⃣ 雪花算法分布式 ID

用 Mysql 自增 id 会暴露客户量等，也容易被遍历造成安全问题。

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
// 初始化：设置起始时间和机器 ID
sf.Epoch = st.UnixNano() / 1000000
node, _ = sf.NewNode(machineID)
// 生成唯一 ID
id := node.Generate().Int64()
</pre>

雪花算法思路是利用时间和机器 id 生成唯一的随机 id（时间戳 + 机器 ID + 序列号的组合设计）

### 4️⃣ MySQL 建表管理

利用 gorm 和 tag 快速管理 Mysql 数据实例：

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
// Tag 映射
type Post struct {
    ID          int64     `gorm:"primary_key" json:"id"`
    Title       string    `gorm:"size:128" json:"title"`
    CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
}
</pre>

### 5️⃣ 统一错误码设计

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
type ResCode int64
const (
    CodeSuccess ResCode = 1000 + iota
    CodeInvalidParam
    CodeUserExist
    CodeNeedLogin
    CodeInvalidToken
)
// 错误码映射消息
var codeMsgMap = map[ResCode]string{
    CodeSuccess: "success",
    CodeInvalidParam: "请求参数错误",
}
</pre>

再通过统一的响应格式封装：

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
ResponseSuccess(c, data)  // {"code":1000, "msg":"success", "data":...}
ResponseError(c, CodeNeedLogin)  // {"code":1014, "msg":"需要登录", "data":null}
</pre>

这样实现前后端交互也方便问题定位

### 6️⃣ Air 与 Makefile

air 可以实现项目的热重载，Makefile 则实现编译的打包；  
注意 air 可能会占领终端，所以错误最好写入日志而不是打印在终端：

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
zap.L().Error("mysql.GetPostById failed",
    zap.Int64("pid", pid),
    zap.Error(err))

zap.L().Debug("checking redis token",
    zap.Int64("userID", mc.UserID))
</pre>

### 7️⃣ Docker 部署

利用 docker 容器化项目可以方便部署服务器和生成镜像

### 8️⃣ Swagger 文档

按照 swagger 定义给接口写注释可以同步生成带测试的接口文档。

---

## 🚧 三、项目开发时遇到的问题

1. 后端传递给前端的 json 包含 int64 的变量如用户 id 时，由于前端 JavaScript 文件不含 int64 类型会造成精度丢失，问题难以发现和修改。  
解决方案：用 string 传递，可以利用 tag 快速转化：

<pre style="font-family: monospace; background-color: #f0f4f8; padding: 16px; border-radius: 8px; border: 1px solid #d0d7de; line-height: 1.45; white-space: pre-wrap; word-wrap: break-word;">
ID int64 `json:"id,string"`    //json传给前端时候转化
ID int64 `db:"post_id,string"` //或者传给数据库的时候转化
</pre>
