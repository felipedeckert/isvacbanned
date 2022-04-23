// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"isvacbanned/domain/telegrambot"
	"sync"
)

// Ensure, that UnfollowPlayerFollowPersistenceGatewayMock does implement telegrambot.UnfollowPlayerFollowPersistenceGateway.
// If this is not the case, regenerate this file with moq.
var _ telegrambot.UnfollowPlayerFollowPersistenceGateway = &UnfollowPlayerFollowPersistenceGatewayMock{}

// UnfollowPlayerFollowPersistenceGatewayMock is a mock implementation of telegrambot.UnfollowPlayerFollowPersistenceGateway.
//
// 	func TestSomethingThatUsesUnfollowPlayerFollowPersistenceGateway(t *testing.T) {
//
// 		// make and configure a mocked telegrambot.UnfollowPlayerFollowPersistenceGateway
// 		mockedUnfollowPlayerFollowPersistenceGateway := &UnfollowPlayerFollowPersistenceGatewayMock{
// 			UnfollowSteamUserFunc: func(ctx context.Context, userID int64, steamID string) (int64, error) {
// 				panic("mock out the UnfollowSteamUser method")
// 			},
// 		}
//
// 		// use mockedUnfollowPlayerFollowPersistenceGateway in code that requires telegrambot.UnfollowPlayerFollowPersistenceGateway
// 		// and then make assertions.
//
// 	}
type UnfollowPlayerFollowPersistenceGatewayMock struct {
	// UnfollowSteamUserFunc mocks the UnfollowSteamUser method.
	UnfollowSteamUserFunc func(ctx context.Context, userID int64, steamID string) (int64, error)

	// calls tracks calls to the methods.
	calls struct {
		// UnfollowSteamUser holds details about calls to the UnfollowSteamUser method.
		UnfollowSteamUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID int64
			// SteamID is the steamID argument value.
			SteamID string
		}
	}
	lockUnfollowSteamUser sync.RWMutex
}

// UnfollowSteamUser calls UnfollowSteamUserFunc.
func (mock *UnfollowPlayerFollowPersistenceGatewayMock) UnfollowSteamUser(ctx context.Context, userID int64, steamID string) (int64, error) {
	callInfo := struct {
		Ctx     context.Context
		UserID  int64
		SteamID string
	}{
		Ctx:     ctx,
		UserID:  userID,
		SteamID: steamID,
	}
	mock.lockUnfollowSteamUser.Lock()
	mock.calls.UnfollowSteamUser = append(mock.calls.UnfollowSteamUser, callInfo)
	mock.lockUnfollowSteamUser.Unlock()
	if mock.UnfollowSteamUserFunc == nil {
		var (
			nOut   int64
			errOut error
		)
		return nOut, errOut
	}
	return mock.UnfollowSteamUserFunc(ctx, userID, steamID)
}

// UnfollowSteamUserCalls gets all the calls that were made to UnfollowSteamUser.
// Check the length with:
//     len(mockedUnfollowPlayerFollowPersistenceGateway.UnfollowSteamUserCalls())
func (mock *UnfollowPlayerFollowPersistenceGatewayMock) UnfollowSteamUserCalls() []struct {
	Ctx     context.Context
	UserID  int64
	SteamID string
} {
	var calls []struct {
		Ctx     context.Context
		UserID  int64
		SteamID string
	}
	mock.lockUnfollowSteamUser.RLock()
	calls = mock.calls.UnfollowSteamUser
	mock.lockUnfollowSteamUser.RUnlock()
	return calls
}