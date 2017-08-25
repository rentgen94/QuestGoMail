package mocks

import (
	"errors"
	"github.com/rentgen94/QuestGoMail/entities"
)

type ExistPlayerDAOMock struct{}

func (*ExistPlayerDAOMock) FindPlayer(player *entities.Player) (*entities.Player, error) {
	return &entities.Player{Id: 2, Login: "qqq", Password: "111"}, nil
}

func (*ExistPlayerDAOMock) CreatePlayer(player *entities.Player) error {
	return nil
}

func (*ExistPlayerDAOMock) Exist(player *entities.Player) (bool, error) {
	return true, nil
}

func (*ExistPlayerDAOMock) FindPlayerById(id int) (*entities.Player, error) {
	return &entities.Player{Id: 2, Login: "qqq", Password: "111"}, nil
}

func (*ExistPlayerDAOMock) SelectAllPlayers() ([]*entities.Player, error) {
	return nil, nil
}

type NotExistPlayerDAOMock struct{}

func (*NotExistPlayerDAOMock) FindPlayer(player *entities.Player) (*entities.Player, error) {
	return nil, errors.New("")
}

func (*NotExistPlayerDAOMock) CreatePlayer(player *entities.Player) error {
	return errors.New("\"Player register error.\"")
}

func (*NotExistPlayerDAOMock) Exist(player *entities.Player) (bool, error) {
	return false, nil
}

func (*NotExistPlayerDAOMock) FindPlayerById(id int) (*entities.Player, error) {
	return nil, errors.New("")
}

func (*NotExistPlayerDAOMock) SelectAllPlayers() ([]*entities.Player, error) {
	return nil, nil
}
