// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"isvacbanned/domain/entities"
	"isvacbanned/domain/telegrambot"
	"sync"
)

// Ensure, that ShowPlayersFollowPersistenceGatewayMock does implement telegrambot.ShowPlayersFollowPersistenceGateway.
// If this is not the case, regenerate this file with moq.
var _ telegrambot.ShowPlayersFollowPersistenceGateway = &ShowPlayersFollowPersistenceGatewayMock{}

// ShowPlayersFollowPersistenceGatewayMock is a mock implementation of telegrambot.ShowPlayersFollowPersistenceGateway.
//
// 	func TestSomethingThatUsesShowPlayersFollowPersistenceGateway(t *testing.T) {
//
// 		// make and configure a mocked telegrambot.ShowPlayersFollowPersistenceGateway
// 		mockedShowPlayersFollowPersistenceGateway := &ShowPlayersFollowPersistenceGatewayMock{
// 			GetUsersFollowedFunc: func(ctx context.Context, userID int64) ([]entities.UsersFollowed, error) {
// 				panic("mock out the GetUsersFollowed method")
// 			},
// 		}
//
// 		// use mockedShowPlayersFollowPersistenceGateway in code that requires telegrambot.ShowPlayersFollowPersistenceGateway
// 		// and then make assertions.
//
// 	}
type ShowPlayersFollowPersistenceGatewayMock struct {
	// GetUsersFollowedFunc mocks the GetUsersFollowed method.
	GetUsersFollowedFunc func(ctx context.Context, userID int64) ([]entities.UsersFollowed, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetUsersFollowed holds details about calls to the GetUsersFollowed method.
		GetUsersFollowed []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID int64
		}
	}
	lockGetUsersFollowed sync.RWMutex
}

// GetUsersFollowed calls GetUsersFollowedFunc.
func (mock *ShowPlayersFollowPersistenceGatewayMock) GetUsersFollowed(ctx context.Context, userID int64) ([]entities.UsersFollowed, error) {
	callInfo := struct {
		Ctx    context.Context
		UserID int64
	}{
		Ctx:    ctx,
		UserID: userID,
	}
	mock.lockGetUsersFollowed.Lock()
	mock.calls.GetUsersFollowed = append(mock.calls.GetUsersFollowed, callInfo)
	mock.lockGetUsersFollowed.Unlock()
	if mock.GetUsersFollowedFunc == nil {
		var (
			usersFollowedsOut []entities.UsersFollowed
			errOut            error
		)
		return usersFollowedsOut, errOut
	}
	return mock.GetUsersFollowedFunc(ctx, userID)
}

// GetUsersFollowedCalls gets all the calls that were made to GetUsersFollowed.
// Check the length with:
//     len(mockedShowPlayersFollowPersistenceGateway.GetUsersFollowedCalls())
func (mock *ShowPlayersFollowPersistenceGatewayMock) GetUsersFollowedCalls() []struct {
	Ctx    context.Context
	UserID int64
} {
	var calls []struct {
		Ctx    context.Context
		UserID int64
	}
	mock.lockGetUsersFollowed.RLock()
	calls = mock.calls.GetUsersFollowed
	mock.lockGetUsersFollowed.RUnlock()
	return calls
}