package controllers

import (
	models "lalela-backend/internal/pkg/models"
	services "lalela-backend/internal/pkg/services"
	"log"
	"net/http"
)

type KanbanCon struct{}

// Done

func (t *KanbanCon) GetBoard(r *http.Request, args *models.KanbanGetBoardRequest, reply *models.KanbanGetBoardResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	reply.Cards = services.GetKanbanCards()
	reply.Lists = services.GetKanbanLists()
	reply.Members = services.GetKanbanMembers()

	return nil
}
func (t *KanbanCon) GetCardDetails(r *http.Request, args *models.KanbanGetCardDetailsRequest, reply *models.KanbanGetCardDetailsResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}
	reply.Card = services.GetKanbanCardById(args.CardId)
	reply.CardHistory = services.GetKanbanEventsForCard(reply.Card)
	return nil
}
func (t *KanbanCon) CreateCard(r *http.Request, args *models.KanbanCreateCardRequest, reply *models.KanbanCreateCardResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}
	reply.Message, _ = services.CreateKanbanCard(args.Card, args.UserId)
	return nil
}
func (t *KanbanCon) AssignToCard(r *http.Request, args *models.KanbanAssignCardRequest, reply *models.KanbanAssignCardResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}
	reply.Message = services.AssignKanbanCard(services.GetKanbanCardById(args.CardId), args.UserId, args.Action, args.AssigneeId)
	return nil
}
func (t *KanbanCon) MoveCard(r *http.Request, args *models.KanbanMoveCardRequest, reply *models.KanbanMoveCardResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	reply.Message = services.MoveCardToList(args.UserId, args.CardId, args.ListID)

	return nil
}
func (t *KanbanCon) GetKanbanMembers(r *http.Request, args *models.KanbanGetKanbanMembersRequest, reply *models.KanbanGetKanbanMembersResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	reply.Members = services.GetKanbanMembers()

	return nil
}
func (t *KanbanCon) UpdateCard(r *http.Request, args *models.KanbanUpdateCardRequest, reply *models.KanbanUpdateCardResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	reply.Message = services.UpdateCard(args.Card, args.UserId)
	return nil
}
func (t *KanbanCon) GetKanbanEventsAll(r *http.Request, args *models.KanbanGetEventsAllRequest, reply *models.KanbanGetEventsAllResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	reply.Events = services.GetKanbanEvents(args)
	return nil
}

// Future

func (t *KanbanCon) ArchiveCard(r *http.Request, args *models.KanbanGetBoardRequest, reply *models.KanbanGetBoardResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	return nil
}
