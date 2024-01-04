package main
import (
	"fmt"

	"log"

	ansible "github.com/febrianrendak/go-ansible"
)
const DefauttPathInventory ="./inventory"
type Executor interface {
	Execute(command string, args []string, prefix string) error
}
type MyExecutor struct{}

func (e *MyExecutor) Execute(command string, args []string, prefix string) error {
	fmt.Println("I am doing nothing")

	return nil
}
func RunAnsible(playbookName, inventory string) error {

	ansiblePlaybookConnectionOptions := &ansible.AnsiblePlaybookConnectionOptions{
		AskPass:    false,
		Connection: "local",
	}
	ansiblePlaybookOptions := &ansible.AnsiblePlaybookOptions{
		Inventory: DefauttPathInventory,
	}
	playbook := &ansible.AnsiblePlaybookCmd{
		Playbook:          inventory,
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
func main() {
	RunAnsible("my play", "172.17.04")
}
