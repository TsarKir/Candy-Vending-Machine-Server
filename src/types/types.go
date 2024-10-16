package types

type Request struct {
	Money      int    `json:"money,omitempty"`
	CandyType  string `json:"candyType,omitempty"`
	CandyCount int    `json:"candyCount,omitempty"`
}

type Response struct {
	Thanks string `json:"thanks,omitempty"`
	Change int    `json:"change,omitempty"`
	Error  string `json:"error,omitempty"`
}
