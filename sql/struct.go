package sql

type Danmaku struct {
	DanmakuId    int    `db:"danmaku_id" json:"danmaku_id"`
	SendMode     int    `db:"send_mode" json:"send_mode"`
	SendFontSize int    `db:"send_font_size" json:"send_font_size"`
	DanmakuColor int64  `db:"danmaku_color" json:"danmaku_color"`
	Time         int64  `db:"time" json:"time"`
	Dmid         int64  `db:"dmid" json:"dmid"`
	MsgType      int    `db:"msg_type" json:"msg_type"`
	Bubble       string `db:"bubble" json:"bubble"`
	Content      string `db:"content" json:"content"`
	Mid          int64  `db:"mid" json:"mid"`
	Uname        string `db:"uname" json:"uname"`
	RoomAdmin    int    `db:"room_admin" json:"room_admin"`
	Vip          int    `db:"vip" json:"vip"`
	Svip         int    `db:"svip" json:"svip"`
	Rank         int    `db:"rank" json:"rank"`
	MobileVerify int    `db:"mobile_verify" json:"mobile_verify"`
	UnameColor   string `db:"uname_color" json:"uname_color"`
	MedalName    string `db:"medal_name" json:"medal_name"`
	UpName       string `db:"up_name" json:"up_name"`
	MedalLevel   int    `db:"medal_level" json:"medal_level"`
	UserLevel    int    `db:"user_level" json:"user_level"`
}

type Hot struct {
	Time    int64 `db:"time" json:"time"`
	Hot     int64 `db:"hot" json:"hot"`
	Watched int64 `db:"watched" json:"watched"`
}

type Gift struct {
	Id        int64  `db:"id" json:"id"`
	Uname     string `db:"uname" json:"uname"`
	Uid       int64  `db:"uid" json:"uid"`
	Rnd       string `db:"rnd" json:"rnd"`
	GiftName  string `db:"gift_name" json:"gift_name"`
	GiftNum   int    `db:"gift_num" json:"gift_num"`
	GiftID    int64  `db:"gift_id" json:"gift_id"`
	Action    string `db:"action" json:"action"`
	GiftPrice int    `db:"gift_price" json:"gift_price"`
	Time      int64  `db:"time" json:"time"`
	Num       int    `db:"num" json:"num"`
}

type Sc struct {
	ScId      int64  `db:"sc_id" json:"sc_id"`
	ID        int64  `db:"id" json:"id"`
	Uname     string `db:"uname" json:"uname"`
	Uid       int64  `db:"uid" json:"uid"`
	Time      int64  `db:"time" json:"time"`
	StartTime int64  `db:"start_time" json:"start_time"`
	Message   string `db:"message" json:"message"`
	Price     int    `db:"price" json:"price"`
}

type Viewer struct {
	ViewerId   int64  `db:"viewer_id" json:"viewer_id"`
	Uid        int64  `db:"uid" json:"uid"`
	Uname      string `db:"uname" json:"uname"`
	Time       int64  `db:"time" json:"time"`
	Score      int64  `db:"score" json:"score"`
	Dmscore    int    `db:"dmscore" json:"dmscore"`
	Medallevel int64  `db:"medallevel" json:"medallevel"`
	Medalname  string `db:"medalname" json:"medalname"`
	Targetid   int64  `db:"targetid" json:"targetid"`
}
