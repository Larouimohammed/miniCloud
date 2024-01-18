package command

import (
	"context"

	log "github.com/Larouimohammed/miniCloud.git/infra/logger"
	consul "github.com/Larouimohammed/miniCloud.git/miniCloudCore/core/plugin/consulproxy"
	cli "github.com/docker/docker/client"
)

func Update(cli *cli.Client, containername, image, subnet, installWithAnsiblePath string, numberofistance int32, command []string, log log.Log, consulproxy *consul.ConsulProxy) error {
	instance, err := GetInstance(context.Background(), cli, containername, log)
	if err != nil {
		log.Logger.Sugar().Error("number of instance is indectectible  %v", err)
		return err
	}
	if err := StopandDropContainer(cli, containername, instance); err != nil {
		log.Logger.Sugar().Error(" droping infra  error  %v", err)
		return err
	}
	if err := ProvApply(cli, containername, image, subnet, installWithAnsiblePath, numberofistance, command, log, consulproxy); err != nil {
		log.Logger.Sugar().Error(" provisionning error  %v", err)
		return err

	}
	return nil
}
