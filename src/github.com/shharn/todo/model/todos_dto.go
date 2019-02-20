package model

type TodosDto struct {
	Todos []Todo `json:"todos"`
	TotalCount int `json:"totalCount"`
}