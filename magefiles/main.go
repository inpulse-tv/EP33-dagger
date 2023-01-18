package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"dagger.io/dagger"
	"github.com/sourcegraph/conc/pool"
	"go.uber.org/multierr"
)

var platforms = []dagger.Platform{
	"linux/amd64", // a.k.a. x86_64
	"linux/arm64", // a.k.a. aarch64
	// "linux/s390x", // a.k.a. IBM S/390
}

func Engine() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}

func runUnitTest(client *dagger.Client, workdir *dagger.Directory, ctx context.Context, platform dagger.Platform) error {
	_, err := client.
		Container(dagger.ContainerOpts{Platform: platform}).
		From("golang:1.19.3-alpine").
		WithDirectory("/go/src/dagger-demo", workdir).
		WithWorkdir("/go/src/dagger-demo").
		WithExec([]string{"/usr/local/go/bin/go", "test", "./math", "-v"}).
		Stdout(ctx)
	return err
}

func Testconcurrent() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	p := pool.NewWithResults[error]()

	workdir := client.
		Host().
		Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"magefiles", "go.work"},
		})

	start := time.Now()
	for _, platform := range platforms {
		p.Go(func() error {
			return runUnitTest(client, workdir, ctx, platform)
		})
	}
	res := p.Wait()
	fmt.Printf("%d\n", time.Since(start).Milliseconds())
	var errors error
	for _, e := range res {
		if e != nil {
			return multierr.Append(errors, e)
		}
	}

	return errors
}

func Test() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	workdir := client.
		Host().
		Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"magefiles", "go.work"},
		})

	var errors error
	start := time.Now()
	for _, platform := range platforms {
		err := runUnitTest(client, workdir, ctx, platform)
		if err != nil {
			return multierr.Append(errors, err)
		}
	}
	fmt.Printf("%d\n", time.Since(start).Milliseconds())
	return errors
}

func innerBuild(client *dagger.Client) []*dagger.Container {
	workdir := client.Host().
		Directory(".", dagger.HostDirectoryOpts{Exclude: []string{"magefiles"}})
	containers := []*dagger.Container{}
	for _, platform := range platforms {
		containers = append(containers, client.
			Container(dagger.ContainerOpts{Platform: platform}).
			From("golang:1.19.3-alpine").
			WithDirectory("/go/src/dagger-demo", workdir).
			WithWorkdir("/go/src/dagger-demo").
			WithEnvVariable("CGO_ENABLED", "0").
			WithExec([]string{"/usr/local/go/bin/go", "build", "-o", "./output/dagger-demo"}))
	}
	return containers
}

func Build() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}

	for _, container := range innerBuild(client) {
		container.Stdout(ctx)
	}
	if err != nil {
		return err
	}

	return nil
}

func Export() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}

	platformVariants := []*dagger.Container{}
	for i, container := range innerBuild(client) {
		output := container.Directory("./output/")
		platformVariants = append(platformVariants, client.Container(dagger.ContainerOpts{Platform: platforms[i]}).WithRootfs(output))
	}

	fmt.Println(len(platformVariants))
	_, err = client.
		Container().
		Export(ctx, "dagger-demo.tar",
			dagger.ContainerExportOpts{PlatformVariants: platformVariants},
		)

	if err != nil {
		return err
	}

	return nil
}
