package provisioner

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/emicklei/go-restful/v3/log"
)

type Provisioner struct {
}
type Config struct {
	Containername string `json:"containername"`
	Image         string `json:"image"`
	Subnet        string `json:"subnet"`
	Nunofinstance int32  `json:"numofvms"`
}

func (P *Provisioner) ContainerProvisioner(config Config) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Logger.Printf(err.Error())
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, config.Image, types.ImagePullOptions{})
	if err != nil {
		log.Logger.Printf(err.Error())
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)
	for i := 0; i < int(config.Nunofinstance); i++ {
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Hostname: config.Containername + fmt.Sprint(i),
			Image:    config.Image,
			Cmd:      []string{"sleep", "10"},
			Tty:      false,
		}, nil, nil, nil, config.Containername+fmt.Sprint(i))
		if err != nil {
			log.Logger.Printf(err.Error())

		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			log.Logger.Printf(err.Error())

		}

		// statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		// select {
		// case err := <-errCh:
		// 	if err != nil {
		// 		log.Logger.Printf(err.Error())
		// 	}
		// case <-statusCh:
		// }

		// out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})

		// _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		// if err != nil {
		// 	log.Logger.Printf(err.Error())

		// }
	}
}


