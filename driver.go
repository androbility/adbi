package adbi

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Commander struct {
	cmd        *exec.Cmd
	in         io.WriteCloser
	lastActive time.Time
	m          *sync.Mutex
	stopCh     chan interface{}
}

func New() (*Commander, error) {
	cmd := exec.Command("adb", "shell")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	cmndr := &Commander{
		cmd:        cmd,
		in:         stdin,
		lastActive: time.Now(),
		m:          &sync.Mutex{},
		stopCh:     make(chan interface{}),
	}
	//	go cmndr.ping(179 * time.Second)

	return cmndr, nil
}

func (c *Commander) Signal(key Keyevent) error {
	return c.SignalWithRepeat(key, 1)
}

func (c *Commander) SignalWithRepeat(key Keyevent, n int) error {
	inputEvent := key.TriggerWithRepeat(n)
	if _, err := c.in.Write(inputEvent); err != nil {
		// Communication with the Android device failed.
		log.WithFields(log.Fields{
			"error": err,
			"key":   rune(key),
		}).Error("KeyEvent send failed")

		// We can assume the server is down, or restarting.
		// Let's return an error, kill cmd, and close the channel.
		defer close(c.stopCh)
		defer c.cmd.Wait()

		return errors.New("server connection lost")
	}

	log.Info(strings.Trim(string(inputEvent), "\n"))

	return nil
}

func (c *Commander) Raw(cmd string) error {
	if _, err := c.in.Write([]byte(fmt.Sprintf("input %s\n", cmd))); err != nil {
		// Communication with the Android device failed.
		log.WithFields(log.Fields{
			"error": err,
			"key":   rune('\x00'),
		}).Error(fmt.Sprintf("%s send failed", cmd))

		// We can assume the server is down, or restarting.
		// Let's return an error, kill cmd, and close the channel.
		defer close(c.stopCh)
		defer c.cmd.Wait()

		return errors.New("server connection lost")
	}

	log.Info(strings.Trim(cmd, "\n"))

	return nil
}

func (c *Commander) Quit() {
	c.in.Write([]byte("exit\n"))
	log.Info("Quitting")
	c.cmd.Wait()
	os.Exit(0)
}

// Ensure device stays awake.  Purpose: FireTV Cube.
func (c *Commander) ping(dur time.Duration) {
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if time.Since(c.lastActive) < dur {
				break
			}
			c.m.Lock()
			c.Signal(Keyevent('w'))
			c.m.Unlock()
		case <-c.stopCh:
			return
		}
	}
}

// Wait for server to reconnect.
func WaitForAndroid() {
	log.Info("Waiting for a new adb connection.  (Hint: adb connect <ip[:port]>)")
	exec.Command("adb", "wait-for-device").Run()
	log.Info("Success!")
}
