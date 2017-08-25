package management

import (
	"errors"
	"fmt"
	"time"
	"sync"
)

const (
	managerNotFoundTemplate = "Game %d not found"
	failedOnTimeOut         = "Failed on time out"
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
	managerMutex sync.Mutex
	cntMutex sync.Mutex
	respMutex sync.Mutex
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
		managerMutex: sync.Mutex{},
		cntMutex: sync.Mutex{},
		respMutex: sync.Mutex{},
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

func (pool *ManagerPool) AddManager(manager *PlayerManager) int {
	var gameId = pool.cnt

	pool.cntMutex.Lock()
	pool.cnt++
	pool.cntMutex.Unlock()

	pool.respMutex.Lock()
	pool.respMap[gameId] = make(chan Response, pool.respBuffSize)
	pool.respMutex.Unlock()

	pool.managerMutex.Lock()
	pool.managers[gameId] = manager
	pool.managerMutex.Unlock()

	return gameId
}

func (pool *ManagerPool) DeleteManager(id int) {
	pool.managerMutex.Lock()
	delete(pool.managers, id)
	pool.managerMutex.Unlock()

	pool.respMutex.Lock()
	delete(pool.respMap, id)
	pool.respMutex.Unlock()
}

func (pool *ManagerPool) Running(id int) bool {
	var _, ok = pool.managers[id]
	if !ok {
		return false
	}

	return pool.managers[id].stateCode == managerWorkCode
}

func (pool *ManagerPool) SendCommand(command AddressedCommand) {
	pool.commandChan <- command
}

func (pool *ManagerPool) GetResponseSync(gameId int, timeout time.Duration) (Response, error) {
	var ch, ok = pool.respMap[gameId]
	if !ok {
		return Response{}, errors.New(fmt.Sprint(managerNotFoundTemplate, gameId))
	}

	for {
		select {
		case result := <-ch:
			return result, nil
		case <-time.After(timeout):
			return Response{}, errors.New(failedOnTimeOut)
		}
	}
	return <-ch, nil
}

func (pool *ManagerPool) stop() {
	pool.running = false
	pool.stopChan <- 1

	pool.managerMutex.Lock()
	for _, manager := range pool.managers {
		manager.Stop()
	}
	pool.managerMutex.Unlock()
}

func (pool *ManagerPool) monitorManagers() {
	for pool.running {

		pool.managerMutex.Lock()
		for key, manager := range pool.managers {
			if manager.Finished() {
				delete(pool.managers, key)
			}
		}
		pool.managerMutex.Unlock()

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
