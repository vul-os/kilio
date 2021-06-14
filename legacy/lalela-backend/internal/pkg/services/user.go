package services

import (
	"lalela-backend/internal/pkg/middleware"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
	"fmt"
	"github.com/o1egl/govatar"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GetUsers(user *models.User) ([]models.User, error) {

	// Init DB
	db := utils.GetDB()

	// Init RBCA Util
	e := utils.InitRBCA()

	// Grab User ID
	db.Where("email = ?", user.Email).First(user)

	// Define Slices
	var users []models.User
	var userRows []models.User

	// Check if User is allowed to view
	if utils.PermissionCanViewUser(e, fmt.Sprint(user.ID), fmt.Sprint(user.UserGroupId)) || user.UserGroupId == 1 {

		/// true
		if user.UserGroupId == 1 {

			// Grab All Users
			db.Find(&users)

		} else {

			// Grab All Users Within Same Group
			db.Where("user_group_id = ?", user.UserGroupId).Find(&users)
		}

		// loop
		for _, v := range users {

			// Create struct with ID of User ID
			result := models.User{ID: v.ID}

			// Fill in struct based on User ID
			db.Where("id = ?", v.ID).First(&result)

			// Add to Slice
			userRows = append(userRows, result)

		}

		return userRows, nil

	} else {

		return userRows, nil

	}
}
func GetUser(user models.User) (models.UserGetReport, error) {

	// Init DB
	db := utils.GetDB()

	// Grab User ID
	db.First(&user)
	userFinal := models.UserGetReport{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		UserGroupId: GetUserGroup(user.UserGroupId),
		Avatar:      user.Avatar,
		IsTrail:     user.IsTrail,
		TrailExpire: user.TrailExpire,
		LastLogin:   GetUserLastLogin(user.ID),
	}

	return userFinal, nil

}
func GetUserRaw(user models.User) (models.User, error) {

	// Init DB
	db := utils.GetDB()

	// Grab User ID
	db.Where("email = ?", user.Email).First(&user)

	return user, nil

}

func GetUserLastLogin(id uint) string {
	db := utils.GetDB()
	var event models.Event
	db.Where("user_id = ?", id).Last(&event)
	return event.Date.Format("2006-01-02 15:04:05")
}

func GetUserGroup(id int) string {
	db := utils.GetDB()
	var group models.Group
	db.Where("id = ?", id).Find(&group)
	return group.Name
}

func GetUsersRaw() []models.User {
	db := utils.GetDB()
	var users []models.User
	db.Find(&users)
	return users
}

func GetUserExists(user *models.User) (bool, error) {
	db := utils.GetDB()
	var users models.User
	if len(user.Email) != 0 {
		if err := db.Where(&models.User{Email: user.Email}).Find(&users).Error; err != nil {
			log.Printf("Error: %s", err)
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return true, nil
	}
}

func AddUser(user *models.User) bool {
	db := utils.GetDB()

	// todo :: Find a better way to do emtpy checks

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Print(middleware.NewError(err))
		return false
	}

	user.Password = string(pass)

	if user.FirstName == "" || user.Email == "" {
		return false
	} else {
		db.Create(&user)
		// Todo check if added
		return true
	}
	return false
}

func CreateUser(user *models.User) (*models.User, error) {
	db := utils.GetDB().Table("users")
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(middleware.NewError(err))
		return nil, err
	}

	user.Password = string(pass)
	user.RoleID = 3

	rand.Seed(time.Now().UnixNano())

	user.ResetToken = ""
	user.ResetTokenExpiry = time.Now()

	userDb := db.Create(user)

	if userDb.Error != nil {
		log.Print(middleware.NewError(userDb.Error))
		fmt.Println(userDb.Error)
		return nil, userDb.Error
	}

	return user, nil
}

func UpdateUser(user *models.User) {
	db := utils.GetDB().Table("users")
	db.Save(&user)
}

func DeleteUser(user *models.User) {
	db := utils.GetDB().Table("users")
	db.Delete(&user)
}

func GetTokensPassword() (string, time.Time) {
	resetToken := GenerateToken()
	timein := time.Now().Add(time.Hour*0 + time.Minute*10 + time.Second*0)

	return resetToken, timein
}

func GenerateToken() string {
	rand.Seed(time.Now().UnixNano())
	resetToken := randSeq(25)
	return resetToken
}

func FindByEmail(user *models.User) (*models.User, error) {
	db := utils.GetDB().Table("users")
	userDb := &models.User{}
	fmt.Println(user.Email)
	if err := db.Where("Email = ?", user.Email).First(userDb).Error; err != nil {
		log.Print(middleware.NewError(err))
		fmt.Println(err)
		return nil, err
	}

	return userDb, nil
}

func SetAvatar(user *models.User, r *http.Request) (*models.User, error) {

	err := govatar.GenerateFileForUsername(govatar.MALE, user.Email, "./avatars/"+user.Email+".jpg")
	if err != nil {
		log.Print(middleware.NewError(err))
	}
	user.Avatar = "http://" + r.Host + "/avatars/" + user.Email + ".jpg"
	UpdateUser(user)
	return user, nil
}

// Added To Not Mess With Inital Func
// todo :: Remove inital function
func FindByEmailString(email string) (*models.User, error) {
	db := utils.GetDB().Table("users")
	user := models.User{Email: email}
	userDb := &models.User{}
	if err := db.Where("Email = ?", user.Email).First(userDb).Error; err != nil {
		log.Print(middleware.NewError(err))
		fmt.Println(err)
		return nil, err
	}
	return userDb, nil
}

func FindByToken(token string) (*models.User, error) {
	db := utils.GetDB().Table("users")
	userDb := &models.User{}

	if err := db.Where("reset_token = ?", token).First(userDb).Error; err != nil {
		// var resp = map[string]interface{}{"status": false, "message": "Invalid Token"}
		// json.NewEncoder(w).Encode(resp)
		// return
		log.Print(middleware.NewError(err))
		return nil, err
	}

	return userDb, nil
}

func FindByValidationToken(token string) (*models.User, error) {
	db := utils.GetDB().Table("users")
	userDb := &models.User{}

	if err := db.Where("validation_token = ?", token).First(userDb).Error; err != nil {
		// var resp = map[string]interface{}{"status": false, "message": "Invalid Token"}
		// json.NewEncoder(w).Encode(resp)
		// return
		log.Print(middleware.NewError(err))
		return nil, err
	}

	return userDb, nil
}

func CheckPassword(user *models.User, pw string) error {
	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		log.Print(middleware.NewError(errf))
		fmt.Println(errf)
		return errf
	}
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GiveUserSubAdmin(user uint, group int) {
	utils.GroupPolicyExistsAdd("user::"+fmt.Sprint(user), "role::3", fmt.Sprint(group))
}

func FindById(id int) (string, error) {
	db := utils.GetDB().Table("users")
	userDb := &models.User{}

	if err := db.Where("id = ?", id).First(userDb).Error; err != nil {
		log.Print(middleware.NewError(err))
		return "nil", err
	}

	return userDb.Email, nil
}

func FindByIdRaw(id uint) (*models.User, error) {
	db := utils.GetDB().Table("users")
	userDb := &models.User{}

	if err := db.Where("id = ?", id).First(userDb).Error; err != nil {
		log.Print(middleware.NewError(err))
		return nil, err
	}

	return userDb, nil
}

func FindByIdReturnGroup(id int) (string, error) {
	db := utils.GetDB()
	userDb := &models.User{}
	groupD := &models.Group{}

	if err := db.Where("id = ?", id).First(userDb).Error; err != nil {
		log.Print(middleware.NewError(err))
		return "nil", err
	}

	if err := db.Where("id = ?", userDb.UserGroupId).First(&groupD).Error; err != nil {
		log.Print(middleware.NewError(err))
		return "nil", err
	}

	return groupD.Name, nil
}

func FindByIdReturnGroupName(id int) (string, error) {
	db := utils.GetDB()
	groupD := &models.Group{}

	if err := db.Where("id = ?", id).First(&groupD).Error; err != nil {
		log.Print(middleware.NewError(err))
		return "nil", err
	}

	return groupD.Name, nil
}

func UserPermission(ar []string, id int, i uint) {
	var roles = GetRoles()

	for ks, role := range roles {

		_, found := Find(ar, role)
		if found {
			utils.GroupPolicyExistsAdd("user::"+fmt.Sprint(i), ks, fmt.Sprint(id))
		} else {
			utils.GroupPolicyExistsRemove("user::"+fmt.Sprint(i), ks, fmt.Sprint(id))
		}

	}
}

func Find(slice []string, val string) (string, bool) {
	for _, item := range slice {
		if item == val {
			return item, true
		}
	}
	return "", false
}

func GetUserCount() int64 {
	var count int64
	db := utils.GetDB()
	var users []models.User
	db.Find(&users).Count(&count)
	return count
}

func IsUserAdmin(user models.User) bool {
	if user.UserGroupId == 1 {
		if user.Email == "frans.reichert@cas.group" {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
	return false
}
