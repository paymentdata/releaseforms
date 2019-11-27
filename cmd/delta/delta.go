package main

import (
        "bufio"
        "encoding/gob"
        "io"
        "log"
        "os"
        "os/exec"
        "strconv"
)

//git log $lastDeploymentSHA.. --merges --oneline --first-parent | grep -oe '#[0-9]*' | tr -d "#" | tr '\n' ','
func main() {
        var (
                fromSHA, toSHA string
                err            error

                gd = gob.NewEncoder(os.Stdout)
        )

        if len(os.Args) == 1 {
                os.Stdout.Write([]byte("Provide at least a starting SHA for calculating Δ\n"))
                os.Exit(1)
        }

        fromSHA = os.Args[1]
        if len(fromSHA) == 0 {
                os.Exit(1)
        }

        if len(os.Args) == 3 {
                toSHA = os.Args[2]
        }

		//this is definitely a candidate for process refinement, kind of delegating gopher duties to the shell here.
        cmd := exec.Command("sh", "-c", "git log "+fromSHA+".."+toSHA+" --merges --oneline --first-parent | grep -oe '#[0-9]*' | tr -d \"#\" | tr '\n' ','")

        var (
                prIDsource   *bufio.Reader
                stdoutStream io.ReadCloser
                str          string //temp string for PRIDs to cast into int before sending along.
                num          int
        )

        stdoutStream, err = cmd.StdoutPipe()
        if err != nil {
                panic(err)
        }

        prIDsource = bufio.NewReader(stdoutStream)

        if err = cmd.Start(); err != nil {
                panic(err)
        }

        for {
                str, err = prIDsource.ReadString(',')
                if err != nil {
                        if err != io.EOF {
                                panic(err)
                        }
                        break
                }
                num, err = strconv.Atoi(str[:len(str)-1])
                if err != nil {
                        log.Fatalf("expected PullReleaseID, got %s instead which caused err: %+v", str[:len(str)-1], err)
                        continue
                }
                err = gd.Encode(num)
                if err != nil {
                        log.Printf("Err: %+v", err)
                }
        }
}