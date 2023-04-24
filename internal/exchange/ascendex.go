package exchange

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	wsAscendexSchema = "wss"
	wsAscendexHost   = "ascendex.com"
	wsAscendexPath   = "/1/api/pro/v1/stream"
)

type BestOrderBook struct {
	Ask Order `json:"ask"` //asks.Price > any bids.Price
	Bid Order `json:"bid"`
}
type Order struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}

//go:generate moq -out ascendexConn_moq_test.go . ascendexConn
type ascendexConn interface {
	Close() error
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
}

type Ascendex struct {
	l         *logrus.Logger
	Conn      ascendexConn
	channel   string
	writeChan chan []byte
}

func NewAscendex(l *logrus.Logger) *Ascendex {
	return &Ascendex{
		l:         l,
		channel:   "bbo",
		writeChan: make(chan []byte),
	}
}

func (a *Ascendex) Connection() error {
	u := url.URL{Scheme: wsAscendexSchema, Host: wsAscendexHost, Path: wsAscendexPath}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("error dial : %s", err.Error())
	}
	a.Conn = c
	return nil
}

func (a *Ascendex) Disconnect() {
	if err := a.Conn.Close(); err != nil {
		a.l.Error("error disconection :", err)
	}
}
func (a *Ascendex) SubscribeToChannel(symbol string) error {
	type ascendexSubscribe struct {
		Operation string `json:"op"`
		Id        string `json:"id"`
		Channel   string `json:"ch"`
	}
	sub := ascendexSubscribe{
		Operation: "sub",
		Id:        strconv.Itoa(1),
		Channel:   a.channel + ":" + symbol,
	}
	subByte, err := json.Marshal(sub)
	if err != nil {
		return fmt.Errorf("marshaling error: %v", err)
	}
	if err := a.Conn.WriteMessage(websocket.TextMessage, subByte); err != nil {
		return fmt.Errorf("writing message error : %v", err)
	}
	return nil
}

type ascendexParams struct {
	Bid []string `json:"bid"`
	Ask []string `json:"ask"`
}

type ascendexMsg struct {
	Method string         `json:"m"`
	Symbol string         `json:"symbol"`
	Data   ascendexParams `json:"data"`
}

func (a *Ascendex) ReadMessagesFromChannel(ch chan<- BestOrderBook) {
	for {
		_, resp, err := a.Conn.ReadMessage()
		if err != nil {
			a.l.Errorf("reading message error: %v", err)
			continue
		}

		if len(resp) == 0 {
			continue
		}

		msg := ascendexMsg{}
		if err := json.Unmarshal(resp, &msg); err != nil {
			a.l.Errorf("unmarshaling error: %v, resp: %s", err, string(resp))
			continue
		}

		if msg.Method == "ping" {
			msg := map[string]string{"op": "pong"}
			msgByte, err := json.Marshal(msg)
			if err != nil {
				a.l.Errorf("marshaling error: %v", err)
				continue
			}

			a.writeChan <- msgByte
			continue
		}
		if len(msg.Data.Ask) == 0 || len(msg.Data.Bid) == 0 {
			continue
		}
		priceBid, err := strconv.ParseFloat(msg.Data.Bid[0], 64)
		if err != nil {
			a.l.Warnf("Float parsing error: %v", err)
			continue
		}
		priceAsk, err := strconv.ParseFloat(msg.Data.Ask[0], 64)
		if err != nil {
			a.l.Warnf("Float parsing error: %v", err)
			continue
		}

		amountBid, err := strconv.ParseFloat(msg.Data.Bid[1], 64)
		if err != nil {
			a.l.Warnf("Float parsing error: %v", err)
			continue
		}
		amountAsk, err := strconv.ParseFloat(msg.Data.Ask[1], 64)
		if err != nil {
			a.l.Warnf("Float parsing error: %v", err)
			continue
		}

		ch <- BestOrderBook{
			Ask: Order{
				Amount: amountAsk,
				Price:  priceAsk,
			},
			Bid: Order{
				Amount: amountBid,
				Price:  priceBid,
			},
		}

	}
}

func (a *Ascendex) WriteMessagesToChannel() {
	for {
		msg := <-a.writeChan
		if len(msg) == 0 {
			continue
		}

		if err := a.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			a.l.Errorf("writing message error : %v", err)
		}
	}
}
