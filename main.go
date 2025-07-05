package main

import (
	"log"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

const serviceName string = "AOE4CardinalKiller"
const aoe4process string = "RelicCardinal.exe"

var serviceLog *eventlog.Log

func main() {
	elog, err := eventlog.Open(serviceName)
	if err != nil {
		log.Fatalf("Event log open error: %v", err)
	}

	serviceLog = elog

	if err := runService(); err != nil {
		log.Fatalf("Service error: %v", err)
	}
}

func runService() error {
	isWindowsService, err := svc.IsWindowsService()
	if err != nil {
		return err
	}

	if isWindowsService {
		return svc.Run(serviceName, &KillerService{})
	} else {
		return debug.Run(serviceName, &KillerService{})
	}
}
