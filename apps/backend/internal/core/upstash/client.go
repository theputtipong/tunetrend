package upstash

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var ErrNotConfigured = errors.New("upstash: UPSTASH_REDIS_REST_URL/TOKEN not configured")

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http:    &http.Client{Timeout: 3 * time.Second},
	}
}

func NewFromEnv() *Client {
	return New(os.Getenv("UPSTASH_REDIS_REST_URL"), os.Getenv("UPSTASH_REDIS_REST_TOKEN"))
}

func (c *Client) Configured() bool {
	return c != nil && c.baseURL != "" && c.token != ""
}

type Result struct {
	Value json.RawMessage
	Err   error
}

func (r Result) IsNull() bool {
	return len(r.Value) == 0 || string(r.Value) == "null"
}

func (r Result) Int() (int64, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	if r.IsNull() {
		return 0, nil
	}
	var n int64
	if err := json.Unmarshal(r.Value, &n); err == nil {
		return n, nil
	}
	var s string
	if err := json.Unmarshal(r.Value, &s); err != nil {
		return 0, err
	}
	return strconv.ParseInt(s, 10, 64)
}

func (r Result) String() (string, bool, error) {
	if r.Err != nil {
		return "", false, r.Err
	}
	if r.IsNull() {
		return "", false, nil
	}
	var s string
	if err := json.Unmarshal(r.Value, &s); err != nil {
		return "", false, err
	}
	return s, true, nil
}

type pipelineStep struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error,omitempty"`
}

func (c *Client) Pipeline(cmds ...[]string) ([]Result, error) {
	if !c.Configured() {
		return nil, ErrNotConfigured
	}

	body, err := json.Marshal(cmds)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/pipeline", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstash: unexpected status %d", res.StatusCode)
	}

	var steps []pipelineStep
	if err := json.NewDecoder(res.Body).Decode(&steps); err != nil {
		return nil, err
	}
	if len(steps) != len(cmds) {
		return nil, fmt.Errorf("upstash: expected %d results, got %d", len(cmds), len(steps))
	}

	results := make([]Result, len(steps))
	for i, s := range steps {
		results[i] = Result{Value: s.Result}
		if s.Error != "" {
			results[i].Err = errors.New(s.Error)
		}
	}
	return results, nil
}

func ttlSeconds(ttl time.Duration) string {
	return strconv.FormatInt(int64(ttl.Seconds()), 10)
}

func (c *Client) Get(key string) (string, bool, error) {
	res, err := c.Pipeline([]string{"GET", key})
	if err != nil {
		return "", false, err
	}
	return res[0].String()
}

func (c *Client) GetInt(key string) (int64, error) {
	res, err := c.Pipeline([]string{"GET", key})
	if err != nil {
		return 0, err
	}
	return res[0].Int()
}

func (c *Client) SetEX(key, value string, ttl time.Duration) error {
	res, err := c.Pipeline([]string{"SET", key, value, "EX", ttlSeconds(ttl)})
	if err != nil {
		return err
	}
	return res[0].Err
}

func (c *Client) SetNX(key, value string, ttl time.Duration) (bool, error) {
	res, err := c.Pipeline([]string{"SET", key, value, "NX", "EX", ttlSeconds(ttl)})
	if err != nil {
		return false, err
	}
	if res[0].Err != nil {
		return false, res[0].Err
	}
	return !res[0].IsNull(), nil
}

func (c *Client) IncrByWithTTL(key string, n int64, ttl time.Duration) (int64, error) {
	res, err := c.Pipeline(
		[]string{"INCRBY", key, strconv.FormatInt(n, 10)},
		[]string{"EXPIRE", key, ttlSeconds(ttl), "NX"},
	)
	if err != nil {
		return 0, err
	}
	return res[0].Int()
}

func (c *Client) DecrBy(key string, n int64) error {
	res, err := c.Pipeline([]string{"DECRBY", key, strconv.FormatInt(n, 10)})
	if err != nil {
		return err
	}
	return res[0].Err
}

func (c *Client) Del(key string) error {
	res, err := c.Pipeline([]string{"DEL", key})
	if err != nil {
		return err
	}
	return res[0].Err
}
