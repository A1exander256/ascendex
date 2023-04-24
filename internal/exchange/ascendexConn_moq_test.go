// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package exchange

import (
	"sync"
)

// Ensure, that ascendexConnMock does implement ascendexConn.
// If this is not the case, regenerate this file with moq.
var _ ascendexConn = &ascendexConnMock{}

// ascendexConnMock is a mock implementation of ascendexConn.
//
//	func TestSomethingThatUsesascendexConn(t *testing.T) {
//
//		// make and configure a mocked ascendexConn
//		mockedascendexConn := &ascendexConnMock{
//			CloseFunc: func() error {
//				panic("mock out the Close method")
//			},
//			ReadMessageFunc: func() (int, []byte, error) {
//				panic("mock out the ReadMessage method")
//			},
//			WriteMessageFunc: func(messageType int, data []byte) error {
//				panic("mock out the WriteMessage method")
//			},
//		}
//
//		// use mockedascendexConn in code that requires ascendexConn
//		// and then make assertions.
//
//	}
type ascendexConnMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func() error

	// ReadMessageFunc mocks the ReadMessage method.
	ReadMessageFunc func() (int, []byte, error)

	// WriteMessageFunc mocks the WriteMessage method.
	WriteMessageFunc func(messageType int, data []byte) error

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// ReadMessage holds details about calls to the ReadMessage method.
		ReadMessage []struct {
		}
		// WriteMessage holds details about calls to the WriteMessage method.
		WriteMessage []struct {
			// MessageType is the messageType argument value.
			MessageType int
			// Data is the data argument value.
			Data []byte
		}
	}
	lockClose        sync.RWMutex
	lockReadMessage  sync.RWMutex
	lockWriteMessage sync.RWMutex
}

// Close calls CloseFunc.
func (mock *ascendexConnMock) Close() error {
	if mock.CloseFunc == nil {
		panic("ascendexConnMock.CloseFunc: method is nil but ascendexConn.Close was just called")
	}
	callInfo := struct {
	}{}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	return mock.CloseFunc()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//
//	len(mockedascendexConn.CloseCalls())
func (mock *ascendexConnMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// ReadMessage calls ReadMessageFunc.
func (mock *ascendexConnMock) ReadMessage() (int, []byte, error) {
	if mock.ReadMessageFunc == nil {
		panic("ascendexConnMock.ReadMessageFunc: method is nil but ascendexConn.ReadMessage was just called")
	}
	callInfo := struct {
	}{}
	mock.lockReadMessage.Lock()
	mock.calls.ReadMessage = append(mock.calls.ReadMessage, callInfo)
	mock.lockReadMessage.Unlock()
	return mock.ReadMessageFunc()
}

// ReadMessageCalls gets all the calls that were made to ReadMessage.
// Check the length with:
//
//	len(mockedascendexConn.ReadMessageCalls())
func (mock *ascendexConnMock) ReadMessageCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockReadMessage.RLock()
	calls = mock.calls.ReadMessage
	mock.lockReadMessage.RUnlock()
	return calls
}

// WriteMessage calls WriteMessageFunc.
func (mock *ascendexConnMock) WriteMessage(messageType int, data []byte) error {
	if mock.WriteMessageFunc == nil {
		panic("ascendexConnMock.WriteMessageFunc: method is nil but ascendexConn.WriteMessage was just called")
	}
	callInfo := struct {
		MessageType int
		Data        []byte
	}{
		MessageType: messageType,
		Data:        data,
	}
	mock.lockWriteMessage.Lock()
	mock.calls.WriteMessage = append(mock.calls.WriteMessage, callInfo)
	mock.lockWriteMessage.Unlock()
	return mock.WriteMessageFunc(messageType, data)
}

// WriteMessageCalls gets all the calls that were made to WriteMessage.
// Check the length with:
//
//	len(mockedascendexConn.WriteMessageCalls())
func (mock *ascendexConnMock) WriteMessageCalls() []struct {
	MessageType int
	Data        []byte
} {
	var calls []struct {
		MessageType int
		Data        []byte
	}
	mock.lockWriteMessage.RLock()
	calls = mock.calls.WriteMessage
	mock.lockWriteMessage.RUnlock()
	return calls
}