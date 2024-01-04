package command

import (
	"context"
	"fmt"
	"io"
	"os"
	log "github.com/Larouimohammed/miniCloud.git/logger"
	consul "github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/consulproxy"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	cli "github.com/docker/docker/client"
)

func ProvApply(cli *cli.Client, containername, image, subnet string, numberofistance int32, command, installWithAnsible []string, log log.Log, consulproxy *consul.ConsulProxy) error {
	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {

		return err
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)
	for i := 0; i < int(numberofistance); i++ {
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Hostname:        containername + fmt.Sprint(i),
			Image:           image,
			Cmd:             command,
			Tty:             false,
			NetworkDisabled: false,
		}, &container.HostConfig{PublishAllPorts: true, Privileged: false}, nil, nil, containername+fmt.Sprint(i))
		if err != nil {

			return err

		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			log.Logger.Sugar().Error("start container failled: %v", err)
			return err

		}
		// ruu ansible install
		// ansible.RunAnsible("myplay", subnet)

		// consul service register
		go func(j int) {
			if err := consulproxy.Start(containername+fmt.Sprint(j), resp.ID, "172.17.0.4", 80); err != nil {

				log.Logger.Sugar().Error(err)

			}
		}(i)

	}
	return nil
}
