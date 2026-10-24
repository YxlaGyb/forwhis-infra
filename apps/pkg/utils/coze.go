package utils

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	maxScanTokenSize  = 1024 * 1024
	maxRetries        = 3
	initialRetryDelay = 500 * time.Millisecond
	maxRetryDelay     = 5 * time.Second
)

type CozeClient struct {
	AccessToken string
	BaseURL     string
	BotID       string
	HTTPClient  *http.Client
	DebugMode   bool
}

func NewCozeClient(accessToken, botID string) *CozeClient {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxConnsPerHost: 50,
	}

	return &CozeClient{
		AccessToken: accessToken,
		BaseURL:     "https://api.coze.cn/v3/chat",
		BotID:       botID,
		DebugMode:   false,
		HTTPClient: &http.Client{
			Timeout:   120 * time.Second,
			Transport: transport,
		},
	}
}

func (c *CozeClient) SetDebugMode(enable bool) {
	c.DebugMode = enable
}

type ChatStreamRequest struct {
	Context        context.Context
	ConversationID string
	UserID         string
	Message        string
}

type StreamCallback func(event string, data string) error

func (c *CozeClient) ChatStream(req ChatStreamRequest, callback StreamCallback) error {
	payload := map[string]interface{}{
		"bot_id":  c.BotID,
		"user_id": req.UserID,
		"stream":  true,
		"additional_messages": []map[string]string{
			{
				"role":    "user",
				"content": req.Message,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request failed: %w", err)
	}

	apiURL := c.BaseURL
	if req.ConversationID != "" {
		apiURL = fmt.Sprintf("%s?conversation_id=%s", c.BaseURL, req.ConversationID)
	}

	httpReq, err := http.NewRequest("POST", apiURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	if req.Context != nil {
		httpReq = httpReq.WithContext(req.Context)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.AccessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.doWithRetry(httpReq)
	if err != nil {
		return fmt.Errorf("request failed after retries: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("api error: status=%d, failed to read body: %w", resp.StatusCode, readErr)
		}
		return fmt.Errorf("api error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, maxScanTokenSize)

	var currentEvent string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "event:") {
			currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if data == "" {
				continue
			}
			if currentEvent == "done" {
				continue
			}
			if err := callback(currentEvent, data); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func (c *CozeClient) doWithRetry(req *http.Request) (*http.Response, error) {
	var lastErr error
	var bodyBytes []byte

	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("read request body failed: %w", err)
		}
		req.Body.Close()
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := c.calculateRetryDelay(attempt)
			time.Sleep(delay)
		}

		if bodyBytes != nil {
			req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			if c.isRetryableError(err) {
				if c.DebugMode {
					fmt.Printf("[Coze Retry] attempt %d, error: %v\n", attempt+1, err)
				}
				continue
			}
			return nil, err
		}

		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			lastErr = fmt.Errorf("server error: status=%d", resp.StatusCode)
			resp.Body.Close()
			if c.DebugMode {
				fmt.Printf("[Coze Retry] attempt %d, status: %d\n", attempt+1, resp.StatusCode)
			}
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func (c *CozeClient) calculateRetryDelay(attempt int) time.Duration {
	delay := initialRetryDelay * time.Duration(1<<uint(attempt-1))
	if delay > maxRetryDelay {
		delay = maxRetryDelay
	}
	jitter := time.Duration(rand.Int63n(int64(delay / 2)))
	return delay + jitter
}

func (c *CozeClient) isRetryableError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, context.Canceled) {
		return false
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}

	return true
}

type CozeMessage struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type GetMessagesOptions struct {
	Limit int
}

func (c *CozeClient) GetMessages(conversationID string, opts *GetMessagesOptions) ([]CozeMessage, error) {
	limit := 50
	if opts != nil && opts.Limit > 0 {
		limit = opts.Limit
	}

	apiURL := fmt.Sprintf("https://api.coze.cn/v1/conversation/message/list?conversation_id=%s&order=asc&limit=%d", conversationID, limit)

	httpReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.doWithRetry(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.DebugMode {
		fmt.Printf("[Coze Messages API] status=%d\n", resp.StatusCode)
	}

	var result struct {
		Code int           `json:"code"`
		Msg  string        `json:"msg"`
		Data []CozeMessage `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("api error: code=%d, msg=%s", result.Code, result.Msg)
	}

	return result.Data, nil
}

func ParseConversationID(data string) string {
	var convData struct {
		ConversationID string `json:"conversation_id"`
	}
	if err := json.Unmarshal([]byte(data), &convData); err == nil {
		return convData.ConversationID
	}
	return ""
}

func ParseMessageContent(data string) string {
	var msgData struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(data), &msgData); err == nil {
		return msgData.Content
	}
	return ""
}

func ParseError(data string) string {
	var errData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal([]byte(data), &errData); err == nil {
		return fmt.Sprintf("code:%d msg:%s", errData.Code, errData.Msg)
	}
	return ""
}
