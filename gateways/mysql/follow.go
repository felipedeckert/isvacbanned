package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"isvacbanned/domain/entities"
	"strconv"
)

type Follow struct {
	*sql.DB
}

// FollowSteamUser links a telegram user to a steam user which is being followed
func (f *Follow) FollowSteamUser(ctx context.Context, chatID int64, steamID, currNickname string, userID int64) (int64, error) {
	query := `INSERT INTO follow(chat_id, steam_id, user_id, old_nickname, curr_nickname) VALUES($1, $2, $3, $4, $5)
				RETURNING id`

	var id int64
	err := f.QueryRowContext(ctx, query, chatID, steamID, userID, currNickname, currNickname).Scan(&id)
	if err != nil {
		return -1, errors.New(fmt.Sprintf(`error inactivating user userID:%d err:%s`, userID, err.Error()))
	}

	return id, nil
}

// UnfollowSteamUser sets a followed player flag is_active to false
func (f *Follow) UnfollowSteamUser(ctx context.Context, userID int64, steamID string) (int64, error) {
	query := `UPDATE follow 
				SET is_active = false 
				WHERE user_id = $1 
				AND steam_id  = $2
				RETURNING id`

	var id int64
	err := f.QueryRowContext(ctx, query, userID, steamID).Scan(&id)
	if err != nil {
		return -1, errors.New(fmt.Sprintf(`error inactivating user userID:%d err: %s`, userID, err.Error()))
	}

	return id, nil
}

//GetFollowerCountBySteamID get the number of followers of a steam player
func (f *Follow) GetFollowerCountBySteamID(ctx context.Context, steamID string) (int64, error) {

	query := `SELECT COUNT(f.id) as count 
				FROM follow f 
				WHERE f.steam_id = $1`

	var count int64
	err := f.QueryRowContext(ctx, query,steamID).Scan(&count)

	return count, err
}

//GetAllIncompleteFollowedUsers get all fallowed steam user for every telegram user
func (f *Follow) GetAllIncompleteFollowedUsers(ctx context.Context) (map[int64][]entities.UsersFollowed, error) {

	query := `SELECT f.id, f.chat_id, f.steam_id, f.old_nickname, f.curr_nickname 
				FROM follow f 
				JOIN user u 
					ON f.user_id = u.id 
				WHERE f.is_completed <> true 
				AND f.is_active = true 
				AND u.is_active = true`

	rows, err := f.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[int64][]entities.UsersFollowed)
	for rows.Next() {
		var (
			id           int64
			chatID       int64
			steamID      string
			oldNickname  string
			currNickName string
		)

		if err = rows.Scan(&id, &chatID, &steamID, &oldNickname, &currNickName); err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}

			return m, err
		}

		m[chatID] = append(m[chatID], entities.UsersFollowed{ID: id, SteamID: steamID, OldNickname: oldNickname, CurrNickname: currNickName})
	}

	return m, nil
}

//GetUsersFollowed gets a slice os the nicknames (the old ones) of players followed by a user
func (f *Follow) GetUsersFollowed(ctx context.Context, userID int64) ([]entities.UsersFollowed, error) {

	query  := `SELECT old_nickname, is_completed 
				FROM follow 
				WHERE user_id = $1`

	rows, err := f.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]entities.UsersFollowed, 0)
	for rows.Next() {
		var (
			oldNickname string
			isCompleted bool
		)

		if err = rows.Scan(&oldNickname, &isCompleted); err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}

			return users, nil
		}

		users = append(users, entities.UsersFollowed{OldNickname: oldNickname, IsCompleted: isCompleted})
	}

	return users, nil
}

//SetCurrNickname updates de curr_nickname of a given player
func (f *Follow) SetCurrNickname(ctx context.Context, userID int64, sanitizedActualNickname string) error {

	query := `UPDATE follow 
				SET curr_nickname = $1 
				WHERE id = $2
				RETURNING id`

	var id int64
	err := f.QueryRowContext(ctx, query, sanitizedActualNickname, userID).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

// SetFollowedUserToCompleted sets a player status to completed, and it will not be followed anymore
func (f *Follow) SetFollowedUserToCompleted(ctx context.Context, id []int64) {
	query := `UPDATE follow 
				SET is_completed = true 
				WHERE id IN($1)`

	f.QueryRowContext(ctx, query, sliceToStringParam(id))
}

// GetUsersFollowedSummary returns a summary of players followed, separated by ban status
func (f *Follow) GetUsersFollowedSummary(ctx context.Context, userID int64) (map[bool]int, error) {

	query := `SELECT COUNT(id) as count, is_completed
				FROM follow
				WHERE user_id = $1
				AND is_active = 1
				GROUP BY is_completed`

	rows, err := f.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := make(map[bool]int, 0)
	for rows.Next() {
		var (
			count       int
			isCompleted bool
		)

		if err = rows.Scan(&count, &isCompleted); err != nil {
			if err != sql.ErrNoRows {
				return  nil, err
			}
			continue
		}

		summary[isCompleted] = count
	}

	return summary, nil
}

// IsFollowed checks if a user already follows a player
func (f *Follow) IsFollowed(ctx context.Context, steamID string, userID int64) (string, int64, error) {

	query := `SELECT old_nickname, id 
				FROM follow 
				WHERE steam_id = $1 
				AND user_id = $2`

	var (
		oldNickname string
		id int64
	)
	err := f.QueryRowContext(ctx, query, steamID, userID).Scan(&oldNickname, &id)

	return oldNickname, id, err
}

func sliceToStringParam(ids []int64) string {
	if len(ids) == 0 {
		return ""
	}

	// Appr. 5 chars per num plus the comma.
	estimate := len(ids) * 6
	b := make([]byte, 0, estimate)

	for _, n := range ids {
		b = strconv.AppendInt(b, n, 10)
		b = append(b, ',')
	}
	b = b[:len(b)-1]
	return string(b)
}