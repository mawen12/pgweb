# Middleware

## dbCheck

用于校验数据库连接状态。
默认为单 session 模式，使用 `api.DbClient` 来管理连接，且永不过期。
当配置了`--sessions`时，切换到多 session 模式，使用 `SessionManager` 来管理连接。
如果同时配置了`--no-idle-timeout`，则其会话连接永不过期；反之通过`--idle-timeout`
来控制连接过期。

## Cors

允许跨域操作。
仅当配置了 `--cors` 时才起效。
