package ginserver

import "fmt"

const (
	success = "success"
)

// State base state struct
type State struct {
	// 0 success, other error
	Code int `json:"code"`
	// success or error message
	Msg string `json:"msg"`
}

// Response base response struct
type Response struct {
	State State `json:"state"`
	// multiple data is a list, single data is a object
	Data interface{} `json:"data,omitempty"`
}

// ExceptResponse except response
func ExceptResponse(code int, msg ...interface{}) *Response {
	return &Response{
		State: State{
			Code: code,
			Msg:  fmt.Sprint(msg...),
		},
	}
}

// DataResponse single data response
func DataResponse(data interface{}) *Response {
	return &Response{
		State: State{
			Code: 0,
			Msg:  success,
		},
		Data: data,
	}
}

// ListData multiple data struct
type ListData struct {
	// multiple data
	Items interface{} `json:"items,omitempty"`
	// total count
	Total int `json:"total"`
}

// ListResponse multiple data response
func ListResponse(total int, data interface{}) *Response {
	return &Response{
		State: State{
			Code: 0,
			Msg:  success,
		},
		Data: ListData{
			Items: data,
			Total: total,
		},
	}
}

// OkResponse success response without data
func OkResponse() *Response {
	return &Response{
		State: State{
			Code: 0,
			Msg:  success,
		},
	}
}
