package monitors

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	DEGRADED_DEFAULT_THRESHHOLD = 2000 * time.Millisecond
)

// TODO: add zap for logging

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
	lastHB              time.Time
	running             bool
}

func (m *HTTPMonitor) GetHBInterval() time.Duration {
	return m.HBInterval
}

func (m *HTTPMonitor) GetLastHB() time.Time {
	return m.lastHB
}

func (m *HTTPMonitor) SetLastHB(hbTime time.Time) {
	m.lastHB = hbTime
}

func (m *HTTPMonitor) IsRunning() bool {
	return m.running
}

func (m *HTTPMonitor) SetRunning(b bool) {
	m.running = b
}

func (m *HTTPMonitor) GetID() string {
	return m.ID
}

func (m *HTTPMonitor) Check(ctx context.Context) Result {
	currentRetries := 0
PerformCheck:
	if currentRetries >= m.Retries+1 {
		return Result{
			MonitorID: m.ID,
			Type:      "http/https",
			Status:    StatusDown,
			Success:   false,
			Error:     "max retries reached",
		}
	}

	checkResult, err := m.performCheck()
	if err != nil {
		currentRetries++
		time.Sleep(m.RetryInterval)
		goto PerformCheck
	}
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
		m.DegradedThreshold = DEGRADED_DEFAULT_THRESHHOLD
	}
	if checkResult.duration >= m.DegradedThreshold {
		res.Status = StatusDegraded
		res.Message = fmt.Sprintf("Response slow: %v > %v", checkResult.duration, m.DegradedThreshold)
	}
	return res
}

type CheckResult struct {
	isUp     bool
	start    time.Time
	end      time.Time
	duration time.Duration
	code     int
}

func (m *HTTPMonitor) performCheck() (CheckResult, error) {
	m.HTTPMethod = strings.ToUpper(m.HTTPMethod)

	var client *http.Client

	switch m.IPFamily {
	case "v4":
		dialer := &net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 15 * time.Second,
		}

		transport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, "tcp4", addr)
			},
		}

		client = &http.Client{
			Transport: transport,
			Timeout:   m.ReqTimeout * time.Second,
		}
	case "v6":

	default:
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	var reqRes *http.Response
	var reqErr error

	var startTime time.Time

	switch m.HTTPMethod {
	case "GET":
		startTime = time.Now()
		reqRes, reqErr = client.Get(m.URL)
	case "POST":
		startTime = time.Now()
		reqRes, reqErr = client.Post(m.URL, "text/plain", nil)
	default:
		reqErr = fmt.Errorf("invalid method")
	}

	if reqErr != nil {
		return CheckResult{
				isUp:     false,
				start:    startTime,
				end:      startTime,
				duration: 0,
				code:     0,
			},
			fmt.Errorf("failed to make get request: %w", reqErr)
	}
	defer reqRes.Body.Close()
	endTime := time.Now()

	ok := isStatusAccepted(reqRes.StatusCode, m.AcceptedStatusCodes)
	if !ok {
		return CheckResult{
				isUp:     false,
				start:    startTime,
				end:      endTime,
				duration: endTime.Sub(startTime),
				code:     reqRes.StatusCode,
			},
			fmt.Errorf("service responded with unaccepted status")
	}
	return CheckResult{
			isUp:     true,
			start:    startTime,
			end:      endTime,
			duration: endTime.Sub(startTime),
			code:     reqRes.StatusCode,
		},
		nil
}

func isStatusAccepted(status int, accepted []int) bool {
	for _, acceptedStatus := range accepted {
		if status == acceptedStatus {
			return true
		}
		continue
	}
	return false
}
