package api

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sosedoff/pgweb/pkg/client"
	"github.com/sosedoff/pgweb/pkg/metrics"
)

// 会话管理器，基于 sync.Mutex + map 实现
type SessionManager struct {
	logger      *logrus.Logger            // 日志
	sessions    map[string]*client.Client // 不同会话的客户端
	mu          sync.Mutex                // 锁
	idleTimeout time.Duration             // 超时时间
}

func NewSessionManager(logger *logrus.Logger) *SessionManager {
	return &SessionManager{
		logger:   logger,
		sessions: map[string]*client.Client{},
		mu:       sync.Mutex{},
	}
}

func (m *SessionManager) SetIdleTimeout(timeout time.Duration) {
	m.idleTimeout = timeout
}

// 获取所有的会话 id
func (m *SessionManager) IDs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	ids := []string{}
	for k := range m.sessions {
		ids = append(ids, k)
	}

	return ids
}

// 获取所有的会话的副本
func (m *SessionManager) Sessions() map[string]*client.Client {
	m.mu.Lock()
	defer m.mu.Unlock()

	sessions := make(map[string]*client.Client, len(m.sessions))
	for k, v := range m.sessions {
		sessions[k] = v
	}

	return sessions
}

// 获取指定的会话
func (m *SessionManager) Get(id string) *client.Client {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.sessions[id]
}

// 添加新的会话
func (m *SessionManager) Add(id string, conn *client.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessions[id] = conn
	metrics.SetSessionsCount(len(m.sessions))
}

// 移除会话
func (m *SessionManager) Remove(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, ok := m.sessions[id]
	if ok {
		conn.Close()
		delete(m.sessions, id)
	}

	metrics.SetSessionsCount(len(m.sessions))
	return ok
}

// 会话总数
func (m *SessionManager) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.sessions)
}

// 清理
func (m *SessionManager) Cleanup() int {
	if m.idleTimeout == 0 {
		return 0
	}

	removed := 0

	m.logger.Debug("starting idle sessions cleanup")
	defer func() {
		m.logger.Debug("removed idle sessions:", removed)
	}()

	// 获取可清理的会话
	for _, id := range m.staleSessions() {
		m.logger.WithField("id", id).Debug("closing stale session")
		// 清理会话
		if m.Remove(id) {
			removed++
		}
	}

	return removed
}

// 每分钟执行一次清理
func (m *SessionManager) RunPeriodicCleanup() {
	m.logger.WithField("timeout", m.idleTimeout).Info("session manager cleanup enabled")

	for range time.Tick(time.Minute) {
		m.Cleanup()
	}
}

// 返回超时会话的id列表
func (m *SessionManager) staleSessions() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	ids := []string{}

	for id, conn := range m.sessions {
		if now.Sub(conn.LastQueryTime()) > m.idleTimeout {
			ids = append(ids, id)
		}
	}

	return ids
}
