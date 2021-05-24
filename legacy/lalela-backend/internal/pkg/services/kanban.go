package services

import (
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
	"strconv"
	"time"
)

// Scoped Models

type KanbanCardInfo struct {
	CardId    uint
	Action    string
	Value     string
	UserId    uint
	CreatedOn time.Time
}

// Gets
func GetKanbanLists() []models.KanbanListReturn {
	db := utils.GetDB()
	var lists []models.KanbanList
	db.Find(&lists)

	var listsFinal []models.KanbanListReturn

	for _, list := range lists {

		listsFinal = append(listsFinal, models.KanbanListReturn{
			Id:      list.Id,
			Name:    list.Name,
			CardIds: GetKanbanCardsByListId(list.Id),
		})
	}

	return listsFinal
}

func GetKanbanCards() []models.KanbanCard {
	db := utils.GetDB()
	var cards []models.KanbanCard
	db.Find(&cards)
	var cardsFinal []models.KanbanCard

	for _, card := range cards {
		card.MembersId = GrabAllAssigned(card.Id)
		card.History = GetKanbanEventsForCard(card)
		cardsFinal = append(cardsFinal, card)
	}
	return cardsFinal
}

func GetKanbanMembers() []models.KanbanMemberReturn {
	db := utils.GetDB()
	var members []models.KanbanMember
	db.Find(&members)

	var membersFinal []models.KanbanMemberReturn

	for _, member := range members {
		i, _ := FindByIdRaw(member.UserID)
		membersFinal = append(membersFinal, models.KanbanMemberReturn{
			Id:     i.ID,
			Avatar: i.Avatar,
			Name:   KanbanMemberBuildName(i),
			Email:  i.Email,
		})
	}

	return membersFinal
}

func GetKanbanEventsForCard(card models.KanbanCard) []models.KanbanEvents {
	db := utils.GetDB()
	var events []models.KanbanEvents
	db.Order("id asc").Where("card_id = ?", card.Id).Find(&events)

	return events
}

func GetKanbanEvents(args *models.KanbanGetEventsAllRequest) []models.KanbanEvents {
	db := utils.GetDB()
	var events []models.KanbanEvents
	db.Order(args.Column + " " + args.Order).Limit(args.Limit).Find(&events)

	return events
}

func GetKanbanCardsByListId(listId uint) []string {
	db := utils.GetDB()
	var cards []models.KanbanCard
	var cardsFinal []string
	db.Where("list_id = ?", listId).Find(&cards)

	for _, card := range cards {
		cardsFinal = append(cardsFinal, strconv.Itoa(int(card.Id)))
	}
	return cardsFinal
}

func GetKanbanCard(card models.KanbanCard) models.KanbanCard {
	db := utils.GetDB()
	db.Find(&card, card)
	return card
}

func GetKanbanCardById(cardId uint) models.KanbanCard {
	db := utils.GetDB()
	var card models.KanbanCard
	db.Where("Id = ?", cardId).First(&card)
	card.MembersId = GrabAllAssigned(card.Id)
	return card
}

// Actions
func CreateKanbanCard(card models.KanbanCard, userId uint) (string, bool) {
	db := utils.GetDB()
	if CheckIfKanbanCardExists(card) {

		cardTemp := GetKanbanCard(card)
		cardTemp.Repeat++
		db.Save(&cardTemp)
		AddKanbanEvent(cardTemp.Id, "Card Repeated", strconv.Itoa(cardTemp.Repeat), userId)
		return "Exists: Updated", true
	} else {

		card.CreatedOn = time.Now()
		card.Due = time.Now().AddDate(0, 0, 7)
		db.Create(&card)
		AddKanbanEvent(card.Id, "Card Created", "List::"+strconv.Itoa(int(card.ListId)), userId)
		return "Created", true
	}

}

func AssignKanbanCard(card models.KanbanCard, userId uint, action string, newUser uint) string {
	db := utils.GetDB()
	db.First(&card)

	if action == "ADD" {
		s := CheckIfAssignedAlready(GrabAllAssigned(card.Id), newUser)
		if !s {
			tc := models.KanbanAssignLookUp{
				CardId: card.Id,
				UserId: newUser,
			}
			db.Save(&tc)
			AddKanbanEvent(card.Id, "Assigned", strconv.Itoa(int(newUser)), userId)
			return "Assigned " + strconv.Itoa(int(newUser)) + " to " + strconv.Itoa(int(card.Id))
		} else {
			return "Already Assigned " + strconv.Itoa(int(newUser)) + " to " + strconv.Itoa(int(card.Id))
		}

	} else if action == "REMOVE" {
		s := CheckIfAssignedAlready(GrabAllAssigned(card.Id), newUser)
		if s {
			tc := models.KanbanAssignLookUp{}
			db.Where(&models.KanbanAssignLookUp{
				CardId: card.Id,
				UserId: newUser,
			}).First(&tc)
			db.Delete(&tc)
			AddKanbanEvent(card.Id, "Removed", strconv.Itoa(int(newUser)), userId)
			return "Removed " + strconv.Itoa(int(newUser)) + " From " + strconv.Itoa(int(card.Id))
		} else {
			return strconv.Itoa(int(newUser)) + " Not Assigned to " + strconv.Itoa(int(card.Id))
		}
	} else {
		return "Nope"
	}
}

func MoveCardToList(userID uint, cardID uint, listID uint) string {
	db := utils.GetDB()

	cardT := models.KanbanCard{Id: cardID}

	db.Find(&cardT)
	if cardT.ListId == listID {
		return "Card already part of " + strconv.Itoa(int(listID))
	} else {
		cardT.ListId = listID
		AddKanbanEvent(cardT.Id, "Moved to List", strconv.Itoa(int(listID)), userID)
		db.Save(&cardT)
		return "Card Has Been Moved To " + strconv.Itoa(int(listID))
	}

}

func UpdateCard(card models.KanbanCard, userId uint) string {
	db := utils.GetDB()

	cardT := models.KanbanCard{Id: card.Id}

	db.Find(&cardT)

	cardT.Name = card.Name
	cardT.Branch = card.Branch
	cardT.Description = card.Description
	cardT.Division = card.Division
	cardT.Region = card.Region
	cardT.Status = card.Status

	db.Save(&cardT)
	AddKanbanEvent(cardT.Id, "Updated", "NULL", userId)
	return "Card Has Been Updated"
}

func AddKanbanEvent(card uint, action string, value string, userid uint) bool {
	db := utils.GetDB()
	kanbanEventNew := models.KanbanEvents{
		CardId:    card,
		Action:    action,
		Value:     value,
		CreatedOn: time.Now(),
		UserId:    userid,
	}

	if err := db.Create(&kanbanEventNew).Error; err != nil {
		return false
	} else {
		return true
	}
}

func GrabAllAssigned(cardID uint) []uint {
	db := utils.GetDB()
	var assigned []models.KanbanAssignLookUp
	var assignedFinal []uint
	db.Where("card_id = ?", cardID).Find(&assigned)

	for _, assign := range assigned {

		assignedFinal = append(assignedFinal, assign.UserId)
	}
	return assignedFinal
}

func CheckIfKanbanCardExists(card models.KanbanCard) bool {
	db := utils.GetDB()
	var cards models.KanbanCard
	if err := db.Where(&card).Find(&cards).Error; err != nil {
		return false
	} else {
		return true
	}
}

func CheckIfAssignedAlready(memberId []uint, user uint) bool {
	found := FindUint(memberId, user)
	if !found {
		return false
	} else {
		return true
	}
}

// Util
func KanbanSerializer() {

}

func FindUint(slice []uint, val uint) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func RemoveUint(slice []uint, val uint) []uint {
	slice[val] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func KanbanMemberBuildName(i *models.User) string {
	if i.LastName == "" {
		return i.FirstName
	} else {
		return i.FirstName + " " + i.LastName
	}
}
