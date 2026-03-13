# E-commerce Backend (2026.03)

基于Go+Gin+GORM+MySql+JWT的电商后端系统，实现了完整的购物流程。

## 技术栈
- Go + Gin框架
- MySQL + GORM
- JWT认证
- API-Service-DAO分层架构

## 功能模块
- ✅ 用户认证 (JWT)
- ✅ 商品管理 (列表/搜索/详情)
- ✅ 购物车 (添加/查看/删除)
- ✅ 订单系统 (下单/列表/详情)

## 项目亮点

### 1. 原子函数设计
将复杂业务逻辑拆分成小函数，提高可测试性和可维护性：

``` go
// 下单流程拆分成原子函数
func MakeOrder(userID int64, address string) (int64, error)
func getCartItems(userID int64) ([]domain.Cart, error)
func createOrder(userID int64, address string, carts []domain.Cart) (int64, error)
func processCartItem(orderID int64, cart domain.Cart) error
func checkStock(product *domain.Product, quantity int) error
```

### 2. 分层架构
- **API层**: HTTP请求处理
- **Service层**: 业务逻辑实现  
- **DAO层**: 数据库操作封装
- **Domain层**: 数据模型定义

### 3. 错误处理标准化
``` go
type BusinessError struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}
```

## 数据库设计
``` go
// 订单表
type Order struct {
    OrderID   int64           `gorm:"primaryKey"`
    UserID    int64           `gorm:"index"`
    Address   string          `gorm:"not null"`
    Total     decimal.Decimal `gorm:"type:decimal(10,2)"`
    Status    string          `gorm:"default:'pending'"`
    CreatedAt time.Time       `gorm:"autoCreateTime"`
}

// 订单商品表
type OrderItem struct {
    OrderItemID int64           `gorm:"primaryKey"`
    OrderID     int64           `gorm:"index"`
    ProductID   int64           
    ProductName string          `gorm:"not null"`
    Quantity    int             `gorm:"not null"`
    Price       decimal.Decimal `gorm:"type:decimal(10,2)"`
}
```

## API接口
``` go
POST /order/create     # 下单
GET  /order/list      # 订单列表
GET  /order/info/:id  # 订单详情

POST /cart/add/:id    # 加入购物车
GET  /cart/list       # 购物车列表
DELETE /cart/remove/:id # 移除购物车

GET  /product/list    # 商品列表
POST /product/search  # 搜索商品
GET  /product/info/:id # 商品详情
```

## 性能压测结果

### 基准测试（无Redis）

#### 2026.03.14 - 商品列表接口（读操作）
- **工具**: hey
- **并发**: 50
- **时长**: 30秒
- **QPS**: 10,196 请求/秒
- **平均延迟**: 4.9ms
- **错误率**: 0%
- **备注**: 纯MySQL版本，读操作性能优秀

#### 2026.03.14 - 订单创建接口（写操作）
- **工具**: hey
- **并发**: 100
- **时长**: 30秒
- **QPS**: 完全失败
- **平均延迟**: 超时/失败
- **错误率**: 100%
- **备注**: 写操作存在严重并发冲突，需要优化

### 性能分析结论
1. ✅ **读操作性能极佳** - 纯MySQL即可达到10K+ QPS
2. ❌ **写操作严重瓶颈** - 高并发下完全失败
3. 🔧 **优化方向** - 需要添加乐观锁、事务管理等并发控制机制

## 学习收获
1. **电商系统架构** - 完整购物流程实现
2. **Go开发技能** - Gin+GORM+MySql+JWT 实战经验
3. **代码质量** - 原子函数设计思想
4. **工程化思维** - 分层架构和错误处理

## 经验总结与吐槽 💡
- **分层架构的思想很重要**，现在重构的项目就是从这个思路出发的，由**api层、service层、dao层、domain层**组成。大一下写这个项目的时候，把这四层揉到一起了，业务逻辑极其混乱，导致我在初期重读代码的时候特别头疼。后面废了大力气，才把业务逻辑拆分成不同的层，代码结构变得清晰很多。
---
-   最开始在写登录接口的时候，**请求和响应**没有**分开处理**，也就是说，我前端传来请求，我直接bind到user结构体，而不是registerRequest。这导致了个什么后果呢？这导致我返回响应时，把user的password也返回了。这是个严重的安全漏洞。于是我就推倒重写，为**请求和响应分别设计结构体**，这样就解决了传入和返回不必要的字段的问题。我觉得这个思路太巧妙了哈哈哈。
---
- 我最开始写的时候没有注意**错误处理**，又是 errors.New 的，又是 fmt.Errorf 的，格式特别混乱。而且很多时候错误处理形同虚设，就是报错了我也找不出是哪里的问题。后来我就设计了一个BusinessError结构体，用来统一错误处理，然后const了很多错误码，这样返回的错误既清晰又统一，也方便后续维护和拓展，感觉这个设计很不错。
---
- 大一写项目的时候，根本就不懂**原子函数**，把所有业务逻辑都写到一起，比如说service层的login函数，里面包含了用户验证、token生成、数据库操作等多个步骤（也如第一点吐槽说的，api，service，dao层逻辑叠到一个函数写，你敢信吗，我api层直接调用dao层函数，写的简直是狗屎），这导致了大量的代码耦合，也不利于后续的维护和拓展。现在重构项目，我就把业务逻辑拆分成小函数，每个函数只负责一件事，这样就提高了代码的可维护性和复用性。有那么一小会儿，觉得自己像架构师哈哈哈。
---
- 这次重构项目的时候，我写了大量的**通用函数**，比如通用插入，通用查询，通用删除等等。之前学概念的时候，对**泛型**的理解一直不太深刻。但这次写通用函数，就发现泛型的优势很大。比如我写一个通用的插入函数，我可以传入任意结构体，它就会根据结构体的字段，自动生成对应的SQL语句。这就避免了我每次写插入函数都要写一遍SQL语句的问题，也提高了代码的复用性。
---
- **代码命名一定一定要规范！**，变量名、函数名、结构体名等都要符合命名规范。我大一写的时候，那几个ID（userID，productID等），全部没做区分，直接用ID来表示，导致代码可读性非常糟糕，看半天才理解业务逻辑。现在重构项目的时候，就对这部分做了优化。
---
- 写注释真是个好习惯啊。大一写的代码，有部分注释了有部分没有。我读代码的时候，明显感觉到有注释的好读得多，即使过了一年时间，也能快速上手。相对的，那些没注释过的代码（加上我命名也不规范），读起来就会头疼很多。
---
- **一定一定要测试**!当初大一写的代码，我现在测试根本跑不通，那个时候觉得不报错就是对的，现在看来荒唐无比。这次重构项目就是完成一个板块，就测试一个板块。比如我写完商户接口，就直接测试商户接口，然后再写其他接口，保证之前的接口不会出错。有好多问题goland是根本不会报错的，甚至服务端也能200，但是业务逻辑就是有问题。我印象最深的就是，商品重名，会给商户返回未知错误（正常情况会返回商品重名的友好提示）。这种问题不在apifox上测试，根本不可能发现，这个错误的原因是，我有个地方忘记包装bizErr了，是直接返回的原生err，导致我api有个区分业务错误和原生错误的错误鉴定函数错误判断，把商品名重复当成原生错误，最后生成了个未知错误返回。所以测试真的是重中之重。
---

> 这个项目让我深刻体会到：**代码不只是让机器运行，更是给人看的**。好的代码架构和规范，能让后续维护轻松很多！

## 后续优化
- [ ] Redis 集成
- [ ] Docker 部署
- [ ] MQ 消峰