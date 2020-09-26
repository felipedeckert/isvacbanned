package service

type UserModelClient interface {
	GetUserID(telegramID int) (int64, error)
	CreateUser(firstName, username string, telegramID int) int64
}
