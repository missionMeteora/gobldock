package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/mflag.v1"
)

var (
	packageName string
	outputPath  string
	goPath      string
	silent      bool
)

func init() {
	mflag.BoolVar(&silent, []string{"s", "-silent"}, false, "Be silent, no fancy stuff")
	mflag.StringVar(&outputPath, []string{"o", "-output"}, ".", "Specify the path for resulting exectuable")
	mflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <package>\n", os.Args[0])
		mflag.PrintDefaults()
	}
	mflag.Parse()
	packageName = mflag.Arg(0)
	if len(packageName) == 0 {
		mflag.Usage()
		os.Exit(1)
	}
	if goPath = os.Getenv("GOPATH"); len(goPath) == 0 {
		fatal("$GOPATH env variable is not defined")
	}
}

func main() {
	dockerPath, err := getDocker()
	if err != nil {
		fatal("error finding docker executable:", err)
	} else if !silent {
		fmt.Println("found Docker:", dockerPath)
	}

	ver, err := getDockerVersion(dockerPath)
	if err != nil {
		fatal("Docker cannot run:", err)
	} else if !silent {
		fmt.Println(ver)
	}

	dir, name, err := getDirAndFile(outputPath)
	if err != nil {
		fatal("error while checking output path:", err)
	} else if err := os.MkdirAll(dir, 0755); err != nil {
		fatal("error creating output dir:", err)
	}
	if len(name) == 0 {
		name = filepath.Base(packageName)
	}
	if !silent {
		fmt.Println("compiling into:", filepath.Join(dir, name))
	}

	if err := runDockerGoBuild(&buildSpec{
		dockerPath:  dockerPath,
		goPath:      goPath,
		outDir:      dir,
		execName:    name,
		packageName: packageName,
	}, os.Stdout, os.Stderr); err != nil {
		fatal("error running go build using Docker:", err)
	}
}

func getDirAndFile(path string) (dir string, name string, err error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", "", err
	} else if outputPath == "." {
		return abs, "", nil
	} else if info, err := os.Stat(name); err == nil && info.IsDir() {
		return abs, "", errors.New("the path should specify a file, not dir")
	}
	dir = filepath.Dir(abs)
	name = filepath.Base(abs)
	return
}

type buildSpec struct {
	dockerPath  string
	goPath      string
	outDir      string
	execName    string
	packageName string
}

func runDockerGoBuild(spec *buildSpec, stdout, stderr io.Writer) error {
	// docker run -it --rm -v $GOPATH:/go -v $OUT/bin:/data golang go build -o /data/$NAME -i $PACKAGE
	build := exec.Cmd{
		Path: spec.dockerPath,
		Args: []string{
			spec.dockerPath,
			"run", "-it", "--rm",
			"-v", fmt.Sprintf("%s:/go", spec.goPath),
			"-v", fmt.Sprintf("%s:/data", spec.outDir),
			"golang", "go", "build",
			"-o", filepath.Join("/data", spec.execName),
			"-i", spec.packageName,
		},
		Stdout: stdout,
		Stderr: stderr,
	}
	return build.Run()
}

func getDocker() (string, error) {
	out, err := exec.Command("which", "docker").Output()
	if err != nil {
		return "", errors.New("executable file not found in $PATH")
	}
	return string(bytes.TrimSpace(out)), nil
}

func getDockerVersion(path string) (string, error) {
	out, err := exec.Command(path, "-v").Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(out)), nil
}

func fatal(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}
