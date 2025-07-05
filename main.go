package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yusufpapurcu/wmi"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

const serviceName string = "AOE4CardinalKiller"
const mainProcess string = "RelicCardinal.exe"
const secondaryProcess string = "RelicCardinal.exe"

type KillerService struct{}

func (m *KillerService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	var ticker time.Ticker = *time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var isMainStarted bool = false
mainLoop:
	for {
		select {
		case <-ticker.C:
			var processes []Win32_Process
			var query string = fmt.Sprintf("SELECT * FROM Win32_Process WHERE Name = '%s'", mainProcess)

			err := wmi.Query(query, &processes)
			if err != nil {
				log.Fatal(err)
			}

			if !isMainStarted && len(processes) != 0 && processes[0].ThreadCount >= 60 {
				isMainStarted = true
				fmt.Printf("Process %s started\n", mainProcess)
				continue
			}

			if isMainStarted && processes[0].ThreadCount < 60 {
				isMainStarted = false
				terminateProcessesByName(secondaryProcess)
				fmt.Printf("Process %s stopped\n", mainProcess)
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

	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService() error {
	isWindowsService, err := svc.IsWindowsService()
	if err != nil {
		return err
	}

	if isWindowsService {
		return svc.Run(serviceName, &KillerService{})
	}

	elog, err := eventlog.Open(serviceName)
	if err != nil {
		return err
	}

	elog.Info(1, "Open in console mode")
	debug.Run(serviceName, &KillerService{})
	return nil
}

type Win32_Process struct {
	Name        string
	ProcessID   uint32
	ThreadCount uint32
}

func terminateProcess(pid uint32) error {
	_, err := wmi.CallMethod(nil, fmt.Sprintf("Win32_Process.Handle=\"%d\"", pid), "Terminate", nil)
	return err
}

func terminateProcessesByName(name string) {
	var processes []Win32_Process
	query := fmt.Sprintf("SELECT * FROM Win32_Process WHERE Name = '%s'", name)

	err := wmi.Query(query, &processes)
	if err != nil {
		log.Fatal(err)
	}

	for _, process := range processes {
		terminateProcess(process.ProcessID)
	}
}

func main() {
	if err := runService(); err != nil {
		log.Fatalf("Service error: %v", err)
	}
}
