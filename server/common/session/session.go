package session

import (
	"strconv"
	"strings"
)

// Session 用户会话信息 记录用户当前在线状态、连接信息、身份认证等数据
type Session string

// NewSession 创建一个会话
// @param name 用户名
// @param token 用户令牌
// @param id 用户ID
// @return Session 会话
func NewSession(name, token string, id uint64) Session {
	if len(name) == 0 || len(token) == 0 {
		panic("name or token is empty")
	}
	idstr := strconv.FormatUint(id, 10)
	return Session(name + ":" + token + ":" + idstr)
}

// FromString 从字符串中解析出会话
// @param str 字符串
// @return Session 会话
func FromString(str string) Session {
	return Session(str)
}

// Name 获取会话的用户名
// @return string 用户名
func (s Session) Name() string {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}
	return arr[0]
}

// Token 获取会话的用户令牌
// @return string 用户令牌
func (s Session) Token() string {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}
	return arr[1]
}

// Id 获取会话的用户ID
// @return uint64 用户ID
func (s Session) Id() uint64 {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}
	id, err := strconv.ParseUint(arr[2], 10, 64)
	if err != nil {
		panic("invalid id")
	}
	return id
}

// Info 获取会话的用户信息
// @return string 用户名
// @return string 用户令牌
// @return uint64 用户ID
func (s Session) Info() (string, string, uint64) {
	arr := strings.Split(string(s), ":")
	if len(arr) != 3 {
		panic("invalid session")
	}
	id, err := strconv.ParseUint(arr[2], 10, 64)
	if err != nil {
		panic("invalid id")
	}
	return arr[0], arr[1], id
}

// String 将会话转换为字符串
// @return string 字符串
func (s Session) String() string {
	return string(s)
}
