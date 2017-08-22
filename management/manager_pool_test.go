package management

import (
	"fmt"
	"testing"
)

func TestManagerPool_Run(t *testing.T) { // smoke test :)
	var player1 = getPlayer()
	var player2 = getPlayer()
	var manager1, _ = NewPlayerManager(player1, 10, 10)
	var manager2, _ = NewPlayerManager(player2, 10, 10)
	var pool = NewManagerPool(10, 10)

	go pool.Run()
	defer pool.Stop()

	pool.AddManager(manager1)
	pool.AddManager(manager2)

	var command = AddressedCommand{
		Address: 0,
		Command: NewCommand(GetSlotsCode, "", nil, nil),
	}

	pool.SendCommand(command)
	var resp = pool.GetResponseSync(0)
	fmt.Println(resp)
}
