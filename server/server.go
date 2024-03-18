package server

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type gameServer struct {
	subscriberMessageBuffer int
	publishLimiter          *rate.Limiter
	logf                    func(f string, v ...interface{})
	serveMux                http.ServeMux
	subscribersMu           sync.Mutex
	subscribers             map[*subscriber]struct{}
}

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

func (gs *gameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gs.serveMux.ServeHTTP(w, r)
}

func NewGameServer() *gameServer {
	gs := &gameServer{
		subscriberMessageBuffer: 16,
		logf:                    log.Printf,
		subscribers:             make(map[*subscriber]struct{}),
		publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
	gs.serveMux.HandleFunc("/subscribe", gs.subscribeHandler)
	gs.serveMux.HandleFunc("/publish", gs.publishHandler)

	return gs
}

func (gs *gameServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := gs.subscribe(r.Context(), w, r)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		gs.logf("%v", err)
		return
	}
}

func (gs *gameServer) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := io.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	gs.publish(msg)

	w.WriteHeader(http.StatusAccepted)
}

func (gs *gameServer) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	subs := &subscriber{
		msgs: make(chan []byte, gs.subscriberMessageBuffer),
		closeSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
			}
		},
	}
	gs.addSubscriber(subs)
	defer gs.deleteSubscriber(subs)

	c2, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	mu.Lock()
	if closed {
		mu.Unlock()
		return net.ErrClosed
	}
	c = c2
	mu.Unlock()
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)

	for {
		select {
		case msg := <-subs.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (gs *gameServer) publish(msg []byte) {
	gs.subscribersMu.Lock()
	defer gs.subscribersMu.Unlock()

	gs.publishLimiter.Wait(context.Background())

	for s := range gs.subscribers {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}

// addSubscriber registers a subscriber.
func (gs *gameServer) addSubscriber(s *subscriber) {
	gs.subscribersMu.Lock()
	gs.subscribers[s] = struct{}{}
	gs.subscribersMu.Unlock()
}

// deleteSubscriber deletes the given subscriber.
func (gs *gameServer) deleteSubscriber(s *subscriber) {
	gs.subscribersMu.Lock()
	delete(gs.subscribers, s)
	gs.subscribersMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
