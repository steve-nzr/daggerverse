// A generated module for GCC functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/ci/internal/dagger"
	"fmt"
)

func New(
	// +optional
	// +default="bookworm"
	debianVersion string,

	// +optional
	// +default="12"
	gccVersion string,

	// Your C++ project directory. It must contain a main.cpp
	sourceDirectory *dagger.Directory,
) *GCC {
	return &GCC{
		DebianVersion:   debianVersion,
		GccVersion:      gccVersion,
		SourceDirectory: sourceDirectory,
	}
}

type GCC struct {
	DebianVersion   string
	GccVersion      string
	SourceDirectory *dagger.Directory

	buildContainer   *dagger.Container
	compiledProgram  *dagger.File
	programContainer *dagger.Container
}

// Run your program in a container
//
// Example:
// $ dagger call --source-directory=$HOME/dev/cpp/src run
func (m *GCC) Run(ctx context.Context) (string, error) {
	m.withProgramContainer()

	return m.programContainer.
		WithExec([]string{"/main"}).
		Stdout(ctx)
}

// Creates a container from your program
//
// This will compile and put your application inside a distroless container
// that you can publish or run
//
// Example:
// $ dagger call --source-directory=$HOME/dev/cpp/src container publish --address="$addr"
func (m *GCC) Container() *dagger.Container {
	m.withProgramContainer()

	return m.programContainer.
		WithEntrypoint([]string{"/main"})
}

func (m *GCC) withBuildContainer() *GCC {
	if m.buildContainer != nil {
		return m
	}

	m.buildContainer = dag.
		Container().
		From(fmt.Sprintf("gcc:%s-%s", m.GccVersion, m.DebianVersion))

	return m
}

func (m *GCC) withCompiledProgram() *GCC {
	if m.compiledProgram != nil {
		return m
	}

	m.withBuildContainer()

	m.compiledProgram = m.buildContainer.
		WithDirectory("/src", m.SourceDirectory).
		WithExec([]string{"g++", "/src/main.cpp"}).
		File("./a.out")

	return m
}

func (m *GCC) withProgramContainer() *GCC {
	if m.programContainer != nil {
		return m
	}

	m.withCompiledProgram()

	m.programContainer = dag.
		Container().
		From("gcr.io/distroless/cc-debian12").
		WithFile("/main", m.compiledProgram)

	return m
}
