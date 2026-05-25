package general_dto

type Element struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Origin struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
