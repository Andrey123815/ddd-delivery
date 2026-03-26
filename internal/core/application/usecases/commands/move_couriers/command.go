package moveCouriers

type MoveCouriersCommand struct {}

func NewMoveCouriersCommand() (*MoveCouriersCommand, error) {	
	return &MoveCouriersCommand{}, nil
}
