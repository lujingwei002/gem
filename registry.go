package gem

type Registry interface {
	// 注册账号
	AddUser(userID UserID, logout ForceLogout) error
}
