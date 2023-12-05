package command

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func Watching(cli *client.Client, cn string) (int32, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf(" initialisation docker client error : %v", err)
		return 0, err
	}
	defer cli.Close()

	list, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Printf(" listing error : %v", err)
		return 0, err
	}

	instance := 0
	for i, l := range list {
		for _, n := range l.Names {
			fmt.Println(n)
			fmt.Println(cn + fmt.Sprint(i))
			if cn+fmt.Sprint(i) == string(n) {
				instance++

			}
		}

		// contains := slices.Contains((l.Names), cn+fmt.Sprint(i))
		// if contains {
		// 	instance++
		// }

	}
	return int32(instance), nil

}
