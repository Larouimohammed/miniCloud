package ansible

import (
	"fmt"
	"strings"

	log "github.com/Larouimohammed/miniCloud.git/infra/logger"

	ansible "github.com/febrianrendak/go-ansible"
)

type Executor interface {
	Execute(command string, args []string, prefix string) error
}
type MyExecutor struct{}

func (e *MyExecutor) Execute(command string, args []string, prefix string) error {
	fmt.Println("I am doing nothing")

	return nil
}
func RunAnsible(AnsiblePlaybookPath string, inventory []string, log log.Log) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{
		AskPass:    false,
		Connection: "local",
	}

	ansiblePlaybookPrivilegeEscalationOptions := &ansible.AnsiblePlaybookPrivilegeEscalationOptions{
		// Become:        true,
		// AskBecomePass: true,
	}
	strInventory := strings.Join(inventory, ",")
	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		Inventory: strInventory,
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
