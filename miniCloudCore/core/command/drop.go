package command

import (
	"context"
	"fmt"
	"log"

	consul "github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/consulproxy"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var c *consul.ConsulProxy

func StopandDropContainer(cli *client.Client, containername string, numberofistance int32) error {
	// c := P.NewProxy()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf(" initialisation docker client error : %v", err)
		return err
	}
	defer cli.Close()
	ctx := context.Background()
	for i := 0; i < int(numberofistance); i++ {
		if err := cli.ContainerStop(ctx, containername+fmt.Sprint(i), container.StopOptions{}); err != nil {
			log.Printf("Unable to stop container %s: %s", containername+fmt.Sprint(i), err)
		}

		removeOptions := types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}

		if err := cli.ContainerRemove(ctx, containername+fmt.Sprint(i), removeOptions); err != nil {
			log.Printf("Unable to remove container: %s", err)
			return err
		}
		// c.Deregisterservice(containername+fmt.Sprint(i))

	}
	return nil
}
