package memory

type StateRepo struct {
	states map[string]bool
}

func NewStateRepo() *StateRepo {
	return &StateRepo{
		states: make(map[string]bool),
	}
}

func (r *StateRepo) AddState(state string) error {
	r.states[state] = true
	return nil
}

func (r *StateRepo) Exists(state string) (bool, error) {
	_, ok := r.states[state]
	return ok, nil
}

func (r *StateRepo) Remove(state string) error {
	delete(r.states, state)
	return nil
}
