package fake

type IDGenerator struct {
	ReturnID string
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{ReturnID: "FAKE_ID"}
}

func (g *IDGenerator) Generate() string {
	return g.ReturnID
}
