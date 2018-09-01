package main

import (
  "log"
  "fmt"
  "os"
  "os/exec"
  "context"
  "bufio"
  "io"
  "encoding/json"
)

type Events struct {
  Status  string  `json:"status,omitempty"`
  ID	  string  `json:"id,omitempty"`
  From	  string  `json:"from,omitempty"`
  Type	  string  `json:"Type,omitempty"`
  Action  string  `json:"Action,omitempty"`
  Actor	  Actor	  `json:"Actor,omitempty"`
}

type Actor struct {
  ID	      string	  `json:ID,omitempty`
  Attributes  Attributes  `json:"Attributes,omitempty"`
}

type Attributes struct {
  Image string	`json:"image,omitempty"`
  Name	string	`json:"name,omitempty"`
}

func events(ctx context.Context, event chan<- []byte) error {
  var (
    cmd	    *exec.Cmd
    err	    error
    stdout  io.ReadCloser
    scanner *bufio.Scanner
  )

  cmd = exec.CommandContext(ctx, "docker", "events", "--format", "{{json .}}")
  cmd.Stderr = os.Stderr

  if stdout, err = cmd.StdoutPipe(); err != nil {
    return err
  }

  scanner = bufio.NewScanner(bufio.NewReader(stdout))

  go func() {
    for scanner.Scan() {
      event <- scanner.Bytes()
    }
  }()

  if err = cmd.Start(); err != nil {
    return err
  }

  if err = scanner.Err(); err != nil {
    return err
  }

  cmd.Wait()

  return nil
}

func main() {
  var (
    err	error
    errc = make(chan error, 1)
    event = make(chan []byte)
  )

  go func() {
    if err = events(context.Background(), event); err != nil {
      errc <- err
    }
  }()

  for {
    select {
    case msg := <-event:
      var ev Events

      if err = json.Unmarshal(msg, &ev); err != nil {
	log.Fatal(err)
      }

      fmt.Println(string(msg))
      fmt.Println(ev)
    case e := <-errc:
      log.Fatal(e)
    }
  }

  fmt.Println("vim-go")
}
