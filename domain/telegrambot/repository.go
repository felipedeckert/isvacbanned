package telegrambot

import (
	"context"
	"isvacbanned/domain/entities"
)

type Repository interface {
	FollowSteamUser(ctx context.Context, chatID int64, steamID, currNickname string, userID int64) (int64, error)
	UnfollowSteamUser(ctx context.Context, userID int64, steamID string) (int64, error)
	GetFollowerCountBySteamID(ctx context.Context, steamID string) (int64, error)
	GetAllIncompleteFollowedUsers(ctx context.Context) (map[int64][]entities.UsersFollowed, error)
	GetUsersFollowed(ctx context.Context, userID int64) ([]entities.UsersFollowed, error)
	SetCurrNickname(ctx context.Context, userID int64, sanitizedActualNickname string) error
	SetFollowedUserToCompleted(ctx context.Context, id []int64)
	GetUsersFollowedSummary(ctx context.Context, userID int64) (map[bool]int, error)
	IsFollowed(ctx context.Context, steamID string, userID int64) (string, int64, error)
}
