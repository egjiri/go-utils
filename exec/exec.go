package exec

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/egjiri/go-kit/ui/print"
)

// ExecuteWithHeading displays a message and then executes the command
func ExecuteWithHeading(heading, command string) {
	print.Heading(heading)
	print.Command(command)
	ExecuteBash(command)
}

// ExecuteBash packages the passed command as a bash file and executes it
func ExecuteBash(command string) {
	tmpfile, err := ioutil.TempFile("", "cli")
	defer os.Remove(tmpfile.Name()) // clean up
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(command)); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	Execute("/bin/bash", tmpfile.Name())
}

// Execute runs a command directly on the underlying OS
func Execute(commands ...string) error {
	args := buildCommandArgs(commands...)
	name, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	p, err := os.StartProcess(name, args, &procAttr)
	if err != nil {
		return err
	}
	ps, err := p.Wait()
	if err != nil {
		return err
	}
	if !ps.Success() {
		return errors.New("command execution failed")
	}
	return nil
}

func buildCommandArgs(commands ...string) []string {
	var args []string
	re := regexp.MustCompile("^\"([^\"]*)\"$")
	for _, c := range commands {
		if re.MatchString(c) { // command does not get split if it includes double quotes
			args = append(args, re.ReplaceAllString(c, "$1"))
		} else { // split the command as arguments on spaces
			for _, sub := range strings.Split(c, " ") {
				args = append(args, sub)
			}
		}
	}
	return args
}
