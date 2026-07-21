package section

type (
	Processor struct {
		WebServer ProcessorWebServer `json:"web_server" split_words:"true" env:"WEB_SERVER"`
	}

	ProcessorWebServer struct {
		ListenPort uint32 `json:"listen_port" env:"LISTEN_PORT" default:"8080"`
	}
)
