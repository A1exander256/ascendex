package exchange

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewAscendex(t *testing.T) {
	t.Run("must_create_conn", func(t *testing.T) {
		l := logrus.New()
		a := NewAscendex(l)
		require.NotNil(t, a)
		require.IsType(t, &Ascendex{}, a)
	})
}

func TestConnection(t *testing.T) {
	t.Run("must_establish_a_connection", func(t *testing.T) {
		a := setup(t)
		err := a.Connection()
		require.Nil(t, err)
		require.NotNil(t, a.Conn)
	})
}

func TestDisconnect(t *testing.T) {
	a := setup(t)
	t.Run("must_return_nil", func(t *testing.T) {
		connMoq := ascendexConnMock{
			CloseFunc: func() error {
				return nil
			},
		}
		a.Conn = &connMoq
		output := &bytes.Buffer{}
		a.l.Out = output
		a.Disconnect()

		require.Len(t, connMoq.CloseCalls(), 1)
		lines := strings.Split(strings.TrimSpace(output.String()), "\n")
		require.Empty(t, lines[0])
	})

	t.Run("must_return_error", func(t *testing.T) {
		connMoq := &ascendexConnMock{
			CloseFunc: func() error {
				return errors.New("test_error")
			},
		}
		a.Conn = connMoq
		output := &bytes.Buffer{}
		a.l.Out = output
		a.Disconnect()

		require.Len(t, connMoq.CloseCalls(), 1)
		lines := strings.Split(strings.TrimSpace(output.String()), "\n")
		require.NotEmpty(t, lines[0])
		require.Contains(t, lines[0], "error disconection :")
		require.Contains(t, lines[0], "test_error")
	})
}

func TestSubscribeToChannel(t *testing.T) {
	a := setup(t)
	t.Run("must_subscribe_to_channel", func(t *testing.T) {
		connMoq := &ascendexConnMock{
			WriteMessageFunc: func(messageType int, data []byte) error {
				return nil
			},
		}
		a.Conn = connMoq

		err := a.SubscribeToChannel("test")

		require.Nil(t, err)
		require.Len(t, connMoq.WriteMessageCalls(), 1)
	})

	t.Run("must_return_error", func(t *testing.T) {
		connMoq := &ascendexConnMock{
			WriteMessageFunc: func(messageType int, data []byte) error {
				return errors.New("test_error")
			},
		}
		a.Conn = connMoq
		err := a.SubscribeToChannel("test")

		require.NotNil(t, err)
		require.Len(t, connMoq.WriteMessageCalls(), 1)
		require.Contains(t, err.Error(), "writing message error :")
		require.Contains(t, err.Error(), "test_error")
	})
}

func TestReadMessagesFromChannel(t *testing.T) {
	a := setup(t)
	ch := make(chan BestOrderBook)

	msgTest := ascendexMsg{
		Data: ascendexParams{
			Bid: []string{"0.1", "0.1"},
			Ask: []string{"0.2", "0.2"},
		},
	}
	connMoq := &ascendexConnMock{
		ReadMessageFunc: func() (int, []byte, error) {
			jsonByte, err := json.Marshal(msgTest)
			require.Nil(t, err)
			return 0, jsonByte, nil
		},
	}

	a.Conn = connMoq

	go a.ReadMessagesFromChannel(ch)

	msg := <-ch

	require.NotNil(t, msg)
	require.GreaterOrEqual(t, len(connMoq.ReadMessageCalls()), 1)

	require.Equal(t, msgTest.Data.Ask[0], fmt.Sprint(msg.Ask.Price))
	require.Equal(t, msgTest.Data.Ask[1], fmt.Sprint(msg.Ask.Amount))
	require.Equal(t, msgTest.Data.Bid[0], fmt.Sprint(msg.Bid.Price))
	require.Equal(t, msgTest.Data.Bid[1], fmt.Sprint(msg.Bid.Amount))
}

func TestWriteMessagesToChannel(t *testing.T) {
	a := setup(t)
	connMoq := &ascendexConnMock{
		WriteMessageFunc: func(messageType int, data []byte) error {
			return nil
		},
	}

	a.Conn = connMoq
	go a.WriteMessagesToChannel()

	a.writeChan <- []byte("test")
	time.Sleep(time.Millisecond * 100)
	require.GreaterOrEqual(t, len(connMoq.WriteMessageCalls()), 1)
}
func setup(t *testing.T) *Ascendex {
	l := logrus.New()
	a := NewAscendex(l)
	require.NotNil(t, a)
	require.IsType(t, &Ascendex{}, a)
	return a
}
