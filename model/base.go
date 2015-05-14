package model

const (
	loginApi      = "http://douban.fm/j/login"
	captchaApi    = "http://douban.fm/j/new_captcha"
	captchaImgApi = "http://douban.fm/misc/captcha?size=s&id=%s" //size=s小号m中号l大号

	likeListApi = "http://douban.fm/j/play_record?ck=%s&spbid=::%s&type=liked&start=%d"
)
