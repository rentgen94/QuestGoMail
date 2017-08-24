package management

import (
	"fmt"
	"errors"
)

const (
	managerNotFoundTemplate = "Game %d not found"
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
	cnt          int
	workerNum    int
	managers     map[int]*PlayerManager
	commandChan  chan AddressedCommand
	respMap      map[int]chan Response
	stopChan     chan interface{}
	running      bool
	respBuffSize int
}

func NewManagerPool(workerNum int, commandBuffSize int, respBuffSize int) *ManagerPool {
	return &ManagerPool{
		cnt:          0,
		workerNum:    workerNum,
		managers:     make(map[int]*PlayerManager),
		commandChan:  make(chan AddressedCommand, commandBuffSize),
		respMap:      make(map[int]chan Response),
		stopChan:     make(chan interface{}, 1),
		running:      false,
		respBuffSize: respBuffSize,
	}
}

func (pool *ManagerPool) Run() {
	pool.running = true
	for _, manager := range pool.managers {
		go manager.Run()
	}
	go pool.monitorManagers()

	for i := 0; i != pool.workerNum; i++ {
		go pool.handleCommandChanLoop()
	}

	for {
		select {
		case <-pool.stopChan:
			break
		default:
		}
	}

	pool.stop()
}

func (pool *ManagerPool) Stop() {
	if !pool.running {
		return
	}

	pool.stop()
}

func (pool *ManagerPool) AddManager(manager *PlayerManager,gameid int) {
	//var id = pool.cnt
	//pool.cnt++
	pool.managers[gameid] = manager
}

func (pool *ManagerPool) DeleteManager(id int) {
	delete(pool.managers, id)
	delete(pool.respMap, id)
}

func (pool *ManagerPool) SendCommand(command AddressedCommand) {
	pool.commandChan <- command
}

func (pool *ManagerPool) GetResponseSync(gameId int) (Response, error) {
	var ch, ok = pool.respMap[gameId]
	if !ok {
		return Response{}, errors.New (fmt.Sprint(managerNotFoundTemplate, gameId))
	}

	return <-ch, nil
}

func (pool *ManagerPool) stop() {
	pool.running = false
	pool.stopChan <- 1

	for _, manager := range pool.managers {
		manager.Stop()
	}
}

func (pool *ManagerPool) monitorManagers() {
	for pool.running {
		for key, manager := range pool.managers {
			if manager.Finished() {
				delete(pool.managers, key)
			}
		}
	}
}

func (pool *ManagerPool) handleCommandChanLoop() {
	for pool.running {
		pool.handleCommandChan()
	}
}

func (pool *ManagerPool) handleCommandChan() {
	select {
	case command := <-pool.commandChan:
		var manager, ok = pool.managers[command.Address]
		if !ok {
			pool.respMap[command.Address] <- Response{
				ErrMsg: fmt.Sprintf(managerNotFoundTemplate, command.Address),
			}
		}
		manager.CommandChan() <- command.Command
		pool.respMap[command.Address] <- <-manager.RespChan()
	default:
		return
	}
}
