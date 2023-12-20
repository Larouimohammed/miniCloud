package command

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	consul "github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/consulproxy"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var P *consul.ConsulProxy

func ProvApply(cli *client.Client, containername string, image string, subnet string, numberofistance int32) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf(" initialisation docker client error : %v", err)
		return err
	}
	defer cli.Close()

	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		log.Printf("image pull error : %v", err)
		return err
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)
	for i := 0; i < int(numberofistance); i++ {
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Hostname: containername + fmt.Sprint(i),
			Image:    image,
			Cmd:      []string{"sleep", "120"},
			Tty:      false,
		}, nil, nil, nil, containername+fmt.Sprint(i))
		if err != nil {
			log.Printf("create container failled : %v", err)
			return err

		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			log.Printf("start container failled: %v", err)
			return err

		}
		p := P.NewProxy()
		if err := p.Start(containername + fmt.Sprint(i)); err != nil {

			log.Fatal(err)

		}

	}
	return nil
}
