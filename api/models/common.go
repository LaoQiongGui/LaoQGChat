package models

type RestOutDto struct {
	Common RestCommonDto `json:"common"`
	Data   any           `json:"data"`
}

type RestCommonDto struct {
	Status      int    `json:"status"`
	MessageCode string `json:"message_code"`
	MessageText string `json:"message_text"`
}
