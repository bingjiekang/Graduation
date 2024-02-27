package response

type MallUserDetailResponse struct {
	NickName      string `json:"nickName"`
	LoginName     string `json:"loginName"`
	UUid          int64  `json:"uUid"`
	IntroduceSign string `json:"introduceSign"`
}
