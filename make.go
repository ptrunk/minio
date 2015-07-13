// +build ignore

/*
 * Makefile alternative for Minimalist Object Storage, (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"
)

type Version struct {
	Date string
	Tag  string
}

func writeVersion(version Version) error {
	var versionTemplate = `// --------  DO NOT EDIT --------
// this is an autogenerated file

package main

import (
	"net/http"
	"time"
)

// Version autogenerated
var Version = {{if .Date}}"{{.Date}}"{{else}}""{{end}}

// getVersion -
func getVersion() string {
	t, _ := time.Parse(time.RFC3339Nano, Version)
	if t.IsZero() {
		return ""
	}
	return t.Format(http.TimeFormat)
}
`
	t := template.Must(template.New("version").Parse(versionTemplate))
	versionFile, err := os.OpenFile("version.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer versionFile.Close()
	err = t.Execute(versionFile, version)
	if err != nil {
		return err
	}
	return nil
}

type command struct {
	cmd    *exec.Cmd
	stderr *bytes.Buffer
	stdout *bytes.Buffer
}

func (c command) runCommand() error {
	c.cmd.Stdout = c.stdout
	c.cmd.Stderr = c.stderr
	return c.cmd.Run()
}
func (c command) String() string {
	message := c.stderr.String()
	message += c.stdout.String()
	return message
}

func runMinioInstall() {
	minioGenerate := command{exec.Command("godep", "go", "generate", "./..."), &bytes.Buffer{}, &bytes.Buffer{}}
	minioBuild := command{exec.Command("godep", "go", "build", "-a", "./..."), &bytes.Buffer{}, &bytes.Buffer{}}
	minioTest := command{exec.Command("godep", "go", "test", "-race", "./..."), &bytes.Buffer{}, &bytes.Buffer{}}
	minioInstall := command{exec.Command("godep", "go", "install", "-a", "github.com/minio/minio"), &bytes.Buffer{}, &bytes.Buffer{}}
	minioGenerateErr := minioGenerate.runCommand()
	if minioGenerateErr != nil {
		fmt.Print(minioGenerate)
		os.Exit(1)
	}
	fmt.Print(minioGenerate)
	minioBuildErr := minioBuild.runCommand()
	if minioBuildErr != nil {
		fmt.Print(minioBuild)
		os.Exit(1)
	}
	fmt.Print(minioBuild)
	minioTestErr := minioTest.runCommand()
	if minioTestErr != nil {
		fmt.Println(minioTest)
		os.Exit(1)
	}
	fmt.Print(minioTest)
	minioInstallErr := minioInstall.runCommand()
	if minioInstallErr != nil {
		fmt.Println(minioInstall)
		os.Exit(1)
	}
	fmt.Print(minioInstall)
}

func runMinioRelease() {
	t := time.Now().UTC()
	date := t.Format(time.RFC3339Nano)
	version := Version{Date: date}
	err := writeVersion(version)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func main() {
	releaseFlag := flag.Bool("release", false, "make a release")
	installFlag := flag.Bool("install", false, "install minio")

	flag.Parse()

	if *releaseFlag {
		runMinioRelease()
	}
	if *installFlag {
		runMinioInstall()
	}
}
