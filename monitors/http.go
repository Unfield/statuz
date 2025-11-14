package monitors

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	DEGRADED_DEFAULT_THRESHOLD = 2000 * time.Millisecond
)

type HTTPMonitor struct {
	ID                  string        `json:"id"`
	URL                 string        `json:"url"`
	HBInterval          time.Duration `json:"hb_interval"`
	Retries             int           `json:"retries"`
	RetryInterval       time.Duration `json:"retry_interval"`
	ReqTimeout          time.Duration `json:"req_timeout"`
	MaxRedirects        int           `json:"max_redirects"`
	AcceptedStatusCodes []int         `json:"accepted_status_codes"`
	IPFamily            string        `json:"ip_family"`
	HTTPMethod          string        `json:"http_method"`
	DegradedThreshold   time.Duration `json:"degraded_threshold"`

	lastHB  time.Time
	running bool

	ctx    context.Context
	cancel context.CancelFunc
	client *http.Client
	mu     sync.Mutex
}

func (m *HTTPMonitor) Start(parent context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cancel != nil {
		return
	}

	m.ctx, m.cancel = context.WithCancel(parent)
	m.initHTTPClient()
}

func (m *HTTPMonitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}
}

func (m *HTTPMonitor) GetContext() context.Context {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.ctx
}

func (m *HTTPMonitor) GetHBInterval() time.Duration {
	return m.HBInterval
}

func (m *HTTPMonitor) GetLastHB() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastHB
}

func (m *HTTPMonitor) SetLastHB(hbTime time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastHB = hbTime
}

func (m *HTTPMonitor) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.running
}

func (m *HTTPMonitor) SetRunning(b bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.running = b
}

func (m *HTTPMonitor) GetID() string {
	return m.ID
}

func (m *HTTPMonitor) Check(ctx context.Context) Result {
	if m.ctx == nil {
		m.ctx = context.Background()
	}

	retries := max(0, m.Retries)

	for attempt := 0; attempt <= retries; attempt++ {
		select {
		case <-m.ctx.Done():
			return Result{
				MonitorID: m.ID,
				Type:      "http/https",
				Status:    StatusDown,
				Success:   false,
				Error:     "monitor cancelled",
			}
		default:
		}

		checkResult, err := m.performCheck()
		if err == nil {
			res := Result{
				MonitorID: m.ID,
				Type:      "http/https",
				Status:    StatusUp,
				StartTime: checkResult.start,
				EndTime:   checkResult.end,
				Duration:  checkResult.duration,
				Success:   checkResult.isUp,
				CheckedAt: checkResult.start,
			}

			if m.DegradedThreshold <= 0 {
				m.DegradedThreshold = DEGRADED_DEFAULT_THRESHOLD
			}
			if checkResult.duration >= m.DegradedThreshold {
				res.Status = StatusDegraded
				res.Message = fmt.Sprintf(
					"Response slow: %v > %v",
					checkResult.duration,
					m.DegradedThreshold,
				)
			}
			return res
		}

		// Retry logic
		if attempt < retries {
			time.Sleep(m.RetryInterval)
			continue
		}

		// Final failure
		return Result{
			MonitorID: m.ID,
			Type:      "http/https",
			Status:    StatusDown,
			Success:   false,
			Error:     fmt.Sprintf("check failed after %d retries: %v", retries, err),
		}
	}

	// Should never reach this point
	return Result{
		MonitorID: m.ID,
		Type:      "http/https",
		Status:    StatusDown,
		Success:   false,
		Error:     "unexpected error",
	}
}

// CheckResult stores the outcome of a single check execution.
type CheckResult struct {
	isUp     bool
	start    time.Time
	end      time.Time
	duration time.Duration
	code     int
}

// performCheck executes an HTTP request using the configured client and context.
func (m *HTTPMonitor) performCheck() (CheckResult, error) {
	method := strings.ToUpper(m.HTTPMethod)
	if method == "" {
		method = "GET"
	}

	req, err := http.NewRequestWithContext(m.ctx, method, m.URL, nil)
	if err != nil {
		return CheckResult{}, fmt.Errorf("failed to create request: %w", err)
	}

	start := time.Now()
	resp, err := m.client.Do(req)
	if err != nil {
		return CheckResult{
			isUp:     false,
			start:    start,
			end:      start,
			duration: 0,
			code:     0,
		}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	end := time.Now()
	duration := end.Sub(start)

	if !isStatusAccepted(resp.StatusCode, m.AcceptedStatusCodes) {
		return CheckResult{
				isUp:     false,
				start:    start,
				end:      end,
				duration: duration,
				code:     resp.StatusCode,
			},
			fmt.Errorf("unaccepted status code: %d", resp.StatusCode)
	}

	return CheckResult{
		isUp:     true,
		start:    start,
		end:      end,
		duration: duration,
		code:     resp.StatusCode,
	}, nil
}

// initHTTPClient builds an HTTP client with proper IP family and timeouts.
func (m *HTTPMonitor) initHTTPClient() {
	m.HTTPMethod = strings.ToUpper(m.HTTPMethod)
	if m.HTTPMethod == "" {
		m.HTTPMethod = "GET"
	}

	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 15 * time.Second,
	}

	var transport *http.Transport

	switch m.IPFamily {
	case "v4":
		transport = &http.Transport{
			DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, "tcp4", addr)
			},
		}
	case "v6":
		transport = &http.Transport{
			DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, "tcp6", addr)
			},
		}
	default:
		// Allow both IPv4 and IPv6
		transport = &http.Transport{
			DialContext: dialer.DialContext,
		}
	}

	timeout := m.ReqTimeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	m.client = &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

// isStatusAccepted returns true if the response code is within the accepted list.
func isStatusAccepted(status int, accepted []int) bool {
	if len(accepted) == 0 {
		return status >= 200 && status < 400
	}
	for _, acceptedStatus := range accepted {
		if status == acceptedStatus {
			return true
		}
	}
	return false
}
