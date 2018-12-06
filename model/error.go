package model

type Error struct {
	IsErr   bool   `json:"isErr"`
	Message string `json:"message"`
}

func NewError(errMsg string) Error {
	return Error{
		IsErr:   true,
		Message: errMsg,
	}
}
