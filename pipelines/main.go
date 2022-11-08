package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

func main() {
	if err := test(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func test(ctx context.Context) error {
	fmt.Println("Testing with Dagger")

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Get the current working directroy
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	srcPath := filepath.Join(pwd, "../myapp")

	// get reference to the mix project (pwd ../myapp)
	src := client.Host().Directory(srcPath)

	// Get an Elixir image
	container := client.Container().From("elixir:1.14.1-alpine")

	// mount cloned repository into `elixir` image
	container = container.WithMountedDirectory("/src", src).WithWorkdir("/src")
	container = container.Exec(dagger.ContainerExecOpts{
		Args: []string{"mix", "compile"},
	})

	// get reference to build output directory in container
	output := container.Directory("/src")

	// write contents of container build/ directory to the host
	_, err = output.Export(ctx, srcPath)
	if err != nil {
		return err
	}

	return nil
}
