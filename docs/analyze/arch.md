# Arch

## 后端

- `gin` web 通信框架
- `jmoiron/sqlx` 来执行查询 SQL

- `tuvistavie/securerandom` 生成安全的 session id
- `mitchellh/go-homedir` 获取用户目录
- `jessevdk/go-flags` 解析命令行参数

###

- api 层 => gin + middleware + routes + sessionManager
  - `client.Client`
  - `SessionManager`
    - `client.Client`
  - `queries.Store`
- client 层 => Client + bookmarks + history + sql
  - `sqlx.DB`
  - `history.Record`
- db 层 => sqlx

- 启动 => cli + command

## Logic

- 1.pkg/command
- 2.pkg/client
- 3.pkg/cli

## 前端

- `bootstrap` 样式管理
- `font-awesome` 字体

### 前端布局

- `main`

  - `nav` 顶部的选项管理，功能有：`Rows`, `Structure`, `Indexes`, `Constraints`, `Query`, `History`, `Activity`, `Connection`, `Edit Connnection`, `Close Connection`
  - `sidebar` 展示当前连接的 Database， 搜索区域，Tables 列表，Table Info 区域
  - `body` 根据顶部不同的选项展示不同的内容
    - `input` 输入
    - `output` 输出
    - `pagination` 分页，对于表格的数据，使用分页进行处理

- `content-modal` 双击表内容展示
- `connection_widow` 连接页面，支持 Scheme, Standard, SSH
- `tables_context_menu` Table 支持的操作列表
- `view_context_menu` Vie 支持的操作列表
- `current_database_context_menu` 当前数据库支持的操作列表
- `results_header_menu` 当前 header 支持的操作列表
- `results_row_menu` 当前 row 支持的操作列表
- `error_banner` 报错 banner
