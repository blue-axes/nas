package types

type UserInfo struct {
	Username string
	CanRead  bool
	CanWrite bool
	IsAdmin  bool
}