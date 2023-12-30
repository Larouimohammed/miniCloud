package main

import (
	"fmt"

	"log"

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

func main() {
	RunAnsible("dqd", "qdq")
}
func RunAnsible(playbookName, Inventory string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{
		Connection: "local",
	}
	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		Inventory: "172.17.0.4,172.17.0.5,172.17.0.6",
	}
	playbook := &ansible.AnsiblePlaybookCmd{
		Playbook:          playbookName,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		Exec:              &MyExecutor{},
	}
    playbook.Command()
	err := playbook.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
