# Api

## Static

| Method | URI                 | Desc                           |
| ------ | ------------------- | ------------------------------ |
| `GET`  | `/`                 | 获取主页                       |
| `GET`  | `/static/*path`     | 获取静态资源                   |
| `GET`  | `/connec/:resource` | 使用第三方来拦截，默认在配置中 |

## API 接口

| Method | URI                              | Desc                                                                             |
| ------ | -------------------------------- | -------------------------------------------------------------------------------- |
| `GET`  | `/api/info`                      | 获取 PGWEB 信息                                                                  |
| `POST` | `/api/connect`                   | 连接 Postgres                                                                    |
| `POST` | `/api/disconnect`                | 断开链接                                                                         |
| `POST` | `/api/switchdb`                  | 切换数据库                                                                       |
| `GET`  | `/api/databases`                 | 获取数据库                                                                       |
| `GET`  | `/api/connection`                | 获取链接                                                                         |
| `GET`  | `/api/server_settings`           | 获取服务端配置                                                                   |
| `GET`  | `/api/activity`                  | 获取当前活跃的查询                                                               |
| `GET`  | `/api/schemas`                   | 获取 schemas                                                                     |
| `GET`  | `/api/objects`                   | 获取 页面左侧的对象                                                              |
| `GET`  | `/api/tables/:table`             | 获取 获取指定表的结构信息，支持 materialized_view/function/table，以表格形式返回 |
| `GET`  | `/api/tables/:table/rows`        | 获取 获取指定表的行记录，以表格形式返回                                          |
| `GET`  | `/api/tables/:table/info`        | 获取 表的信息                                                                    |
| `GET`  | `/api/tables/:table/indexes`     | 获取 表的索引，以表格形式返回                                                    |
| `GET`  | `/api/tables/:table/constraints` | 获取 表的约束，以表格形式返回                                                    |
| `GET`  | `/api/tables/:table/constraints` | 获取 表的约束，以表格形式返回                                                    |
| `GET`  | `/api/table_stats`               | 获取 表的可导出信息，支持 json/xml/csv 格式                                      |
| `GET`  | `/api/functions/:id`             | 获取 函数详情                                                                    |
| `GET`  | `/api/query`                     | 执行查询                                                                         |
| `POST` | `/api/query`                     | 执行查询                                                                         |
| `GET`  | `/api/explain`                   | 执行解释                                                                         |
| `POST` | `/api/explain`                   | 执行解释                                                                         |
| `GET`  | `/api/analyze`                   | 执行分析                                                                         |
| `POST` | `/api/analyze`                   | 执行分析                                                                         |
| `GET`  | `/api/history`                   | 获取历史                                                                         |
| `GET`  | `/api/bookmarks`                 | 获取书签                                                                         |
| `GET`  | `/api/export`                    | 导出数据                                                                         |
| `GET`  | `/api/local_queries`             | 获取本地查询列表                                                                 |
| `GET`  | `/api/local_queries/:id`         | 执行本地查询                                                                     |
| `POST` | `/api/local_queries/:id`         | 执行本地查询                                                                     |

## Metric

| Method | URI        | Desc                                                           |
| ------ | ---------- | -------------------------------------------------------------- |
| `GET`  | `/metrics` | 获取性能指标，仅当使用 `--metrics`，且未配置远程 metric server |
