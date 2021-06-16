package casbin

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"lalela-backend/internal/pkg/mongo"
	"runtime"
	"strings"
	"time"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

const defaultTimeout = 30 * time.Second
const defaultDatabaseName string = "casbin"
const defaultCollectionName string = "casbin_rule"

// CasbinRule represents a rule in Casbin.
type CasbinRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

// adapter represents the MongoDB adapter for policy storage.
type adapter struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// finalizer is the destructor for adapter.
func finalizer(a *adapter) {

}

// NewAdapter is the constructor for Adapter. If database name is not provided
// in the Mongo URL, 'casbin' will be used as database name.
// 'casbin_rule' will be used as a collection name.
func NewAdapter(database *mongo.Database) (persist.BatchAdapter, error) {
	collection := database.Collection("casbin")
	a := &adapter{
		collection: collection,
		timeout:    defaultTimeout,
	}

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)
	return a, nil
}

func (a *adapter) dropTable() error {
	//ctx, cancel := context.WithTimeout(context(), a.timeout)
	//defer cancel()
	//
	//err := a.collection.Drop(ctx)
	//if err != nil {
	//	return err
	//}
	return nil
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	var p = []string{line.PType,
		line.V0, line.V1, line.V2, line.V3, line.V4, line.V5}
	var lineText string
	if line.V5 != "" {
		lineText = strings.Join(p, ", ")
	} else if line.V4 != "" {
		lineText = strings.Join(p[:6], ", ")
	} else if line.V3 != "" {
		lineText = strings.Join(p[:5], ", ")
	} else if line.V2 != "" {
		lineText = strings.Join(p[:4], ", ")
	} else if line.V1 != "" {
		lineText = strings.Join(p[:3], ", ")
	} else if line.V0 != "" {
		lineText = strings.Join(p[:2], ", ")
	}

	persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy loads policy from database.
func (a *adapter) LoadPolicy(model model.Model) error {
	line := CasbinRule{}
	ctx, cancel := context.WithTimeout(context.TODO(), a.timeout)
	defer cancel()
	cursor, err := a.collection.GetMongoCollection().Find(ctx, bson.D{})
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		err := cursor.Decode(&line)
		if err != nil {
			return err
		}
		loadPolicyLine(line, model)
	}

	return nil
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{
		PType: ptype,
	}

	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy saves policy to database.
func (a *adapter) SavePolicy(model model.Model) error {
	if err := a.dropTable(); err != nil {
		return err
	}

	var lines []interface{}

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, &line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, &line)
		}
	}
	ctx, cancel := context.WithTimeout(context.TODO(), a.timeout)
	defer cancel()

	if _, err := a.collection.GetMongoCollection().InsertMany(ctx, lines); err != nil {
		return err
	}

	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a *adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)

	ctx, cancel := context.WithTimeout(context.TODO(), a.timeout)
	defer cancel()

	if _, err := a.collection.GetMongoCollection().InsertOne(ctx, line); err != nil {
		return err
	}

	return nil
}

// AddPolicies adds policy rules to the storage.
func (a *adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	var lines []CasbinRule
	for _, rule := range rules {
		line := savePolicyLine(ptype, rule)
		lines = append(lines, line)
	}

	for _, line := range lines {
		ctx, cancel := context.WithTimeout(context.TODO(), a.timeout)
		defer cancel()
		if _, err := a.collection.GetMongoCollection().InsertOne(ctx, line); err != nil {
			return err
		}
	}

	return nil
}

// RemovePolicies removes policy rules from the storage.
func (a *adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	var lines []CasbinRule
	for _, rule := range rules {
		line := savePolicyLine(ptype, rule)
		lines = append(lines, line)
	}

	for _, line := range lines {
		ctx, cancel := context.WithTimeout(context.TODO(), a.timeout)
		defer cancel()
		if _, err := a.collection.GetMongoCollection().DeleteOne(ctx, line); err != nil {
			return err
		}
	}

	return nil
}

// RemovePolicy removes a policy rule from the storage.
func (a *adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)

	ctx, cancel := context.WithTimeout(context.TODO(), a.timeout)
	defer cancel()

	if _, err := a.collection.GetMongoCollection().DeleteOne(ctx, line); err != nil {
		return err
	}

	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	selector := make(map[string]interface{})
	selector["ptype"] = ptype

	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		if fieldValues[0-fieldIndex] != "" {
			selector["v0"] = fieldValues[0-fieldIndex]
		}
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		if fieldValues[1-fieldIndex] != "" {
			selector["v1"] = fieldValues[1-fieldIndex]
		}
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		if fieldValues[2-fieldIndex] != "" {
			selector["v2"] = fieldValues[2-fieldIndex]
		}
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		if fieldValues[3-fieldIndex] != "" {
			selector["v3"] = fieldValues[3-fieldIndex]
		}
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		if fieldValues[4-fieldIndex] != "" {
			selector["v4"] = fieldValues[4-fieldIndex]
		}
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		if fieldValues[5-fieldIndex] != "" {
			selector["v5"] = fieldValues[5-fieldIndex]
		}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), a.timeout)
	defer cancel()

	if _, err := a.collection.GetMongoCollection().DeleteMany(ctx, selector); err != nil {
		return err
	}

	return nil
}
