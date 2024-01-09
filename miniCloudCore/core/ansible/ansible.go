package ansible

import (
	"fmt"

	log "github.com/Larouimohammed/miniCloud.git/logger"

	ansible "github.com/febrianrendak/go-ansible"
)

const DefauttPathInventory = "./inventory"

type Executor interface {
	Execute(command string, args []string, prefix string) error
}
type MyExecutor struct{}

func (e *MyExecutor) Execute(command string, args []string, prefix string) error {
	fmt.Println("I am doing nothing")

	return nil
}
func RunAnsible(AnsiblePlaybookPath, inventory string, log log.Log) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{
		AskPass: false,
		 Connection: "local",
	}

	ansiblePlaybookPrivilegeEscalationOptions := &ansible.AnsiblePlaybookPrivilegeEscalationOptions{
		// Become:        true,
		// AskBecomePass: true,
	}
	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		Inventory: 		"172.17.0.4,172.17.0.5",

		// "192.168.1.103,",
	}
	playbook := &ansible.AnsiblePlaybookCmd{
		Playbook:                   AnsiblePlaybookPath,
		ConnectionOptions:          ansiblePlaybookConnectionOptions,
		PrivilegeEscalationOptions: ansiblePlaybookPrivilegeEscalationOptions,
		ExecPrefix:                 "Go-ansible example",
		Options:                    ansiblePlaybookOptions,

		// Exec: &MyExecutor{},
	}

	ansiblelog, ansibleError := playbook.Command()
	log.Logger.Sugar().Info(ansiblelog)
	if ansibleError != nil {
		log.Logger.Sugar().Error(ansibleError)
		return ansibleError
	}
	err := playbook.Run()
	if err != nil {
		log.Logger.Sugar().Error(err)
		return err
	}
	return nil
}

//  func main() {
// 	RunAnsible("/home/ubuntu/miniCloud/miniCloud/miniCloudCore/core/ansible/ansible.yml", "172.17.04",*log.Newlogger())
//  }
