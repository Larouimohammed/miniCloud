package command

import (
	"context"
	"fmt"
	"log"
	"regexp"

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
	for _, l := range list {
		for _, n := range l.Names {
			result := n
			var reNameBlacklist = regexp.MustCompile(`(&|>|/|0|1|2|3|4|5|6|7|9|<|\/|:|\n|\r)*`)
			cnn := reNameBlacklist.ReplaceAllString(result, "")
			fmt.Println(cnn)
			fmt.Println(cn)
			if cn == cnn {
				instance++

			}
		}

		}
	return int32(instance), nil

}
