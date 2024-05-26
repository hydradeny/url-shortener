package session

type Session struct {
	SessionID string
	UserID    uint32
}

type CreateSession struct {
	UserID uint32
}

type DestroySession struct {
	SessionID string
}

type DestroyAllSession struct {
	UserID uint32
}
type CheckSession struct {
	SessionID string
}
