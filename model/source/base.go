package source

type Agent interface {
	Search(string) (Song, bool)
	GetDatil(*Song) bool
}

type Song struct {
	Name   string
	Artist string
	Album  string
	Url    string
}

var (
	Agents []Agent
)

func init() {
	Agents = []Agent{new(Netease)} //注册所有的Agent
}

func SearchAll(key string) (Song, bool) {
	for _, agent := range Agents {
		if song, hasResult := agent.Search(key); hasResult {
			return song, true
		}
	}
	return Song{}, false
}
