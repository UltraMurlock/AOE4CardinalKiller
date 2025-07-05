package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
)

type KillerService struct{}

func (m *KillerService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	serviceLog.Info(1, "Service started")

	var ticker *time.Ticker = time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var isMainStarted bool = false
mainLoop:
	for {
		select {
		case <-ticker.C:
			processes := getProcessesByName(aoe4process)

			if !isMainStarted && len(processes) != 0 && processes[0].ThreadCount >= 60 {
				serviceLog.Info(1, "Game started")
				isMainStarted = true
				continue
			}

			if isMainStarted && processes[0].ThreadCount < 60 {
				serviceLog.Info(1, "Game stopped")
				isMainStarted = false
				err := terminateProcess(processes[0])

				if err != nil {
					serviceLog.Error(1, fmt.Sprintf("Process termination error: %v", err))
				} else {
					serviceLog.Info(1, "Process terminated successfully")
				}
			}

		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break mainLoop
			}
		}
	}

	changes <- svc.Status{State: svc.State(svc.StopPending)}
	changes <- svc.Status{State: svc.Stopped}
	serviceLog.Info(1, "Service stopped")
	return
}
