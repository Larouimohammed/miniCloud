package command

import (
	"context"
	"fmt"

	consul "github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/consulproxy"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var c *consul.ConsulProxy

func StopandDropContainer(cli *client.Client, containername string, numberofistance int32) error {
	ctx := context.Background()
	for i := 0; i < int(numberofistance); i++ {
		if err := cli.ContainerStop(ctx, containername+fmt.Sprint(i), container.StopOptions{}); err != nil {
			return err
		}

		removeOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}

		if err := cli.ContainerRemove(ctx, containername+fmt.Sprint(i), removeOptions); err != nil {
			return err
		}
		// c.Deregisterservice(containername+fmt.Sprint(i))

	}
	return nil
}
