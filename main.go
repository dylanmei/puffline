package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/mgutz/ansi"
)

const PathStyle = "250+h:238"
const GitStyle = "250:22"
const KubernetesStyle = "250:21"
const AWSStyle = "250:130"
const TerraformStyle = "250:93"
const ExitCodeStyle = "250:124"

func main() {
	flag.Parse()

	cmd := exec.Command("pwd")
	stdout, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	pwd := strings.TrimSpace(string(stdout))

	fmt.Fprintln(os.Stdout)
	path(os.Stdout, pwd)
	git(os.Stdout, pwd)
	kubernetes(os.Stdout)
	aws(os.Stdout)
	terraform(os.Stdout, pwd)

	fmt.Fprintln(os.Stdout)
	exitCode(os.Stdout, flag.Args())
	prompt(os.Stdout)

	//colors(os.Stdout)
}

func path(w io.Writer, pwd string) {
	path := pwd
	if user, err := user.Current(); err == nil {
		if strings.HasPrefix(path, user.HomeDir) {
			path = "~" + path[len(user.HomeDir):]
		}
	}

	pathSegments := strings.Split(path, "/")
	viewSegments := []string{}
	for i := 0; i < len(pathSegments)-1; i++ {
		pathSegment := pathSegments[i]
		viewSegments = append(viewSegments, pathSegment[0:1])
	}

	viewSegments = append(viewSegments, pathSegments[len(pathSegments)-1])
	fmt.Fprint(w, ansi.Color(fmt.Sprintf(" %s ", strings.Join(viewSegments, "/")), PathStyle))
	return
}

func git(w io.Writer, pwd string) {
	cmd := exec.Command("git", "-C", pwd, "rev-parse", "--git-dir")
	stdout, _ := cmd.Output()

	if len(stdout) == 0 {
		return
	}

	cmd = exec.Command("git", "-C", pwd, "symbolic-ref", "HEAD")
	stdout, err := cmd.Output()

	if err != nil {
		cmd = exec.Command("git", "-C", pwd, "describe", "--tags", "--exact-match", "HEAD")
		stdout, err = cmd.Output()
	}

	if err != nil {
		cmd = exec.Command("git", "-C", pwd, "rev-parse", "--short", "HEAD")
		stdout, err = cmd.Output()
	}

	if err != nil {
		fmt.Fprint(w, strings.TrimRight(err.Error(), "\n"))
		return
	}

	if ref := strings.Replace(string(stdout), "refs/heads/", "", 1); ref != "" {
		cmd = exec.Command("git", "-C", pwd, "status", "--porcelain")
		stdout, err = cmd.Output()

		if len(stdout) == 0 {
			fmt.Fprint(w, ansi.Color(fmt.Sprintf(" %s ", strings.TrimRight(ref, "\n")), GitStyle))
		} else {
			fmt.Fprint(w, ansi.Color(fmt.Sprintf(" %s + ", strings.TrimRight(ref, "\n")), GitStyle))
		}
	}
}

func kubernetes(w io.Writer) {
	cmd := exec.Command("kubectl", "config", "current-context")
	stdout, _ := cmd.Output()

	if len(stdout) == 0 {
		return
	}

	context := strings.TrimRight(string(stdout), "\n")
	fmt.Fprintf(w, ansi.Color(fmt.Sprintf(" %s ", context), KubernetesStyle))
}

func aws(w io.Writer) {
	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		fmt.Fprintf(w, ansi.Color(fmt.Sprintf(" %s ", profile), AWSStyle))
	}
}

func terraform(w io.Writer, pwd string) {
	cmd := exec.Command("terraform", "workspace", "show")
	stdout, _ := cmd.Output()

	if len(stdout) == 0 {
		return
	}

	if workspace := strings.TrimRight(string(stdout), "\n"); workspace != "default" {
		fmt.Fprintf(w, ansi.Color(fmt.Sprintf(" %s ", workspace), TerraformStyle))
	}
}

func exitCode(w io.Writer, args []string) {
	if len(args) > 0 && args[0] != "0" {
		fmt.Fprint(w, ansi.Color(fmt.Sprintf(" %s ", args[0]), ExitCodeStyle))
	}
}

func prompt(w io.Writer) {
	fmt.Fprint(w, " % ")
}

func colors(w io.Writer) {
	for i := 0; i < 255; i++ {
		fmt.Fprintln(w, ansi.Color(fmt.Sprintf(" color %d ", i), fmt.Sprintf("white:%d", i)))
	}
}
