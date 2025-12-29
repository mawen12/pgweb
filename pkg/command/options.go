package command

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/jackc/pgpassfile"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

const (
	// Prefix to use for all pgweb env vars, ie PGWEB_HOST, PGWEB_PORT, etc
	// PGWEB 所使用的变量前缀
	envVarPrefix = "PGWEB_"
)

type Options struct {
	// 仅输出版本，然后退出
	Version bool `short:"v" long:"version" description:"Print version"`
	// 开启后，日志级别强制为 debug，每隔30s输出一次内存使用情况
	Debug bool `short:"d" long:"debug" description:"Enable debugging mode"`
	// 日志级别，默认为 info，支持 panic, fatal, error, warn, info, debug, trace
	LogLevel string `long:"log-level" description:"Logging level" default:"info"`
	// 日志格式，默认为 text，支持 text, json
	LogFormat        string `long:"log-format" description:"Logging output format" default:"text"`
	LogForwardedUser bool   `long:"log-forwarded-user" description:"Log user information available in X-Forwarded-User/Email headers"`
	URL              string `long:"url" description:"Database connection string"`
	Host             string `long:"host" description:"Server hostname or IP" default:"localhost"`
	Port             int    `long:"port" description:"Server port" default:"5432"`
	User             string `long:"user" description:"Database user"`
	Pass             string `long:"pass" description:"Password for user"`
	Passfile         string `long:"passfile" description:"Local passwords file location"`
	DbName           string `long:"db" description:"Database name"`
	SSLMode          string `long:"ssl" description:"SSL mode"`
	SSLRootCert      string `long:"ssl-rootcert" description:"SSL certificate authority file"`
	SSLCert          string `long:"ssl-cert" description:"SSL client certificate file"`
	SSLKey           string `long:"ssl-key" description:"SSL client certificate key file"`
	OpenTimeout      int    `long:"open-timeout" description:"Maximum wait time for connection, in seconds" default:"30"`
	RetryDelay       uint   `long:"open-retry-delay" description:"Number of seconds to wait before retrying the connection" default:"3"`
	RetryCount       uint   `long:"open-retry" description:"Number of times to retry establishing connection" default:"0"`
	HTTPHost         string `long:"bind" description:"HTTP server host" default:"localhost"`
	HTTPPort         uint   `long:"listen" description:"HTTP server listen port" default:"8081"`
	AuthUser         string `long:"auth-user" description:"HTTP basic auth user"`
	AuthPass         string `long:"auth-pass" description:"HTTP basic auth password"`
	SkipOpen         bool   `short:"s" long:"skip-open" description:"Skip browser open on start"`
	// 多会话模式标志，未开启时，使用 DBClient，否则使用SessionManager，会自动清理
	Sessions bool   `long:"sessions" description:"Enable multiple database sessions"`
	Prefix   string `long:"prefix" description:"Add a url prefix"`
	// 只读模式
	ReadOnly bool `long:"readonly" description:"Run database connection in readonly mode"`
	// 单数据库连接模式，不允许创建新连接
	LockSession bool `long:"lock-session" description:"Lock session to a single database connection"`
	// 开启书签模式
	Bookmark string `short:"b" long:"bookmark" description:"Bookmark to use for connection. Bookmark files are stored under $HOME/.pgweb/bookmarks/*.toml" default:""`
	// 书签目录
	BookmarksDir string `long:"bookmarks-dir" description:"Overrides default directory for bookmark files to search" default:""`
	// 仅支持从书签创建连接
	BookmarksOnly bool `long:"bookmarks-only" description:"Allow only connections from bookmarks"`
	// 本地查询目录
	QueriesDir        string `long:"queries-dir" description:"Overrides default directory for local queries"`
	DisablePrettyJSON bool   `long:"no-pretty-json" description:"Disable JSON formatting feature for result export"`
	DisableSSH        bool   `long:"no-ssh" description:"Disable database connections via SSH"`
	ConnectBackend    string `long:"connect-backend" description:"Enable database authentication through a third party backend"`
	ConnectToken      string `long:"connect-token" description:"Authentication token for the third-party connect backend"`
	ConnectHeaders    string `long:"connect-headers" description:"List of headers to pass to the connect backend"`
	// 禁用链接存储超时，则不再检测链接是否超时，应用于 session manager
	DisableConnectionIdleTimeout bool `long:"no-idle-timeout" description:"Disable connection idle timeout"`
	// 设置链接超时时间，默认 180m
	ConnectionIdleTimeout int `long:"idle-timeout" description:"Set connection idle timeout in minutes" default:"180"`
	// 设置查询超时时间，默认 300s
	QueryTimeout uint `long:"query-timeout" description:"Set global query execution timeout in seconds" default:"300"`
	// 跨域拦截
	Cors bool `long:"cors" description:"Enable Cross-Origin Resource Sharing (CORS)"`
	// 当开启跨域拦截时有效
	CorsOrigin string `long:"cors-origin" description:"Allowed CORS origins" default:"*"`
	// 二进制编码，默认为 none，支持 none, hex, base58, base64
	BinaryCodec string `long:"binary-codec" description:"Codec for binary data serialization, one of 'none', 'hex', 'base58', 'base64'" default:"none"`
	// 开启 metrics 标识
	MetricsEnabled bool   `long:"metrics" description:"Enable Prometheus metrics endpoint"`
	MetricsPath    string `long:"metrics-path" description:"Path prefix for Prometheus metrics endpoint" default:"/metrics"`
	MetricsAddr    string `long:"metrics-addr" description:"Listen host and port for Prometheus metrics server"`
}

var Opts Options

// ParseOptions returns a new options struct from the input arguments
func ParseOptions(args []string) (Options, error) {
	var opts = Options{}

	// 从命令行解析上述参数
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return opts, err
	}

	// 解析日志级别（作者自行开发的工具库）
	_, err = logrus.ParseLevel(opts.LogLevel)
	if err != nil {
		return opts, err
	}

	// 检查 url -> PGWEB_DATABASE_URL -> DATABASE_URL(输出警告，该版本已经弃用，但值会被保留)
	if opts.URL == "" {
		opts.URL = getPrefixedEnvVar("DATABASE_URL")
	}

	// 检查 prefix -> PGWEB_URL_PREFIX -> URL_PREFIX(输出警告，该版本已经弃用，但值会被保留)
	if opts.Prefix == "" {
		opts.Prefix = getPrefixedEnvVar("URL_PREFIX")
	}

	// 检查 passfile => PGPASSFILE => $HOME/.pgpass
	if opts.Passfile == "" {
		passfile := os.Getenv("PGPASSFILE")
		if passfile == "" {
			passfile = filepath.Join(os.Getenv("HOME"), ".pgpass")
		}

		_, err := os.Stat(passfile)
		if err == nil {
			_, err = pgpassfile.ReadPassfile(passfile)
			if err == nil {
				opts.Passfile = passfile
			} else {
				fmt.Printf("[WARN] Pgpass file unreadable: %s\n", err)
			}
		}
	}

	// Handle edge case where pgweb is started with a default host `localhost` and no user.
	// When user is not set the `lib/pq` connection will fail and cause pgweb's termination.
	// 检查 host 和 user，如果是本地且未指定user，则使用主机的用户名，否则置空
	if (opts.Host == "localhost" || opts.Host == "127.0.0.1") && opts.User == "" {
		if username := getCurrentUser(); username != "" {
			opts.User = username
		} else {
			opts.Host = ""
		}
	}

	// 使用 PGWEB_BOOKMARKS_ONLY 覆盖
	if getPrefixedEnvVar("BOOKMARKS_ONLY") != "" {
		opts.BookmarksOnly = true
	}

	// 使用 PGWEB_SESSIONS 覆盖
	if getPrefixedEnvVar("SESSIONS") != "" {
		opts.Sessions = true
	}

	// 使用 PGWEB_LOCK_SESSION 覆盖
	if getPrefixedEnvVar("LOCK_SESSION") != "" {
		opts.LockSession = true
		opts.Sessions = false
	}

	if opts.Sessions || opts.ConnectBackend != "" {
		opts.Bookmark = ""
		opts.URL = ""
		opts.Host = ""
		opts.User = ""
		opts.Pass = ""
		opts.DbName = ""
		opts.SSLMode = ""
	}

	// 格式化 prefix，确保其末尾带有 /
	if opts.Prefix != "" && !strings.HasSuffix(opts.Prefix, "/") {
		opts.Prefix = opts.Prefix + "/"
	}

	// 检查 auth-user => PGWEB_AUTH_USER 覆盖
	if opts.AuthUser == "" {
		opts.AuthUser = getPrefixedEnvVar("AUTH_USER")
	}

	// 检查 auth-pass => PGWEB_AUTH_PASS 覆盖
	if opts.AuthPass == "" {
		opts.AuthPass = getPrefixedEnvVar("AUTH_PASS")
	}

	// 检查 connect-backend 与 connect-token | connect-headers | sessions 的关联性
	if opts.ConnectBackend != "" {
		if !opts.Sessions {
			return opts, errors.New("--sessions flag must be set")
		}
		if opts.ConnectToken == "" {
			return opts, errors.New("--connect-token flag must be set")
		}
	} else {
		if opts.ConnectToken != "" || opts.ConnectHeaders != "" {
			return opts, errors.New("--connect-backend flag must be set")
		}
	}

	if opts.BookmarksOnly {
		if opts.URL != "" {
			return opts, errors.New("--url not supported in bookmarks-only mode")
		}
		if opts.Host != "" && opts.Host != "localhost" {
			return opts, errors.New("--host not supported in bookmarks-only mode")
		}
		if opts.ConnectBackend != "" {
			return opts, errors.New("--connect-backend not supported in bookmarks-only mode")
		}
	}

	homePath, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[WARN] can't detect home dir: %v", err)
		homePath = os.Getenv("HOME")
	}

	if homePath != "" {
		if opts.BookmarksDir == "" {
			opts.BookmarksDir = filepath.Join(homePath, ".pgweb/bookmarks")
		}

		if opts.QueriesDir == "" {
			opts.QueriesDir = filepath.Join(homePath, ".pgweb/queries")
		}
	}

	return opts, nil
}

// SetDefaultOptions parses and assigns the options
func SetDefaultOptions() error {
	opts, err := ParseOptions([]string{})
	if err != nil {
		return err
	}
	Opts = opts
	return nil
}

// getCurrentUser returns a current user name
func getCurrentUser() string {
	u, _ := user.Current()
	if u != nil {
		return u.Username
	}
	return os.Getenv("USER")
}

// getPrefixedEnvVar returns env var with prefix, or falls back to unprefixed one
func getPrefixedEnvVar(name string) string {
	val := os.Getenv(envVarPrefix + name)
	if val == "" {
		val = os.Getenv(name)
		if val != "" {
			fmt.Printf("[DEPRECATION] Usage of %s env var is deprecated, please use PGWEB_%s variable instead\n", name, name)
		}
	}
	return val
}

// AvailableEnvVars returns list of supported env vars.
//
// TODO: These should probably be embedded into flag parsing logic so we dont have
// to maintain the list manually.
func AvailableEnvVars() string {
	return strings.Join([]string{
		"  " + envVarPrefix + "DATABASE_URL  Database connection string",
		"  " + envVarPrefix + "URL_PREFIX    HTTP server path prefix",
		"  " + envVarPrefix + "SESSIONS      Enable multiple database sessions",
		"  " + envVarPrefix + "LOCK_SESSION  Lock session to a single database connection",
		"  " + envVarPrefix + "AUTH_USER     HTTP basic auth username",
		"  " + envVarPrefix + "AUTH_PASS     HTTP basic auth password",
	}, "\n")
}
