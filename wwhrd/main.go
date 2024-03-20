package main

import (
	"context"
	"fmt"

	"wwhrd/internal/dagger"
)

func New(
	// +optional
	// +default="latest"
	wwhrdVersion string,

	// +optional
	// +default="1.22"
	goVersion string,
) *WWHRD {
	return &WWHRD{
		WWHRDVersion: wwhrdVersion,
		GoVersion:    goVersion,
	}
}

type WWHRD struct {
	WWHRDVersion string
	GoVersion    string

	baseContainer  *dagger.Container
	wwhrdContainer *dagger.Container
}

// Runs the check command of wwhrd
//
// Example:
// $ dagger call --go-version=1.22 --wwhrd-version=latest check --source-dir=$HOME/my-go-project
func (w *WWHRD) Check(ctx context.Context, sourceDir *dagger.Directory) (string, error) {
	w.withWWHRDContainer()

	return w.wwhrdContainer.
		WithDirectory("/src", sourceDir).
		WithWorkdir("/src").
		WithExec([]string{"wwhrd", "check"}).
		Stdout(ctx)
}

func (w *WWHRD) withWWHRDContainer() *WWHRD {
	w.withBaseContainer()

	if w.wwhrdContainer != nil {
		return w
	}

	w.wwhrdContainer = w.baseContainer.
		WithExec([]string{"go", "install", fmt.Sprintf("github.com/frapposelli/wwhrd@%s", w.WWHRDVersion)})

	return w
}

func (w *WWHRD) withBaseContainer() *WWHRD {
	if w.baseContainer != nil {
		return w
	}

	w.baseContainer = dag.
		Container().
		From(fmt.Sprintf("golang:%s", w.GoVersion))

	return w
}
