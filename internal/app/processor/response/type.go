package response

type ProcessedResponse interface {
	SetResourceID(s string) *Response
	SetStatus(s string) *Response
	SetPrivateKey(s string) *Response
	SetDescription(s string) *Response
	SetErrorDescription(s string) *Response
	Build() *Response
}

func NewResponseBuilder() ProcessedResponse {
	return &Response{}
}

func (r *Response) SetResourceID(s string) *Response {
	r.ResourceID = s
	return r
}

func (r *Response) SetStatus(s string) *Response {
	r.Status = s
	return r
}

func (r *Response) SetDescription(s string) *Response {
	r.Description = s
	return r
}

func (r *Response) SetErrorDescription(s string) *Response {
	r.ErrorDescription = s
	return r
}

func (r *Response) SetPrivateKey(s string) *Response {
	r.PrivateKey = s
	return r
}

func (r *Response) Build() *Response {
	return r
}
