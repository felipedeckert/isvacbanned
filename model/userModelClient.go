package model

type UserModelClient interface {
	GetUserID(telegramID int) (int64, error)
	CreateUser(firstName, username string, telegramID int) int64
	ActivateUser(userID int64) int64
	InactivateUser(userID int64) int64
}
