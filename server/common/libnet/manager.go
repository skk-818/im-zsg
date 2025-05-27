package libnet

import (
	"server/common/hash"
	"server/common/session"
	"sync"
)

// sessionMapNum 是 sessionMap 的数量 拿来 做 session 分片，将所有的session分别存储在32和片中，每个sessionMap 存储 1/32 的session
// 这样可以提高 查找session的效率
const sessionMapNum = 32

// Manager 连接管理器 管理所有连接 用户和连接的映射关系
type Manager struct {
	Name        string                    // edge 的名字 表示这个连接管理器是用来管理哪个edge的连接的
	sessionMaps [sessionMapNum]sessionMap //
	disposeFlag bool                      // 是否已经dispose
	disposeOnce sync.Once                 // disposeOnce 用来保证只执行一次dispose操作
	disposeWait sync.WaitGroup            // disposeWait 用来等待所有连接都关闭
}

// sessionMap 是 session 的map 用来存储所有的session
// tokenSessions 是 token 的map 用来存储所有的session
type sessionMap struct {
	sync.RWMutex
	sessions      map[session.Session]*Session // 用户的 session 标识 和连接对象上的 session 映射
	tokenSessions map[string][]session.Session // token 的 session 标识 和连接对象上的 session 映射
}

// NewManager 创建一个连接管理器
func NewManager(name string) *Manager {
	manager := &Manager{
		Name: name,
	}
	for i := 0; i < sessionMapNum; i++ {
		manager.sessionMaps[i].sessions = make(map[session.Session]*Session)
		manager.sessionMaps[i].tokenSessions = make(map[string][]session.Session)
	}

	return manager
}

// GetSession 根据 sessionId 获取 session
func (m *Manager) GetSession(sessionId session.Session) *Session {
	token := sessionId.Token()
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId%sessionMapNum]

	smap.RLock()
	defer smap.RUnlock()

	return smap.sessions[sessionId]
}

// GetTokenSessions 根据 token 获取 session
func (m *Manager) GetTokenSessions(token string) []*Session {
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId%sessionMapNum]

	smap.RLock()
	defer smap.RUnlock()

	sessionIds := smap.tokenSessions[token]

	var sessions []*Session
	for _, sessionId := range sessionIds {
		sessions = append(sessions, smap.sessions[sessionId])
	}

	return sessions
}

// AddSession 添加 session
func (m *Manager) AddSession(session *Session) {
	sessionId := session.Session()
	token := session.token
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()

	smap.sessions[sessionId] = session
	smap.tokenSessions[token] = append(smap.tokenSessions[token], sessionId)
}

// RemoveSession 删除 session
func (m *Manager) RemoveSession(session *Session) {
	sessionId := session.Session()
	token := session.token
	hashId := hash.Hash([]byte(token))
	smap := &m.sessionMaps[hashId%sessionMapNum]

	smap.Lock()
	defer smap.Unlock()

	delete(smap.sessions, sessionId)

	sessionIds := smap.tokenSessions[token]
	for i, sessionId := range sessionIds {
		if sessionId == sessionId {
			sessionIds = append(sessionIds[:i], sessionIds[i+1:]...)
			break
		}
	}
	smap.tokenSessions[token] = sessionIds
}

// Close 关闭连接管理器
func (m *Manager) Close() {
	m.disposeOnce.Do(func() {
		m.disposeFlag = true

		for i := 0; i < sessionMapNum; i++ {
			smap := &m.sessionMaps[i]
			smap.Lock()
			for _, session := range smap.sessions {
				session.Close()
			}
			smap.Unlock()
		}
		m.disposeWait.Wait()
	})
}
