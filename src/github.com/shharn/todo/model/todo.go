package model

var NextId = 1

type Todo struct {
	Id int `json:"id"`
	Content string `json:"content,omitempty"`
	ParentIds []int `json:"parentIds,omitempty"`
	RootId int `json:"rootId,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	ModifiedAt string `json:"modifiedAt,omitempty"`
	IsComplete bool `json:"isComplete,omitempty"`
}