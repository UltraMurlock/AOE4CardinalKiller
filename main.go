package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yusufpapurcu/wmi"
)

type Win32_Process struct {
	Name      string
	ProcessID uint32
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
	var mainProcess string = "notepad.exe"
	var secondaryProcess string = "notepad++.exe"

	var isMainStarted bool = false
	for {
		var processes []Win32_Process
		var query string = fmt.Sprintf("SELECT * FROM Win32_Process WHERE Name = '%s'", mainProcess)
		err := wmi.Query(query, &processes)
		if err != nil {
			log.Fatal(err)
		}

		if len(processes) != 0 {
			isMainStarted = true
			fmt.Printf("Process %s started\n", mainProcess)
		}

		if isMainStarted && len(processes) == 0 {
			isMainStarted = false
			terminateProcessesByName(secondaryProcess)
			fmt.Printf("Process %s stopped\n", mainProcess)
		}

		var sleepDuration time.Duration
		if isMainStarted {
			sleepDuration = 1 * time.Second
		} else {
			sleepDuration = 5 * time.Second
		}

		time.Sleep(sleepDuration)
	}
}
