package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/natefinch/atomic"
	"github.com/urfave/cli/v2"
)

const codeFmt = `package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID(%q)
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}
`

const tomlFmt = `[Details]
Icon = "Icon.png"
Name = %q
ID = %q
Version = "0.0.1"
Build = 1
`

func Init() *cli.Command {
	return &cli.Command{
		Name:      "init",
		Usage:     "Initializes a new Fyne project.",
		ArgsUsage: "[module-path]",
		Action:    initAction,
		Description: "Initializes a new Fyne project in the current directory, including\n" +
			"a go.mod, main.go, and FyneApp.toml file (unless existing).",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "appID",
				Aliases: []string{"id"},
				Usage:   "set appID in reversed domain notation for Android, darwin and Windows targets, or a valid provisioning profile on iOS",
			},
			&cli.StringFlag{
				Name:        "name",
				Usage:       "set name the application",
				DefaultText: "executable file name",
			},
		},
	}
}

func getAppID(modpath string) string {
	p := strings.Split(modpath, "/")
	if len(p) == 0 {
		return ""
	}

	d := strings.Split(p[0], ".")
	r := make([]string, len(p)+len(d)-1)
	for n, e := range d {
		r[len(d)-n-1] = e
	}
	for n, e := range p {
		if n == 0 {
			continue
		}
		r[len(d)+n-1] = e
	}

	return strings.Join(r, ".")
}

func getAppName(modpath string) string {
	p := strings.Split(modpath, "/")
	if len(p) == 0 {
		return ""
	}

	if len(p) > 1 {
		return p[len(p)-1]
	}

	d := strings.Split(p[0], ".")

	return d[0]
}

func checkFileOrDo(file string, cb func() error) error {
	fi, err := os.Stat(file)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if fi != nil {
		return nil
	}

	return cb()
}

func checkFileOrCreate(file, content string) error {
	return checkFileOrDo(file, func() error {
		if err := atomic.WriteFile(file, strings.NewReader(content)); err != nil {
			return err
		}
		return os.Chmod(file, 0644)
	})
}

func initAction(ctx *cli.Context) error {
	modpath := ctx.Args().Get(0)
	appID := ctx.String("appID")
	appName := ctx.String("name")

	if modpath == "" {
		modpath = "example"

		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		if wd != "" && wd != "." {
			modpath = filepath.Base(wd)
		}
	}

	if appID == "" {
		appID = getAppID(modpath)
	}

	if appName == "" {
		appName = getAppName(modpath)
	}

	if err := checkFileOrCreate("main.go", fmt.Sprintf(codeFmt, appID)); err != nil {
		return err
	}

	if err := checkFileOrDo("go.mod", func() error {
		cmd := exec.Command("go", "mod", "init", modpath)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		return cmd.Run()
	}); err != nil {
		return err
	}

	if err := checkFileOrCreate("FyneApp.toml", fmt.Sprintf(tomlFmt, appName, appID)); err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}

	return nil
}
