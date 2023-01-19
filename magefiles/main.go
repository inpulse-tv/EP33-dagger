package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func Test() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}

	workdir := client.
		Host().
		Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"magefiles/", "go.work"},
		})

	_, err = client.
		Container().
		From("golang:alpine").
		WithDirectory("/src", workdir).
		WithWorkdir("/src").
		WithExec([]string{"go", "test", "./math"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}

func Build() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}

	workdir := client.
		Host().
		Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"magefiles/", "go.work"},
		})

	_, err = client.
		Container().
		From("golang:alpine").
		WithDirectory("/src", workdir).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "dagger-demo"}).
		ExitCode(ctx)

	if err != nil {
		return err
	}
	return nil
}

func Run() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}

	workdir := client.
		Host().
		Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"magefiles/", "go.work"},
		})

	bin := client.
		Container().
		From("golang:alpine").
		WithDirectory("/src", workdir).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "dagger-demo"}).
		File("dagger-demo")

	output, err := client.
		Container().
		From("alpine").
		WithFile(".", bin).
		WithExec([]string{"./dagger-demo", "2", "2"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	fmt.Println(output)

	return nil
}
