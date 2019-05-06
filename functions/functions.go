package functions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Record(filename string, duration int) error {
	if duration > 5*60 {
		fmt.Println("duration too lengthy")
		duration = 3
	}
	cmd := exec.Command("rec", "-r", "160000", "-c", "1", filename, "trim", "0", strconv.Itoa(duration)) //

	env := os.Environ()
	// env = append(env, "AUDIODEV=hw:1,0")
	cmd.Env = env
	fmt.Printf("TB: recording audio for %d s...\n", duration)
	return cmd.Run()
}

func Snap(filename string) (error) {
	cmd := exec.Command("fswebcam", "--no-banner", filename)
	return cmd.Run()
}

func Ip() (string, error) {
	cmd := exec.Command("hostname", "-I")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error running terminal cmd: %v", err)
		return "", err
	}
	fmt.Printf("IP output: %v", string(output))
	split := strings.Split(string(output), " ")
	if len(split) > 0 {
		return split[0], nil
	}
	/*return func() string {
		if len(split) > 0 {
			return split[0]
		}
		return ""
	}(), nil*/
	return "", nil
}

