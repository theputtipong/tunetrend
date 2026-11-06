package upstashtest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type Server struct {
	*httptest.Server

	mu       sync.Mutex
	data     map[string]string
	ttls     map[string]int64
	commands [][]string
	Fail     bool
	failKey  string
}

func NewServer(t *testing.T) *Server {
	s := &Server{data: map[string]string{}, ttls: map[string]int64{}}
	s.Server = httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(s.Close)
	return s
}

func (s *Server) Token() string { return "test-token" }

func (s *Server) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Server) Value(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data[key]
	return v, ok
}

func (s *Server) TTL(key string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ttls[key]
}

func (s *Server) SetFail(fail bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Fail = fail
}

func (s *Server) FailKeyPrefix(prefix string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failKey = prefix
}

type step struct {
	Result any    `json:"result"`
	Error  string `json:"error,omitempty"`
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Fail || r.Header.Get("Authorization") != "Bearer "+s.Token() || r.URL.Path != "/pipeline" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cmds [][]string
	if err := json.NewDecoder(r.Body).Decode(&cmds); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	out := make([]step, len(cmds))
	for i, cmd := range cmds {
		s.commands = append(s.commands, cmd)
		out[i] = s.exec(cmd)
	}
	_ = json.NewEncoder(w).Encode(out)
}

func (s *Server) exec(cmd []string) step {
	if len(cmd) < 2 {
		return step{Error: "ERR wrong number of arguments"}
	}
	key := cmd[1]
	if s.failKey != "" && strings.HasPrefix(key, s.failKey) {
		return step{Error: "ERR injected failure"}
	}
	switch strings.ToUpper(cmd[0]) {
	case "GET":
		v, ok := s.data[key]
		if !ok {
			return step{Result: nil}
		}
		return step{Result: v}
	case "SET":
		nx := false
		ttl := int64(-1)
		for i := 3; i < len(cmd); i++ {
			switch strings.ToUpper(cmd[i]) {
			case "NX":
				nx = true
			case "EX":
				if i+1 < len(cmd) {
					ttl, _ = strconv.ParseInt(cmd[i+1], 10, 64)
					i++
				}
			}
		}
		if _, exists := s.data[key]; nx && exists {
			return step{Result: nil}
		}
		s.data[key] = cmd[2]
		if ttl >= 0 {
			s.ttls[key] = ttl
		} else {
			delete(s.ttls, key)
		}
		return step{Result: "OK"}
	case "INCR", "INCRBY", "DECR", "DECRBY":
		delta := int64(1)
		if len(cmd) > 2 {
			d, err := strconv.ParseInt(cmd[2], 10, 64)
			if err != nil {
				return step{Error: "ERR value is not an integer"}
			}
			delta = d
		}
		if strings.HasPrefix(strings.ToUpper(cmd[0]), "DECR") {
			delta = -delta
		}
		cur, _ := strconv.ParseInt(s.data[key], 10, 64)
		cur += delta
		s.data[key] = strconv.FormatInt(cur, 10)
		return step{Result: cur}
	case "EXPIRE":
		if _, ok := s.data[key]; !ok {
			return step{Result: 0}
		}
		nx := len(cmd) > 3 && strings.EqualFold(cmd[3], "NX")
		if _, has := s.ttls[key]; nx && has {
			return step{Result: 0}
		}
		ttl, _ := strconv.ParseInt(cmd[2], 10, 64)
		s.ttls[key] = ttl
		return step{Result: 1}
	case "DEL":
		_, ok := s.data[key]
		delete(s.data, key)
		delete(s.ttls, key)
		if ok {
			return step{Result: 1}
		}
		return step{Result: 0}
	default:
		return step{Error: "ERR unknown command " + cmd[0]}
	}
}
