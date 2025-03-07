package bootstrap

type Constants struct {
	Context  Context
	LogLevel LogLevel
}

type Context struct {
	Translator string
}

type LogLevel struct {
	Debug string
	Info  string
	Warn  string
	Error string
	Fatal string
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator: "translator",
		},
		LogLevel: LogLevel{
			Debug: "debug",
			Info:  "info",
			Warn:  "warn",
			Error: "error",
			Fatal: "fatal",
		},
	}
}
