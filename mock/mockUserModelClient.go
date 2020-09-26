package mock

//UserModelClient is the mock client
type UserModelClient struct {
	GetGetUserIDFunc  func(telegramID int) (int64, error)
	GetCreateUserFunc func(firstName, username string, telegramID int) int64
}

//GetUserID is the mock client's `GetUserID` func
func (m *UserModelClient) GetUserID(telegramID int) (int64, error) {
	return GetGetUserIDFunc(telegramID)
}

//CreateUser is the mock client's `CreateUser` func
func (m *UserModelClient) CreateUser(firstName, username string, telegramID int) int64 {
	return GetCreateUserFunc(firstName, username, telegramID)
}

var (
	GetGetUserIDFunc  func(telegramID int) (int64, error)
	GetCreateUserFunc func(firstName, username string, telegramID int) int64
)
