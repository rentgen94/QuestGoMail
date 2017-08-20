package management

import "fmt"

const (
	managerNotFoundTemplate = "Manager %d not found"
)

type AddressedCommand struct {
	Address int
	Command Command
}

type AddressedResponse struct {
	Address  int
	Response Response
}

type ManagerPool struct {
	cnt         int
	managers    map[int]*PlayerManager
	commandChan chan AddressedCommand
	respChan    chan AddressedResponse
	stopChan    chan interface{}
	running     bool
}

func NewManagerPool() *ManagerPool {
	return &ManagerPool{
		cnt:         0,
		managers:    make(map[int]*PlayerManager),
		commandChan: make(chan AddressedCommand),
		respChan:    make(chan AddressedResponse),
		stopChan:    make(chan interface{}),
		running:     false,
	}
}

func (pool *ManagerPool) Run() {
	pool.running = true
	for _, manager := range pool.managers {
		go manager.Run()
	}

	for {
		pool.handleCommandChan()
		select {
		case <-pool.stopChan:
			break
		default:
		}
	}
}

func (pool *ManagerPool) Stop() {
	if !pool.running {
		return
	}

	pool.running = false
	pool.stopChan <- 1

	for _, manager := range pool.managers {
		manager.Stop()
	}
}

func (pool *ManagerPool) AddManager(manager *PlayerManager) {
	var id = pool.cnt
	pool.cnt++
	pool.managers[id] = manager
}

func (pool *ManagerPool) DeleteManager(id int) {
	delete(pool.managers, id)
}

func (pool *ManagerPool) SendCommand(command AddressedCommand) {
	pool.commandChan <- command
}

func (pool *ManagerPool) GetResponseSync(gameId int) Response {
	var manager, ok = pool.managers[gameId]
	if !ok {
		return Response{
			errMsg: fmt.Sprintf(managerNotFoundTemplate, gameId),
		}
	}

	return <-manager.outChan
}

func (pool *ManagerPool) handleCommandChan() {
	select {
	case command := <-pool.commandChan:
		var manager, ok = pool.managers[command.Address]
		if !ok {
			pool.respChan <- AddressedResponse{
				Address: command.Address,
				Response: Response{
					errMsg: fmt.Sprintf(managerNotFoundTemplate, command.Address),
				},
			}
		}
		manager.CommandChan() <- command.Command
	default:
		return
	}
}
