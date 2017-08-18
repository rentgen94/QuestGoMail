package management

import (
	"testing"
	"fmt"
)

func TestManagerPool_Run(t *testing.T) {	// smoke test :)
	var player1 = getPlayer()
	var player2 = getPlayer()
	var manager1, _ = NewPlayerManager(player1)
	var manager2, _ = NewPlayerManager(player2)
	var pool = NewManagerPool()

	go pool.Run()
	defer pool.Stop()

	pool.AddManager(manager1)
	pool.AddManager(manager2)

	var command = AddressedCommand{
		Address:0,
		Command:NewCommand(getSlotsCode, "", nil, nil),
	}

	pool.SendCommand(command)
	var resp = pool.ReceiveBlock()
	fmt.Println(resp)
}