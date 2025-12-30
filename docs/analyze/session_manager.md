# Session Manager

被 `dbCheck` 中间件所使用，检查数据库连接是否过期。

## 多 session 模式的管理

底层使用 `map[string]*client.Client` 来管理不会会话的客户端。
每隔会话单独一个客户端。

为了确保并发操作的安全性，搭配 `sync.Mutex` 的锁机制，当操作 map 时，会获取锁，完成后释放锁。

### Add

- 1.获取锁
- 2.写入 map
- 3.释放锁

### Remove

- 1.获取锁
- 2.关闭连接
- 3.从 map 中删除
- 4.释放锁

### Get

- 1.获取锁
- 2.从 map 中获取
- 3.释放锁

### Len

- 1.获取锁
- 2.计算 map 长度
- 3.释放锁

### Sessions

- 1.获取锁
- 2.创建 map 的拷贝
- 3.释放锁

### IDs

- 1.获取锁
- 2.迭代 map 的 key
- 3.释放锁

## 多 session 的清理

PGWEB 提供了多个 session 的清理功能，对于一段时间不使用的连接，将会使用 goroutine 进行清理。该会话下次尝试使用时，会因为 `dbCheck` 而自动重定向到首页。

### 清理相关配置

- `--no-idle-timeout` 禁用清理，则连接永不过期，不再释放
- `idle-timeout` 存活时间，用当前时间-连接上次使用时间 > idle-timeout，满足清理条件

### 清理逻辑

当开启了清理功能后，在创建 `NewSessionManager` 之后，使用 goroutine 独立运行 `RunPeriodicCleanup` 来执行清理任务。

- 1.使用 Ticker 来实现每分钟执行一次检查
- 2.计算 now - 连接的 lastQueryTime > idleTimeout，则满足
- 3.执行清理，从 map 中移除

- 1.在 conn 每次执行查询的时候，即调用 `query` 时，更新 `lastQueryTime`，代表其被使用了一次。
