package dto

type RestOutDto struct {
	Common struct {
		Status      int    `json:"status"`
		MessageCode string `json:"message_code"`
		MessageText string `json:"message_text"`
	} `json:"common"`
	Data any `json:"data"`
}
