package management

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestManagerPool_Run_Success(t *testing.T) {
	var player1 = getPlayer()
	var player2 = getPlayer()
	var manager1, _ = NewPlayerManager(player1, 10, 10)
	var manager2, _ = NewPlayerManager(player2, 10, 10)
	var pool = NewManagerPool(3, 10, 10)

	go pool.Run()
	defer pool.Stop()

	pool.AddManager(manager1, 1)
	pool.AddManager(manager2, 2)

	var command = AddressedCommand{
		Address: 1,
		Command: NewCommand(GetSlotsCode, 0, nil, nil),
	}

	pool.SendCommand(command)
	var _, respErr = pool.GetResponseSync(1, 1*time.Second)
	assert.Nil(t, respErr)
}

func TestManagerPool_Run_TriedWrongGame(t *testing.T) {
	var player1 = getPlayer()
	var player2 = getPlayer()
	var manager1, _ = NewPlayerManager(player1, 10, 10)
	var manager2, _ = NewPlayerManager(player2, 10, 10)
	var pool = NewManagerPool(3, 10, 10)

	go pool.Run()
	defer pool.Stop()

	pool.AddManager(manager1, 1)
	pool.AddManager(manager2, 2)

	var command = AddressedCommand{
		Address: 0,
		Command: NewCommand(GetSlotsCode, 0, nil, nil),
	}

	pool.SendCommand(command)
	var _, respErr = pool.GetResponseSync(1, 1*time.Second)
	assert.Error(t, respErr, failedOnTimeOut)
}
