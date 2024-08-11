package response

type Response struct {
	ResourceID       string `json:"id,omitempty"`
	Status           string `json:"status,omitempty"`
	PrivateKey       string `json:"private_key,omitempty"`
	Description      string `json:"description,omitempty"`
	ErrorDescription string `json:"error,omitempty"`
}
