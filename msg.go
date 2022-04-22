package main

import "encoding/json"

type comboSend struct {
	Action         string             `json:"action"`
	BatchComboID   string             `json:"batch_combo_id"`
	BatchComboNum  int                `json:"batch_combo_num"`
	ComboID        string             `json:"combo_id"`
	ComboNum       int                `json:"combo_num"`
	ComboTotalCoin int                `json:"combo_total_coin"`
	Dmscore        int                `json:"dmscore"`
	GiftID         int64              `json:"gift_id"`
	GiftName       string             `json:"gift_name"`
	GiftNum        int                `json:"gift_num"`
	IsShow         int                `json:"is_show"`
	MedalInfo      comboSendMedalInfo `json:"medal_info"`
	NameColor      string             `json:"name_color"`
	RUname         string             `json:"r_uname"`
	Ruid           int                `json:"ruid"`
	SendMaster     interface{}        `json:"send_master"`
	TotalNum       int                `json:"total_num"`
	UID            int64              `json:"uid"`
	Uname          string             `json:"uname"`
}

type comboSendMedalInfo struct {
	AnchorRoomid     int    `json:"anchor_roomid"`
	AnchorUname      string `json:"anchor_uname"`
	GuardLevel       int    `json:"guard_level"`
	IconID           int    `json:"icon_id"`
	IsLighted        int    `json:"is_lighted"`
	MedalColor       int    `json:"medal_color"`
	MedalColorBorder int    `json:"medal_color_border"`
	MedalColorEnd    int    `json:"medal_color_end"`
	MedalColorStart  int    `json:"medal_color_start"`
	MedalLevel       int    `json:"medal_level"`
	MedalName        string `json:"medal_name"`
	Special          string `json:"special"`
	TargetID         int    `json:"target_id"`
}

func getData(raw []byte) json.RawMessage {
	var d struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &d); err != nil {
		return []byte{}
	}
	return d.Data
}
