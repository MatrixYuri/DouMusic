package model

import "testing"

func testDownloadWorker(t *testing.T) {
	downloadWorker("离歌 (Live)"+".mp3", "http://mr3.douban.com/201505150937/5d9b0a86e5968617532a61f086c93100/view/song/small/p1476943.mp4")
}

func testDownload(t *testing.T) {

}
