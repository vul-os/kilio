package models

import "time"

//// KanBan
type KanbanMember struct {
	Id     uint `json:"id" gorm:"AUTO_INCREMENT"`
	UserID uint `json:"user_id"`
}

type KanbanList struct {
	Id   uint   `json:"id" gorm:"AUTO_INCREMENT"`
	Name string `json:"name"`
}

type KanbanCard struct {
	Id          uint           `json:"id" gorm:"AUTO_INCREMENT"`
	Cover       string         `json:"cover"`
	Description string         `json:"description"`
	CreatedOn   time.Time      `json:"created_on"`
	Due         time.Time      `json:"due"`
	MembersId   []uint         `json:"memberIds" gorm:"-" `
	History     []KanbanEvents `json:"history" gorm:"-" `
	ListId      uint           `json:"list_id"`
	Name        string         `json:"name"`
	Status      string         `json:"status"`
	Client      string         `json:"client"`
	Branch      string         `json:"branch"`
	Region      string         `json:"region"`
	Division    string         `json:"division"`
	MValue      int            `json:"m_value"`
	Repeat      int            `json:"repeat"`
}

type KanbanEvents struct {
	Id        uint      `json:"id" gorm:"AUTO_INCREMENT"`
	CardId    uint      `json:"cardId"`
	Action    string    `json:"action"`
	Value     string    `json:"value"`
	CreatedOn time.Time `json:"createdOn"`
	UserId    uint      `json:"userId"`
}

type KanbanAssignLookUp struct {
	Id     uint `json:"id" gorm:"AUTO_INCREMENT"`
	CardId uint `json:"card_id"`
	UserId uint `json:"user_id"`
}

// Kanban API Calls
type KanbanCreateCardRequest struct {
	UserId uint       `json:"user_id"`
	Card   KanbanCard `json:"card"`
}
type KanbanCreateCardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type KanbanGetCardDetailsRequest struct {
	UserId uint `json:"user_id"`
	CardId uint `json:"card_id"`
}
type KanbanGetCardDetailsResponse struct {
	Card        KanbanCard     `json:"card"`
	CardHistory []KanbanEvents `json:"card_history"`
}

type KanbanGetEventsAllRequest struct {
	UserId uint   `json:"user_id"`
	Limit  int    `json:"limit"`
	Column string `json:"column"`
	Order  string `json:"order"`
}
type KanbanGetEventsAllResponse struct {
	Events []KanbanEvents `json:"events"`
}

type KanbanGetBoardRequest struct {
	UserId uint `json:"user_id"`
}
type KanbanGetBoardResponse struct {
	Cards   []KanbanCard         `json:"cards"`
	Lists   []KanbanListReturn   `json:"lists"`
	Members []KanbanMemberReturn `json:"members"`
}

type KanbanAssignCardRequest struct {
	UserId     uint   `json:"user_id"`
	AssigneeId uint   `json:"assignee_id"`
	CardId     uint   `json:"card_id"`
	Action     string `json:"action"`
}
type KanbanAssignCardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type KanbanMoveCardRequest struct {
	UserId uint `json:"user_id"`
	CardId uint `json:"card_id"`
	ListID uint `json:"list_id"`
}
type KanbanMoveCardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type KanbanUpdateCardRequest struct {
	UserId uint       `json:"user_id"`
	Card   KanbanCard `json:"card"`
}
type KanbanUpdateCardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type KanbanGetKanbanMembersRequest struct {
	UserId uint `json:"user_id"`
}
type KanbanGetKanbanMembersResponse struct {
	Members []KanbanMemberReturn `json:"members"`
}

// Kanban return interfaces
type KanbanMemberReturn struct {
	Id     uint   `json:"id"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type KanbanListReturn struct {
	Id      uint     `json:"id"`
	Name    string   `json:"name"`
	CardIds []string `json:"cardIds"`
}

