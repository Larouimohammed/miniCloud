package command

import (
	"context"
	"fmt"
	"io"
	"os"

	log "github.com/Larouimohammed/miniCloud.git/infra/logger"
	consul "github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/plugin/consulproxy"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	cli "github.com/docker/docker/client"
)

func ProvApply(cli *cli.Client, containername, image, subnet, installWithAnsible string, numberofistance int32, command []string, log log.Log, consulproxy *consul.ConsulProxy) error {
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
			AttachStdin:     true,
			AttachStdout:    true,
			AttachStderr:    true,
			Healthcheck:     &container.HealthConfig{Test: consulproxy.Cli.Headers()[containername+fmt.Sprint(i)]},
		}, &container.HostConfig{PublishAllPorts: true, Privileged: false}, nil, nil, containername+fmt.Sprint(i))
		if err != nil {
			log.Logger.Sugar().Error("create  container failled: %v", err)

			return err

		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			log.Logger.Sugar().Error("start container failled: %v", err)
			return err

		}
		// ruu ansible install

		// ansible.RunAnsible(installWithAnsible, "", log)

		// consul service register
		go func(j int) {
			if err := consulproxy.Start(containername+fmt.Sprint(j), resp.ID, "172.17.0.4", 80); err != nil {

				log.Logger.Sugar().Error("Service registred failled", err)

			}
		}(i)

	}
	// run ansible
	return nil
}
