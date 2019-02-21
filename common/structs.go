package common

// SspResponse is convert to json
type SspResponse struct {
	URL string `json:"url"`
}

// DspResponse is convert to json
type DspResponse struct {
	RequestID string `json:"request_id"`
	URL       string `json:"url"`
	Price     int    `json:"price"`
}

// DspRequest is convert to json
type DspRequest struct {
	SspName     string `json:"ssp_name"`
	RequestTime string `json:"request_time"`
	RequestID   string `json:"request_id"`
	AppID       int    `json:"app_id"`
}

// WinNotice is convert to json
type WinNotice struct {
	RequestID string `json:"request_id"`
	Price     int    `json:"price"`
}


// PriceInfo is convert to json
type PriceInfo struct {
	DspHost     string
	DspResponse DspResponse
}
