package memory

type InMemoryStateRepo struct {
	states map[string]bool
}

func NewInMemoryStateRepo() *InMemoryStateRepo {
	return &InMemoryStateRepo{
		states: make(map[string]bool),
	}
}

func (r *InMemoryStateRepo) AddState(state string) {
	r.states[state] = true
}

func (r *InMemoryStateRepo) StateExists(state string) bool {
	_, ok := r.states[state]
	return ok
}
