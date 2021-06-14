package store

import (
	"context"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/database"
	"lalela-backend/internal/pkg/organizations"
	"time"
)



func CreateOne(orgName string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var org organizations.Organizations
	var collection = database.OpenCollection("organizations")

	org.ID = primitive.NewObjectID()
	org.Name = orgName

	org.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	org.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	insertId, err := collection.InsertOne(ctx, org)
	if err != nil {
		cancel()
		return "", err
	}
	defer cancel()

	return cast.ToString(insertId.InsertedID), err
}

//func getUserOrganizations(email string) ([]organizations.Organizations, error) {
//	foundUser, err := store.FindUserByEmail(email)
//	if err != nil {
//		return nil, err
//	}
//
//	orgs, err := auth.Enforcer.GetDomainsForUser(foundUser.ID.Hex())
//	if err != nil {
//		return nil, err
//	}
//
//	var collection = database.OpenCollection("organizations")
//	var orgz []organizations.Organizations
//
//	for _, orgId := range orgs {
//		var organization organizations.Organizations
//		err := collection.FindOne(context.Background(), bson.M{"_id": orgId}).Decode(&organization)
//		if err != nil {
//			return []organizations.Organizations{}, err
//		}
//		orgz = append(orgz, organization)
//	}
//	return orgz, nil
//}