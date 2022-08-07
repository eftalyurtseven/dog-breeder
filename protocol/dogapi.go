package protocol

// DogApiResp is the response from the dog api
type DogApiResp struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Code    int    `json:"code,omitempty"`
}
