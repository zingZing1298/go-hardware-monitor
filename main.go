// package main

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"sync"
// 	"time"

// 	"github.com/zingZing1298/go-hardware-monitor/pkgs/hardware"
// 	"nhooyr.io/websocket"
// )

// type server struct {
// 	subscriberMessageBuffer int
// 	mux                     http.ServeMux
// 	subscribersMutex        sync.Mutex
// 	subscribers             map[*subscriber]struct{}
// }

// type subscriber struct {
// 	msgs chan []byte
// }

// func MakeNewServer() *server {
// 	s := &server{
// 		subscriberMessageBuffer: 10,
// 		subscribers:             make(map[*subscriber]struct{}),
// 	}
// 	s.mux.Handle("/", http.FileServer(http.Dir("./templates")))
// 	s.mux.HandleFunc("/ws", s.subscribeHandler)
// 	return s
// }

// func (s *server) addSubscriber(subscriber *subscriber) {
// 	s.subscribersMutex.Lock()
// 	s.subscribers[subscriber] = struct{}{}
// 	s.subscribersMutex.Unlock()

// 	fmt.Println("added Subscriber")
// }

// func (s *server) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 	var c *websocket.Conn
// 	subscriber := &subscriber{
// 		msgs: make(chan []byte, s.subscriberMessageBuffer),
// 	}
// 	s.addSubscriber(subscriber)

// 	c, err := websocket.Accept(w, r, nil)

// 	if err != nil {
// 		return err
// 	}
// 	defer c.CloseNow()
// 	ctx = c.CloseRead(ctx)
// 	//Keep looping on thread
// 	for {
// 		select {
// 		case msg := <-subscriber.msgs:
// 			ctx, cancel := context.WithTimeout(ctx, time.Second)
// 			defer cancel()
// 			err := c.Write(ctx, websocket.MessageText, msg)
// 			if err != nil {
// 				return err
// 			}
// 		case <-ctx.Done():
// 			return ctx.Err()

// 		}
// 	}
// }

// func (s *server) subscribeHandler(w http.ResponseWriter, r *http.Request) {
// 	err := s.subscribe(r.Context(), w, r)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

// func (s *server) broadcast(msg []byte) {
// 	s.subscribersMutex.Lock()
// 	for subscriber := range s.subscribers {
// 		subscriber.msgs <- msg
// 	}
// 	s.subscribersMutex.Unlock()
// }

// func main() {
// 	srv := MakeNewServer()
// 	go func(s *server) {
// 		for {
// 			systemDetails, err := hardware.GetSystemDetails()
// 			if err != nil {
// 				fmt.Println(err)
// 			}

// 			diskDetails, err := hardware.GetDiskDetails()
// 			if err != nil {
// 				fmt.Println(err)
// 			}

// 			cpuDetails, err := hardware.GetCPUDetails()
// 			if err != nil {
// 				fmt.Println(err)
// 			}

// 			netDetails, err := hardware.GetNetworkUsage()
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			// fmt.Println(systemDetails)
// 			// fmt.Println(diskDetails)
// 			// fmt.Println(cpuDetails)
// 			// fmt.Println(netDetails)
// 			timeStamp := time.Now().Format("2006-01-02 15:04:05")
// 			html := ` <div hx-swap-oob="innerHTML:#update-timestamp"> ` + timeStamp + `</div>`
// 			s.broadcast([]byte(html))
// 			s.broadcast([]byte(systemDetails))
// 			s.broadcast([]byte(diskDetails))
// 			s.broadcast([]byte(cpuDetails))
// 			s.broadcast([]byte(netDetails))

// 			time.Sleep(3 * time.Second)
// 		}

// 	}(srv)
// 	// time.Sleep(5 * time.Minute)\
// 	// Run server

// 	err := http.ListenAndServe(":8000", &srv.mux)

// 	if err != nil {
// 		fmt.Println("Server crash...\n", err)
// 		os.Exit(1)
// 	}

// }
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/zingZing1298/go-hardware-monitor/pkgs/hardware"
	"nhooyr.io/websocket"
)

type server struct {
	subscriberMessageBuffer int
	mux                     http.ServeMux
	subscribersMutex        sync.Mutex
	subscribers             map[*subscriber]struct{}
}

type subscriber struct {
	msgs chan []byte
}

func MakeNewServer() *server {
	s := &server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*subscriber]struct{}),
	}
	s.mux.Handle("/", http.FileServer(http.Dir("./templates")))
	s.mux.HandleFunc("/ws", s.subscribeHandler)
	return s
}

func (s *server) addSubscriber(subscriber *subscriber) {
	s.subscribersMutex.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscribersMutex.Unlock()

	fmt.Println("added Subscriber")
}

func (s *server) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	opts := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	}
	c, err := websocket.Accept(w, r, opts)
	if err != nil {
		return fmt.Errorf("failed to accept WebSocket connection: %w", err)
	}
	defer c.Close(websocket.StatusInternalError, "internal server error")

	subscriber := &subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}
	s.addSubscriber(subscriber)

	ctx = c.CloseRead(ctx)

	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			err := c.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return fmt.Errorf("failed to write message: %w", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *server) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	go func() {
		for {
			systemDetails, err := hardware.GetSystemDetails()
			if err != nil {
				fmt.Println(err)
			}

			diskDetails, err := hardware.GetDiskDetails()
			if err != nil {
				fmt.Println(err)
			}

			cpuDetails, err := hardware.GetCPUDetails()
			if err != nil {
				fmt.Println(err)
			}

			netDetails, err := hardware.GetNetworkUsage()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(systemDetails)
			fmt.Println(diskDetails)
			fmt.Println(cpuDetails)
			fmt.Println(netDetails)

			time.Sleep(3 * time.Second)
		}
	}()

	srv := MakeNewServer()
	err := http.ListenAndServe("127.0.0.1:8080", &srv.mux)

	if err != nil {
		fmt.Println("Server crash...\n", err)
		os.Exit(1)
	}
}
