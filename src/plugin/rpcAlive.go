package plugin

type Alive int64

func (a *Alive) RegisterAction(args, reply *int64) error {
	*reply = 1
	return nil
}
