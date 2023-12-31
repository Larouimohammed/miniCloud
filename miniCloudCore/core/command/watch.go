package command

import (
	"context"
	"io"
	"os"
	"regexp"

	log "github.com/Larouimohammed/miniCloud.git/logger"

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

				_, err = io.Copy(os.Stdout, out)

				if err != nil {
					log.Logger.Sugar().Error(err.Error())

				}

			}

		}
		// 	statusCh, errCh := cli.ContainerWait(context.Background(), l.ID, container.WaitConditionNotRunning)
		// 	select {
		// 	case err := <-errCh:
		// 		if err != nil {
		// 			log.Logger.Sugar().Error(err.Error())
		// 		}
		// 	case <-statusCh:
		// 	}

		// 	out, err := cli.ContainerLogs(context.Background(), l.ID, types.ContainerLogsOptions{ShowStdout: true})

		// 	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		// 	if err != nil {
		// 		log.Logger.Sugar().Error(err.Error())
		// 	}

		// }

	}
	return int32(instance), nil

}

// func Watching(cli *client.Client, containername string, log log.Log) {
// wg := sync.WaitGroup{}
// wg.Add(1)
// go func() {
// 	for {
// 		msgs, serrs := cli.Events(, types.EventsOptions{})

// 		select {
// 		case msg := <-msgs:
// 			log.Logger.Sugar().Info(msg)

// 			err := stream.Send(&pb.WResp{Wresp: fmt.Sprintf("%+v", msg), Werr: ""})
// 			if err != nil {
// 				S.logger.Logger.Sugar().Error(err)
// 				return err
// 			}
// 		case errs := <-serrs:
// 			S.logger.Logger.Sugar().Error(errs)

// 			err := stream.Send(&pb.WResp{Wresp: "", Werr: errs.Error()})
// 			if err != nil {
// 				S.logger.Logger.Sugar().Error(err)
// 				return err
// 			}

// 		}

// 	}
// }()
// wg.Wait()
// }
