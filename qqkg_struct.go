package main

type PersonHomepage struct {
	Code    int    `json:"code"`
	Subcode int    `json:"subcode"`
	Message string `json:"message"`
	Default int    `json:"default"`
	Data    struct {
		Age        int    `json:"age"`
		AlbumNum   int    `json:"album_num"`
		AuthValue  int    `json:"auth_value"`
		FollowFlag int    `json:"follow_flag"`
		Follower   int    `json:"follower"`
		Following  int    `json:"following"`
		Friends    int    `json:"friends"`
		Gender     int    `json:"gender"`
		GreenLevel int    `json:"green_level"`
		GroupID    int    `json:"group_id"`
		GroupName  string `json:"group_name"`
		HasMore    int    `json:"has_more"`
		HeadImgURL string `json:"head_img_url"`
		IsAnchor   int    `json:"is_anchor"`
		IsGuest    int    `json:"is_guest"`
		IsVipGreen int    `json:"is_vip_green"`
		Isvip      int    `json:"isvip"`
		KgeUID     string `json:"kge_uid"`
		Kid        string `json:"kid"`
		Level      int    `json:"level"`
		Levelname  string `json:"levelname"`
		LiveInfo   struct {
			IDiamondLevel int `json:"iDiamondLevel"`
			IIsAnchor     int `json:"iIsAnchor"`
			IPVNum        int `json:"iPVNum"`
			IRelationID   int `json:"iRelationId"`
			ISelfStatus   int `json:"iSelfStatus"`
			IStatus       int `json:"iStatus"`
			IUsePVNum     int `json:"iUsePVNum"`
			MapExt        struct {
			} `json:"mapExt"`
			StH265Param struct {
				IEnableTransform int `json:"iEnableTransform"`
				ITransformType   int `json:"iTransformType"`
			} `json:"stH265Param"`
			StrAVAudienceRole      string `json:"strAVAudienceRole"`
			StrAnchorMuid          string `json:"strAnchorMuid"`
			StrDiamondLevelName    string `json:"strDiamondLevelName"`
			StrDiamondNum          string `json:"strDiamondNum"`
			StrGroupID             string `json:"strGroupId"`
			StrGroupType           string `json:"strGroupType"`
			StrLiveCoverURL        string `json:"strLiveCoverUrl"`
			StrPopularity          string `json:"strPopularity"`
			StrPushStreamLivingURL string `json:"strPushStreamLivingUrl"`
			StrRoomID              string `json:"strRoomID"`
			UOnlineNum             int    `json:"uOnlineNum"`
			UShowStartTime         int    `json:"uShowStartTime"`
		} `json:"live_info"`
		LoginAuthValue  int    `json:"login_auth_value"`
		LoginHeadImg    string `json:"login_head_img"`
		LoginNickname   string `json:"login_nickname"`
		LoginUID        string `json:"login_uid"`
		Nickname        string `json:"nickname"`
		SCityID         string `json:"sCityId"`
		SCountryID      string `json:"sCountryId"`
		SDistrictID     string `json:"sDistrictId"`
		SProvinceID     string `json:"sProvinceId"`
		SignedAnchor    int    `json:"signed_anchor"`
		TrackTotalCount int    `json:"track_total_count"`
		TreasureLevel   int    `json:"treasure_level"`
		TreasureValue   int    `json:"treasure_value"`
		UgcTotalCount   int    `json:"ugc_total_count"`
		Ugclist         []Ugc  `json:"ugclist"`
	} `json:"data"`
}

type Ugc struct {
	Albumid     int    `json:"albumid"`
	Avatar      string `json:"avatar"`
	ComentCount int    `json:"coment_count"`
	GiftCount   int    `json:"gift_count"`
	IsSegment   bool   `json:"is_segment"`
	KsongMid    string `json:"ksong_mid"`
	PlayCount   int    `json:"play_count"`
	ScoreRank   int    `json:"score_rank"`
	Shareid     string `json:"shareid"`
	Stralbumid  string `json:"stralbumid"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Ugcmask     int    `json:"ugcmask"`
}
