package memory

type InMemoryStateRepo struct {
	states map[string]bool
}

func NewInMemoryStateRepo() *InMemoryStateRepo {
	return &InMemoryStateRepo{
		states: make(map[string]bool),
	}
}

func (r *InMemoryStateRepo) AddState(state string) error {
	r.states[state] = true
	return nil
}

func (r *InMemoryStateRepo) Exists(state string) (bool, error) {
	_, ok := r.states[state]
	return ok, nil
}

func (r *InMemoryStateRepo) Remove(state string) error {
	delete(r.states, state)
	return nil
}
