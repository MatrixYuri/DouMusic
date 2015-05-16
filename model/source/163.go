package source

const (
	songSearchApi_163 = "http://music.163.com/api/search/get/web"
	songDetailApi_163 = "http://music.163.com/api/song/detail"
)

type Netease struct {
	_ Song
}

func (s *Netease) Search(key string) (Song, bool) {
	return Song{}, true
}

func (s *Netease) GetDatil(song *Song) bool {
	return true
}
