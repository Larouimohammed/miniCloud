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
	"github.com/docker/docker/client"
)

func ProvApply(cli *client.Client, containername string, image string, subnet string, numberofistance int32, log log.Log, consulproxy *consul.ConsulProxy) error {

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
			Cmd:             []string{"sleep", "1200"},
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

		go func(j int) {
			if err := consulproxy.Start(containername+fmt.Sprint(j), resp.ID, "172.17.0.4", 80); err != nil {

				log.Logger.Sugar().Error(err)

			}
		}(i)

		// kv := consulproxy.Cli.KV()
		// p := &capi.KVPair{Key: "TTL", Value: []byte("10")}
		// _, err = kv.Put(p, nil)
		// if err != nil {
		// 	panic(err)
		// }
	}
	return nil
}
