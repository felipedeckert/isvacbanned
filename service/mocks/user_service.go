// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"isvacbanned/service"
	"sync"
)

// Ensure, that UserServiceInterfaceMock does implement service.UserServiceInterface.
// If this is not the case, regenerate this file with moq.
var _ service.UserServiceInterface = &UserServiceInterfaceMock{}

// UserServiceInterfaceMock is a mock implementation of service.UserServiceInterface.
//
// 	func TestSomethingThatUsesUserServiceInterface(t *testing.T) {
//
// 		// make and configure a mocked service.UserServiceInterface
// 		mockedUserServiceInterface := &UserServiceInterfaceMock{
// 			getUserIDFunc: func(chat *tb.Chat) int64 {
// 				panic("mock out the getUserID method")
// 			},
// 		}
//
// 		// use mockedUserServiceInterface in code that requires service.UserServiceInterface
// 		// and then make assertions.
//
// 	}
type UserServiceInterfaceMock struct {
	// getUserIDFunc mocks the getUserID method.
	getUserIDFunc func(chat *tb.Chat) int64

	// calls tracks calls to the methods.
	calls struct {
		// getUserID holds details about calls to the getUserID method.
		getUserID []struct {
			// Chat is the chat argument value.
			Chat *tb.Chat
		}
	}
	lockgetUserID sync.RWMutex
}

// getUserID calls getUserIDFunc.
func (mock *UserServiceInterfaceMock) getUserID(chat *tb.Chat) int64 {
	callInfo := struct {
		Chat *tb.Chat
	}{
		Chat: chat,
	}
	mock.lockgetUserID.Lock()
	mock.calls.getUserID = append(mock.calls.getUserID, callInfo)
	mock.lockgetUserID.Unlock()
	if mock.getUserIDFunc == nil {
		var (
			nOut int64
		)
		return nOut
	}
	return mock.getUserIDFunc(chat)
}

// getUserIDCalls gets all the calls that were made to getUserID.
// Check the length with:
//     len(mockedUserServiceInterface.getUserIDCalls())
func (mock *UserServiceInterfaceMock) getUserIDCalls() []struct {
	Chat *tb.Chat
} {
	var calls []struct {
		Chat *tb.Chat
	}
	mock.lockgetUserID.RLock()
	calls = mock.calls.getUserID
	mock.lockgetUserID.RUnlock()
	return calls
}