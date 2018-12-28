// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	goVersion   = "1.11"
	packageName = "github.com/brettpechiney/workout-service"
)

var (
	// Default target to run when none is passed specified.
	Default    = StartRoachContainer
	gocmd      = mg.GoCmd()
	releaseTag = regexp.MustCompile(`^v+[0-9]+\.[0-9]+\.[0-9]+$`)
)

// StartRoachContainer starts the CockroachDB container and kicks off it's
// migration scripts.
func StartRoachContainer() error {
	const MsgPrefix = "in StartRoachContainer"
	stopped := make(chan struct{})
	errchan := make(chan error)
	go func() {
		defer close(stopped)
		if err := sh.Run("cmd", "/C", "start", "docker-compose", "up", "cockroach"); err != nil {
			errchan <- fmt.Errorf("%s: %v", MsgPrefix, err)
		}
	}()
	for {
		select {
		case err := <-errchan:
			return err
		case <-stopped:
			return sh.Run("mage", "-d", "./migrations", "Migrate")
		}
	}
}

// Install runs go install and generates version information in the binary.
func Install() error {
	pkgslice := strings.Split(packageName, "/")
	name := pkgslice[len(pkgslice)-1]
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	gopath, err := sh.Output(gocmd, "env", "GOPATH")
	if err != nil {
		return fmt.Errorf("unable to detect GOPATH: %v", err)
	}
	paths := strings.Split(gopath, string([]rune{os.PathListSeparator}))
	bin := filepath.Join(paths[0], "bin")
	if err := os.Mkdir(bin, 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create %q: %v", bin, err)
	}
	path := filepath.Join(bin, name)
	return sh.Run(gocmd, "build", "-o", path, "-ldflags="+flags(), packageName)
}

// Test runs all tests with the verbose flag set.
func Test() error {
	if !isGoLatest() {
		return fmt.Errorf("go %s.x must be installed", goVersion)
	}
	s, err := sh.Output(gocmd, "test", "./...")
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

// TestRace runs all tests with both the race and verbose flags set.
func TestRace() error {
	s, err := sh.Output(gocmd, "test", "-v", "-race", "./...")
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}

// Fmt runs the gofmt tool on all packages in the project.
func Fmt() error {
	if err := sh.RunV(gocmd, "fmt", "./..."); err != nil {
		return fmt.Errorf("error running go fmt: %v", err)
	}
	return nil
}

// Lint runs the golint tool on all packages in the project.
func Lint() error {
	if err := sh.RunV("golint", "./..."); err != nil {
		return fmt.Errorf("error running golint: %v", err)
	}
	return nil
}

// Vet runs the vet tool with all flags set.
func Vet() error {
	if err := sh.RunV(gocmd, "vet", "-all", "./..."); err != nil {
		return fmt.Errorf("error running go vet: %v", err)
	}
	return nil
}

// Check makes sure the correct Golang version is being used, runs Fmt and Vet
// in parallel, and then runs TestRace.
func Check() error {
	if !isGoLatest() {
		return fmt.Errorf("go %s.x must be installed", goVersion)
	}
	mg.Deps(Fmt, Vet)
	mg.Deps(TestRace)
	return nil
}

// Release creates a new version of the service. It expects the TAG environment
// variable to be set because that is what it will use to create a new git tag.
// This target pushes changes to the Github repo.
func Release() (err error) {
	tag := strings.TrimSpace(os.Getenv("TAG"))
	if !releaseTag.MatchString(tag) {
		return fmt.Errorf("TAG should have format vx.x.x but is %s", tag)
	}

	trimmedTag := strings.TrimPrefix(tag, "v")
	msg := fmt.Sprintf("Release %s", trimmedTag)
	if err := sh.Run("git", "tag", "-a", tag, "-m", msg); err != nil {
		return fmt.Errorf("error adding git tag: %v", err)
	}
	if err := sh.Run("git", "push", "origin", tag); err != nil {
		return fmt.Errorf("error pushing to origin: %v", err)
	}
	defer func() {
		if err != nil {
			_ = sh.Run("git", "tag", "--delete", "$TAG")
			_ = sh.Run("git", "push", "--delete", "origin", "$TAG")
		}
	}()
	return nil
}

const (
	protocmd   = "protoc"
	protoSrc   = "--proto_path=pb/"
	goSrc      = "--proto_path=${GOPATH}/src"
	googleAPIs = "--proto_path=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis"
	compileOut = "--go_out=plugins=grpc:pb"
	proxyOut   = "--grpc-gateway_out=logtostderr=true:pb"
)

// ProtoCompile builds the service's '.proto' files.
func ProtoCompile() error {
	err := sh.RunV(
		protocmd,
		protoSrc,
		compileOut,
		"pb/common.proto",
		"pb/movementservice.proto",
	)
	if err != nil {
		return fmt.Errorf("error running protoc: %v", err)
	}
	return nil
}

// ProtoProxy generates a reverse proxy that translates a RESTful JSON API
// into gRPC calls.
func ProtoProxy() error {
	err := sh.RunV(
		protocmd,
		protoSrc,
		googleAPIs,
		proxyOut,
		"pb/common.proto",
		"pb/movementservice.proto",
	)
	if err != nil {
		return fmt.Errorf("error running protoc: %v", err)
	}
	return nil
}

func isGoLatest() bool {
	return strings.Contains(runtime.Version(), goVersion)
}

func flags() string {
	timestamp := time.Now().Format(time.RFC3339)
	hash := hash()
	tag := tag()
	if tag == "" {
		tag = "dev"
	}
	ts := fmt.Sprintf("%s.timestamp=%s", packageName, timestamp)
	ghash := fmt.Sprintf("%s.commitHash=%s", packageName, hash)
	rtag := fmt.Sprintf("%s.gitTag=%s", packageName, tag)
	return fmt.Sprintf(`-X %s -X %s -X %s `, ts, ghash, rtag)
}

func tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	return s
}

func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}
