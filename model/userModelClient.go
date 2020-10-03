package model

type UserModelClient interface {
	GetUserID(telegramID int64) (int64, error)
	CreateUser(firstName, username string, telegramID int64) int64
	ActivateUser(userID int64) int64
	InactivateUser(userID int64) int64
}
