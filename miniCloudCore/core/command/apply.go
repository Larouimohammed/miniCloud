package command

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Larouimohammed/miniCloud.git/cli/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func ProvApply(cli client.Client, con config.Config) {
	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, con.Image, types.ImagePullOptions{})
	if err != nil {
		log.Printf(err.Error())
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)
	for i := 0; i < int(con.Nunofinstance); i++ {
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Hostname: con.Containername + fmt.Sprint(i),
			Image:    con.Image,
			Cmd:      []string{"sleep", "10"},
			Tty:      false,
		}, nil, nil, nil, con.Containername+fmt.Sprint(i))
		if err != nil {
			log.Printf(err.Error())

		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			log.Printf(err.Error())

		}

	}
}
