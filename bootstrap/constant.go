package bootstrap

type Constants struct {
	Context Context
}

type Context struct {
	Translator string
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator: "translator",
		},
	}
}
