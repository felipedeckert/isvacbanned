package model

type UserMock struct{}

func (u UserMock) GetUserID(telegramID int64) (int64, error) {
	return int64(30), nil
}

// CreateUser inserts a new user in the database
func (u UserMock) CreateUser(firstName, username string, telegramID int64) int64 {
	return int64(31)
}

//InactivateUser sets user flag is_active to false
func (u UserMock) InactivateUser(userID int64) int64 {
	return int64(1)
}

//ActivateUser sets user flag is_active to true
func (u UserMock) ActivateUser(userID int64) int64 {
	return int64(1)
}
