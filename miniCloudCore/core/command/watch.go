package command

import (
	"context"
	"io"
	"os"
	"regexp"

	log "github.com/Larouimohammed/miniCloud.git/infra/logger"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetInstance(ctx context.Context, cli *client.Client, cn string, log log.Log) (int32, error) {
	// we should fic that
	list, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Logger.Sugar().Error(" listing error : %v", err)
		return 0, err
	}
	instance := 0
	for _, l := range list {

		for _, n := range l.Names {
			result := n
			//we should fi that
			var reNameBlacklist = regexp.MustCompile(`(&|>|/|0|1|2|3|4|5|6|7|9|<|\/|:|\n|\r)*`)
			cnn := reNameBlacklist.ReplaceAllString(result, "")
			if cn == cnn {
				instance++
				out, err := cli.ContainerLogs(context.Background(), l.ID, types.ContainerLogsOptions{ShowStdout: true, Details: true})
				if err != nil {
					log.Logger.Sugar().Error(err.Error())

				}
				_, err = io.Copy(os.Stdout, out)

				if err != nil {
					log.Logger.Sugar().Error(err.Error())

				}

			}

		}

	}
	return int32(instance), nil

}
