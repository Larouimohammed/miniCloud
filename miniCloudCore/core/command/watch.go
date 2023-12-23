package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func Watching(cli *client.Client, cn string) (int32, error) {
	// we should fic that
	list, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Printf(" listing error : %v", err)
		return 0, err
	}
	instance := 0
	for _, l := range list {
		for _, n := range l.Names {
			result := n
			//we should fi that
			var reNameBlacklist = regexp.MustCompile(`(&|>|/|0|1|2|3|4|5|6|7|9|<|\/|:|\n|\r)*`)
			cnn := reNameBlacklist.ReplaceAllString(result, "")
			fmt.Println(cnn)
			fmt.Println(cn)
			if cn == cnn {
				instance++

			}

		}

		statusCh, errCh := cli.ContainerWait(context.Background(), l.ID,container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				log.Printf(err.Error())
			}
		case <-statusCh:
		}

		out, err := cli.ContainerLogs(context.Background(), l.ID, types.ContainerLogsOptions{ShowStdout: true})

		_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		if err != nil {
			log.Printf(err.Error())
		}

	}

	return int32(instance), nil
}
