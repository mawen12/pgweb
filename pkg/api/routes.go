package api

import (
	"github.com/gin-gonic/gin"

	"github.com/sosedoff/pgweb/pkg/command"
	"github.com/sosedoff/pgweb/pkg/metrics"
)

// 配置中间件
func SetupMiddlewares(group *gin.RouterGroup) {
	if command.Opts.Cors {
		group.Use(corsMiddleware())
	}

	group.Use(dbCheckMiddleware())
}

// 配置路由
func SetupRoutes(router *gin.Engine) {
	// 添加统一前缀
	root := router.Group(command.Opts.Prefix)

	// / => 获取主页
	root.GET("/", gin.WrapH(GetHome(command.Opts.Prefix)))
	// /static/*path => 获取主页所需的静态资源
	root.GET("/static/*path", gin.WrapH(GetAssets(command.Opts.Prefix)))
	// /connnect/:resource => 连接后端
	root.GET("/connect/:resource", ConnectWithBackend)

	// /api 接口
	api := root.Group("/api")
	// 配置中间件
	SetupMiddlewares(api)

	// 开启多会话时，配置 /api/sessions
	if command.Opts.Sessions {
		api.GET("/sessions", GetSessions)
	}

	// /api/info => PGWEB 应用信息
	api.GET("/info", GetInfo)
	// /api/connect => 常见一个新连接
	api.POST("/connect", Connect)
	// /api/disconnect => 断开连接
	api.POST("/disconnect", Disconnect)
	// /api/switchdb => 切换数据库
	api.POST("/switchdb", SwitchDb)
	// /api/databases => 列举数据库
	api.GET("/databases", GetDatabases)
	// /api/connection => 获取连接信息
	api.GET("/connection", GetConnectionInfo)
	// /api/server_settings => 获取服务设置
	api.GET("/server_settings", GetServerSettings)
	// /api/activity => 获取当前活跃的查询
	api.GET("/activity", GetActivity)
	// /api/schemas => 获取 schema
	api.GET("/schemas", GetSchemas)
	// /api/objects => 获取对象
	api.GET("/objects", GetObjects)
	// /api/tables/:table => 获取指定表
	api.GET("/tables/:table", GetTable)
	// /api/tables/:table/rows => 获取指定表的行记录
	api.GET("/tables/:table/rows", GetTableRows)
	// /api/tables/:table/info => 获取表信息
	api.GET("/tables/:table/info", GetTableInfo)
	// /api/tables/:table/indexes => 获取表索引
	api.GET("/tables/:table/indexes", GetTableIndexes)
	// /api/tables/:table/constraints => 获取表约束
	api.GET("/tables/:table/constraints", GetTableConstraints)
	// /api/tables_stats => 获取表统计数据
	api.GET("/tables_stats", GetTablesStats)
	// /api/functions/:id => 获取函数
	api.GET("/functions/:id", GetFunction)
	// /api/query => 执行查询，GET / POST
	api.GET("/query", RunQuery)
	api.POST("/query", RunQuery)
	// /api/explain => 执行解释，GET / POST
	api.GET("/explain", ExplainQuery)
	api.POST("/explain", ExplainQuery)
	// /api/analyze => 执行分析，GET / POST
	api.GET("/analyze", AnalyzeQuery)
	api.POST("/analyze", AnalyzeQuery)
	// /api/history => 获取SQL查询历史
	api.GET("/history", GetHistory)
	// /api/bookmarks => 获取书签
	api.GET("/bookmarks", GetBookmarks)
	// /api/export => 导出
	api.GET("/export", DataExport)
	// /api/local_queries => 获取本地查询
	api.GET("/local_queries", requireLocalQueries(), GetLocalQueries)
	// /api/local_queries/:id => 获取本地查询，GET / POST
	api.GET("/local_queries/:id", requireLocalQueries(), RunLocalQuery)
	api.POST("/local_queries/:id", requireLocalQueries(), RunLocalQuery)
}

// 配置指标
func SetupMetrics(engine *gin.Engine) {
	if command.Opts.MetricsEnabled && command.Opts.MetricsAddr == "" {
		// NOTE: We're not supporting the MetricsPath CLI option here to avoid the route conflicts.
		engine.GET("/metrics", gin.WrapH(metrics.NewHandler()))
	}
}
