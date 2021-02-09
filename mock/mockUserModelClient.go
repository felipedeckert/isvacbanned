package mock

//UserModelClient is the mock client
type UserModelClient struct {
	GetGetUserIDFunc  func(telegramID int64) (int64, error)
	GetCreateUserFunc func(firstName, username string, telegramID int64) int64
	GetInactivateUser func(userID int64) int64
	GetActivateUser   func(userID int64) int64
}

//GetUserID is the mock client's `GetUserID` func
func (m *UserModelClient) GetUserID(telegramID int64) (int64, error) {
	return GetGetUserIDFunc(telegramID)
}

//CreateUser is the mock client's `CreateUser` func
func (m *UserModelClient) CreateUser(firstName, username string, telegramID int64) int64 {
	return GetCreateUserFunc(firstName, username, telegramID)
}

//InactivateUser is the mock client's `InactivateUser` func
func (m *UserModelClient) InactivateUser(userID int64) int64 {
	return GetInactivateUser(userID)
}

//ActivateUser is the mock client's `ActivateUser` func
func (m *UserModelClient) ActivateUser(userID int64) int64 {
	return GetActivateUser(userID)
}

var (
	GetGetUserIDFunc  func(telegramID int64) (int64, error)
	GetCreateUserFunc func(firstName, username string, telegramID int64) int64
	GetInactivateUser func(userID int64) int64
	GetActivateUser   func(userID int64) int64
)
