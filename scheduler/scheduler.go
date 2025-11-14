package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/Unfield/Statuz/logger"
	"github.com/Unfield/Statuz/monitors"
)

type Scheduler struct {
	ctx           context.Context
	monitors      []monitors.Monitor
	ResultChannel chan monitors.Result
	logger        logger.Logger
}

func NewScheduler(ctx context.Context, ms []monitors.Monitor) *Scheduler {
	return &Scheduler{
		ctx:           ctx,
		monitors:      ms,
		ResultChannel: make(chan monitors.Result, 100),
		logger:        logger.NewLogger(),
	}
}

func (s *Scheduler) Run() {
	for _, m := range s.monitors {
		go s.runMonitorLoop(m)
	}

	<-s.ctx.Done()
	log.Println("Scheduler stopped.")
}

func (s *Scheduler) runMonitorLoop(m monitors.Monitor) {
	ticker := time.NewTicker(m.GetHBInterval())
	defer ticker.Stop()

	s.logger.Infof("📡 [Scheduler] Started monitor %v (interval: %v)", m.GetID(), m.GetHBInterval())

	s.performCheck(m)

	for {
		select {
		case <-ticker.C:
			if !m.IsRunning() {
				go s.performCheck(m)
			}
		case <-s.ctx.Done():
			s.logger.Infof("[Scheduler] Stopping monitor %v", m.GetID())
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

func (s *Scheduler) Context() context.Context {
	return s.ctx
}
