package hub

import (
	"log"
	"strings"
	"sync"

	"github.com/desertbit/glue"
	"github.com/garyburd/redigo/redis"
)

var (
	sockets = map[*glue.Socket]map[string]bool{}
	topics  = map[string]map[*glue.Socket]bool{}

	pubconn redis.Conn
	subconn redis.PubSubConn

	l sync.RWMutex
)

func InitHub(url string) error {
	c, err := redis.DialURL(url)
	if err != nil {
		return err
	}
	pubconn = c

	c, err = redis.DialURL(url)
	if err != nil {
		return err
	}
	subconn = redis.PubSubConn{c}

	go func() {
		for {
			switch v := subconn.Receive().(type) {
			case redis.Message:
				EmitLocal(v.Channel, string(v.Data))

			case error:
				panic(v)
			}
		}
	}()

	return nil
}

func Subscribe(s *glue.Socket, t string) error {
	l.Lock()
	defer l.Unlock()

	_, ok := sockets[s]
	if !ok {
		sockets[s] = map[string]bool{}
	}
	sockets[s][t] = true

	_, ok = topics[t]
	if !ok {
		topics[t] = map[*glue.Socket]bool{}
		err := subconn.Subscribe(t)
		if err != nil {
			return err
		}
	}
	topics[t][s] = true

	return nil
}

func UnsubscribeAll(s *glue.Socket) error {
	l.Lock()
	defer l.Unlock()

	for t := range sockets[s] {
		delete(topics[t], s)
		if len(topics[t]) == 0 {
			delete(topics, t)
			err := subconn.Unsubscribe(t)
			if err != nil {
				return err
			}
		}
	}
	delete(sockets, s)

	return nil
}

func Emit(t string, m string) error {
	_, err := pubconn.Do("PUBLISH", t, m)
	return err
}

func EmitLocal(t string, m string) {
	l.RLock()
	defer l.RUnlock()

	for s := range topics[t] {
		s.Write(m)
	}
}

func HandleSocket(s *glue.Socket) {
	s.OnClose(func() {
		err := UnsubscribeAll(s)
		if err != nil {
			log.Print(err)
		}
	})

	s.OnRead(func(data string) {
		fields := strings.Fields(data)
		if len(fields) == 0 {
			return
		}
		switch fields[0] {
		case "watch":
			if len(fields) != 2 {
				return
			}
			err := Subscribe(s, fields[1])
			if err != nil {
				log.Print(err)
			}

		case "touch":
			if len(fields) != 4 {
				return
			}
			err := Emit(fields[1], "touch:"+fields[2]+","+fields[3])
			if err != nil {
				log.Print(err)
			}
		}
	})
}
