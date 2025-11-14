package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/Unfield/Statuz/logger"
	"github.com/Unfield/Statuz/monitors"
)

type Scheduler struct {
	ctx           context.Context
	ResultChannel chan monitors.Result
	logger        logger.Logger

	mu       sync.RWMutex
	monitors map[string]monitors.Monitor
}

func NewScheduler(ctx context.Context, ms []monitors.Monitor) *Scheduler {
	m := make(map[string]monitors.Monitor, len(ms))
	for _, mon := range ms {
		m[mon.GetID()] = mon
	}

	return &Scheduler{
		ctx:           ctx,
		ResultChannel: make(chan monitors.Result, 100),
		logger:        logger.NewLogger(),
		monitors:      m,
	}
}

func (s *Scheduler) Run() {
	s.logger.Info("[Scheduler] Running dynamic scheduler...")

	s.mu.RLock()
	for _, m := range s.monitors {
		m.Start(s.ctx)
		go s.runMonitorLoop(m)
	}
	s.mu.RUnlock()

	<-s.ctx.Done()
	s.logger.Warn("[Scheduler] Context canceled, stopping all monitors...")

	s.mu.Lock()
	for id, m := range s.monitors {
		m.Stop()
		delete(s.monitors, id)
	}
	s.mu.Unlock()

	s.logger.Info("Scheduler stopped cleanly.")
}

func (s *Scheduler) AddMonitor(m monitors.Monitor) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.monitors[m.GetID()]; exists {
		s.logger.Warnf("[Scheduler] Monitor %s already exists", m.GetID())
		return
	}

	m.Start(s.ctx)
	s.monitors[m.GetID()] = m

	go s.runMonitorLoop(m)
	s.logger.Infof("[Scheduler] Added and started monitor %s", m.GetID())
}

func (s *Scheduler) StopMonitor(monitorID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, exists := s.monitors[monitorID]
	if !exists {
		s.logger.Warnf("[Scheduler] No monitor found with ID %s", monitorID)
		return
	}

	m.Stop()
	delete(s.monitors, monitorID)

	s.logger.Infof("[Scheduler] Stopped and removed monitor %s", monitorID)
}

func (s *Scheduler) ListMonitors() []monitors.Monitor {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]monitors.Monitor, 0, len(s.monitors))
	for _, m := range s.monitors {
		list = append(list, m)
	}
	return list
}

func (s *Scheduler) Context() context.Context {
	return s.ctx
}

func (s *Scheduler) runMonitorLoop(m monitors.Monitor) {
	ticker := time.NewTicker(m.GetHBInterval())
	defer ticker.Stop()

	s.logger.Infof("📡 [Scheduler] Started monitor %v (interval: %v)",
		m.GetID(), m.GetHBInterval())

	s.performCheck(m)

	for {
		select {
		case <-ticker.C:
			if !m.IsRunning() {
				go s.performCheck(m)
			}
		case <-m.GetContext().Done():
			s.logger.Infof("[Scheduler] Monitor %v stopped", m.GetID())
			return
		case <-s.ctx.Done():
			s.logger.Infof("[Scheduler] Monitor %v stopping (global shutdown)", m.GetID())
			return
		}
	}
}

func (s *Scheduler) performCheck(m monitors.Monitor) {
	m.SetRunning(true)
	defer m.SetRunning(false)

	res := m.Check(s.ctx)
	m.SetLastHB(res.EndTime)
	s.ResultChannel <- res
}
