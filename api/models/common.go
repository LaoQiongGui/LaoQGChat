package models

type Response struct {
	Common ResponseCommon `json:"common"`
	Data   any            `json:"data"`
}

type ResponseCommon struct {
	Status      int    `json:"status"`
	MessageCode string `json:"message_code"`
	MessageText string `json:"message_text"`
}

const (
	ResponseCommonStatusSuccess      = 0
	ResponseCommonStatusWarning      = 100
	ResponseCommonStatusServiceError = 200
	ResponseCommonStatusSystemError  = 990
)

var (
	ResponseCommonSuccess = ResponseCommon{
		Status:      ResponseCommonStatusSuccess,
		MessageCode: "N000000",
		MessageText: "",
	}

	ResponseCommonSystemError = ResponseCommon{
		Status:      ResponseCommonStatusSystemError,
		MessageCode: "E999999",
		MessageText: "System Error",
	}
)
