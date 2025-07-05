package main

import (
	"fmt"
	"log"

	"github.com/yusufpapurcu/wmi"
)

type Win32_Process struct {
	Name        string
	ProcessID   uint32
	ThreadCount uint32
}

func getProcessesByName(name string) []Win32_Process {
	var processes []Win32_Process
	var query string = fmt.Sprintf("SELECT Name, ProcessID, ThreadCount FROM Win32_Process WHERE Name = '%s'", name)

	err := wmi.Query(query, &processes)
	if err != nil {
		log.Fatal(err)
	}

	return processes
}

func terminateProcess(process Win32_Process) error {
	_, err := wmi.CallMethod(nil, fmt.Sprintf("Win32_Process.Handle=\"%d\"", process.ProcessID), "Terminate", nil)
	return err
}
