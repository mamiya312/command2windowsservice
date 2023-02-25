package main

import (
	"bufio"
	"context"
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
	"golang.org/x/sys/windows/svc"
)

type service struct {
	stopCh chan<- bool
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s.stopCh <- true
				break loop
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func serve(ctx context.Context, name string, command string, arguments []string) error {
	cmd := exec.CommandContext(ctx, command, arguments...)
	if err := cmd.Run(); err != nil {
		return err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		sc := bufio.NewScanner(stdoutPipe)
		for sc.Scan() {
			log.Println(sc.Text())
		}
	}()
	go func() {
		sc := bufio.NewScanner(stderrPipe)
		for sc.Scan() {
			log.Println(sc.Text())
		}
	}()
	if err := cmd.Wait(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

var StopCh = make(chan bool)

func main() {
	app := &cli.App{
		Name:    "command2windowsservice",
		Version: getVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Required: true},
		},
		Action: func(c *cli.Context) error {
			name := c.String("name")

			isService, err := svc.IsWindowsService()
			if err != nil {
				return err
			}

			command := ""
			arguments := []string{}
			values := c.Args().Slice()
			if len(values) == 0 {
				return errors.New("no command")
			}
			if len(values) == 1 {
				command = values[0]
			} else {
				command = values[0]
				arguments = values[1:]
			}

			if isService {
				go func() {
					err = svc.Run(name, &service{stopCh: StopCh})
					if err != nil {
						log.Fatalf("Failed to start %s service: %v", name, err.Error())
					}
				}()
			}
			ctx, cancel := context.WithCancel(context.Background())

			go func() {
				if err := serve(ctx, name, command, arguments); err != nil {
					log.Fatalf("cannot start %s: %s", name, err.Error())
				}
			}()

			for {
				if <-StopCh {
					cancel()
					log.Printf("Shutting down %s", name)
					break
				}
			}

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
