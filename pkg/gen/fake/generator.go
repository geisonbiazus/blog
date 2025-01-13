package fake

type Generator struct {
	ReturnID string
}

func NewGenerator() *Generator {
	return &Generator{ReturnID: "FAKE_ID"}
}

func (g *Generator) Generate() string {
	return g.ReturnID
}
