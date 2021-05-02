package builder

import (
	"context"
	"reflect"
	"strings"
	"time"
)

type TrackingConfig struct {
	Authorization string `mapstructure:"authorization" json:"authorization,omitempty" gorm:"column:authorization" bson:"authorization,omitempty" dynamodbav:"authorization,omitempty" firestore:"authorization,omitempty"`
	User          string `mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	CreatedBy     string `mapstructure:"created_by" json:"createdBy,omitempty" gorm:"column:createdby" bson:"createdBy,omitempty" dynamodbav:"createdBy,omitempty" firestore:"createdBy,omitempty"`
	CreatedAt     string `mapstructure:"created_at" json:"createdAt,omitempty" gorm:"column:createdat" bson:"createdAt,omitempty" dynamodbav:"createdAt,omitempty" firestore:"createdAt,omitempty"`
	UpdatedBy     string `mapstructure:"updated_by" json:"updatedBy,omitempty" gorm:"column:updatedby" bson:"updatedBy,omitempty" dynamodbav:"updatedBy,omitempty" firestore:"updatedBy,omitempty"`
	UpdatedAt     string `mapstructure:"updated_at" json:"updatedAt,omitempty" gorm:"column:updatedat" bson:"updatedAt,omitempty" dynamodbav:"updatedAt,omitempty" firestore:"updatedAt,omitempty"`
}
type ModelBuilder struct {
	GenerateId     func(ctx context.Context, model interface{}) (int, error)
	Authorization  string
	Key            string
	modelType      reflect.Type
	createdByName  string
	createdAtName  string
	updatedByName  string
	updatedAtName  string
	createdByIndex int
	createdAtIndex int
	updatedByIndex int
	updatedAtIndex int
}

func NewDefaultModelBuilderByConfig(generateId func(context.Context) (string, error), modelType reflect.Type, c TrackingConfig) *ModelBuilder {
	if generateId != nil {
		idGenerator := NewIdGenerator(generateId)
		return NewModelBuilderByConfig(idGenerator.Generate, modelType, c)
	} else {
		return NewModelBuilderByConfig(nil, modelType, c)
	}
}
func NewModelBuilderByConfig(generateId func(context.Context, interface{}) (int, error), modelType reflect.Type, c TrackingConfig) *ModelBuilder {
	return NewModelBuilder(generateId, modelType, c.CreatedBy, c.CreatedAt, c.UpdatedBy, c.UpdatedAt, c.User, c.Authorization)
}
func NewDefaultModelBuilder(generateId func(context.Context) (string, error), modelType reflect.Type, options ...string) *ModelBuilder {
	if generateId != nil {
		idGenerator := NewIdGenerator(generateId)
		return NewModelBuilder(idGenerator.Generate, modelType, options...)
	} else {
		return NewModelBuilder(nil, modelType, options...)
	}
}
func NewModelBuilder(generateId func(context.Context, interface{}) (int, error), modelType reflect.Type, options ...string) *ModelBuilder {
	var createdByName, createdAtName, updatedByName, updatedAtName, key, authorization string
	if len(options) > 0 {
		createdByName = options[0]
	}
	if len(options) > 1 {
		createdAtName = options[1]
	}
	if len(options) > 2 {
		updatedByName = options[2]
	}
	if len(options) > 3 {
		updatedAtName = options[3]
	}
	if len(options) > 4 && len(options[4]) > 0 {
		key = options[4]
	} else {
		key = "userId"
	}
	if len(options) > 5 {
		authorization = options[5]
	}
	createdByIndex := findFieldIndex(modelType, createdByName)
	createdAtIndex := findFieldIndex(modelType, createdAtName)
	updatedByIndex := findFieldIndex(modelType, updatedByName)
	updatedAtIndex := findFieldIndex(modelType, updatedAtName)

	return &ModelBuilder{
		GenerateId:     generateId,
		Authorization:  authorization,
		Key:            key,
		modelType:      modelType,
		createdByName:  createdByName,
		createdAtName:  createdAtName,
		updatedByName:  updatedByName,
		updatedAtName:  updatedAtName,
		createdByIndex: createdByIndex,
		createdAtIndex: createdAtIndex,
		updatedByIndex: updatedByIndex,
		updatedAtIndex: updatedAtIndex,
	}
}

func (c *ModelBuilder) BuildToInsert(ctx context.Context, obj interface{}) (interface{}, error) {
	if c.GenerateId != nil {
		_, er0 := c.GenerateId(ctx, obj)
		if er0 != nil {
			return obj, er0
		}
	}
	v := reflect.Indirect(reflect.ValueOf(obj))
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	userId := fromContext(ctx, c.Key, c.Authorization)
	if v.Kind() == reflect.Struct {
		if c.createdByIndex >= 0 {
			createdByField := reflect.Indirect(v).Field(c.createdByIndex)
			if createdByField.Kind() == reflect.Ptr {
				createdByField.Set(reflect.ValueOf(&userId))
			} else {
				createdByField.Set(reflect.ValueOf(userId))
			}
		}

		if c.createdAtIndex >= 0 {
			createdAtField := reflect.Indirect(v).Field(c.createdAtIndex)
			t := time.Now()
			if createdAtField.Kind() == reflect.Ptr {
				createdAtField.Set(reflect.ValueOf(&t))
			} else {
				createdAtField.Set(reflect.ValueOf(t))
			}
		}

		if c.updatedByIndex >= 0 {
			updatedByField := reflect.Indirect(v).Field(c.updatedByIndex)
			if updatedByField.Kind() == reflect.Ptr {
				updatedByField.Set(reflect.ValueOf(&userId))
			} else {
				updatedByField.Set(reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			updatedAtField := v.Field(c.updatedAtIndex)
			t := time.Now()
			if updatedAtField.Kind() == reflect.Ptr {
				updatedAtField.Set(reflect.ValueOf(&t))
				//updatedAtField = reflect.Indirect(updatedAtField)
			} else {
				updatedAtField.Set(reflect.ValueOf(t))
			}
		}
	} else if v.Kind() == reflect.Map {
		var createdByTag, createdAtTag string
		if c.createdByIndex >= 0 {
			if createdByTag = getJsonName(c.modelType, c.createdByIndex); createdByTag == "" || createdByTag == "-" {
				createdByTag = getBsonName(c.modelType, c.createdByIndex)
			}
			if createdByTag != "" && createdByTag != "-" {
				v.SetMapIndex(reflect.ValueOf(createdByTag), reflect.ValueOf(userId))
			}
		}
		if c.createdAtIndex >= 0 {
			if createdAtTag = getJsonName(c.modelType, c.createdAtIndex); createdAtTag == "" || createdAtTag == "-" {
				createdAtTag = getBsonName(c.modelType, c.createdAtIndex)
			}
			if createdAtTag != "" && createdAtTag != "-" {
				v.SetMapIndex(reflect.ValueOf(createdAtTag), reflect.ValueOf(time.Now()))
			}
		}
		var updatedByTag, updatedAtTag string
		if c.updatedByIndex >= 0 {
			if updatedByTag = getJsonName(c.modelType, c.updatedByIndex); updatedByTag == "" || updatedByTag == "-" {
				updatedByTag = getBsonName(c.modelType, c.updatedByIndex)
			}
			if updatedByTag != "" && updatedByTag != "-" {
				v.SetMapIndex(reflect.ValueOf(updatedByTag), reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			if updatedAtTag = getJsonName(c.modelType, c.updatedAtIndex); updatedAtTag == "" || updatedAtTag == "-" {
				updatedAtTag = getBsonName(c.modelType, c.updatedAtIndex)
			}
			if updatedAtTag != "" && updatedAtTag != "-" {
				v.SetMapIndex(reflect.ValueOf(updatedAtTag), reflect.ValueOf(time.Now()))
			}
		}
	}
	return obj, nil
}

func (c *ModelBuilder) BuildToUpdate(ctx context.Context, obj interface{}) (interface{}, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	userId := fromContext(ctx, c.Key, c.Authorization)
	if v.Kind() == reflect.Struct {
		if c.updatedByIndex >= 0 {
			updatedByField := reflect.Indirect(v).Field(c.updatedByIndex)
			if updatedByField.Kind() == reflect.Ptr {
				updatedByField.Set(reflect.ValueOf(&userId))
			} else {
				updatedByField.Set(reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			updatedAtField := v.Field(c.updatedAtIndex)
			t := time.Now()
			if updatedAtField.Kind() == reflect.Ptr {
				updatedAtField.Set(reflect.ValueOf(&t))
				//updatedAtField = reflect.Indirect(updatedAtField)
			} else {
				updatedAtField.Set(reflect.ValueOf(t))
			}
		}
	} else if v.Kind() == reflect.Map {
		var updatedByTag, updatedAtTag string
		if c.updatedByIndex >= 0 {
			if updatedByTag = getJsonName(c.modelType, c.updatedByIndex); updatedByTag == "" || updatedByTag == "-" {
				updatedByTag = getBsonName(c.modelType, c.updatedByIndex)
			}
			if updatedByTag != "" && updatedByTag != "-" {
				v.SetMapIndex(reflect.ValueOf(updatedByTag), reflect.ValueOf(userId))
			}
		}

		if c.updatedAtIndex >= 0 {
			if updatedAtTag = getJsonName(c.modelType, c.updatedAtIndex); updatedAtTag == "" || updatedAtTag == "-" {
				updatedAtTag = getBsonName(c.modelType, c.updatedAtIndex)
			}
			if updatedAtTag != "" && updatedAtTag != "-" {
				v.SetMapIndex(reflect.ValueOf(updatedAtTag), reflect.ValueOf(time.Now()))
			}
		}
	}

	return obj, nil
}

func (c *ModelBuilder) BuildToPatch(ctx context.Context, obj interface{}) (interface{}, error) {
	return c.BuildToUpdate(ctx, obj)
}

func (c *ModelBuilder) BuildToSave(ctx context.Context, obj interface{}) (interface{}, error) {
	return c.BuildToUpdate(ctx, obj)
}
func fromContext(ctx context.Context, key string, authorization string) string {
	if len(authorization) > 0 {
		token := ctx.Value(authorization)
		if token != nil {
			if authorizationToken, exist := token.(map[string]interface{}); exist {
				return fromMap(key, authorizationToken)
			}
		}
		return ""
	} else {
		u := ctx.Value(key)
		if u != nil {
			v, ok := u.(string)
			if ok {
				return v
			}
		}
		return ""
	}
}
func fromMap(key string, data map[string]interface{}) string {
	u := data[key]
	if u != nil {
		v, ok := u.(string)
		if ok {
			return v
		}
	}
	return ""
}
func getBsonName(modelType reflect.Type, index int) string {
	field := modelType.Field(index)
	if tag, ok := field.Tag.Lookup("bson"); ok {
		return strings.Split(tag, ",")[0]
	}
	return ""
}

func getJsonName(modelType reflect.Type, index int) string {
	field := modelType.Field(index)
	if tag, ok := field.Tag.Lookup("json"); ok {
		return strings.Split(tag, ",")[0]
	}
	return ""
}

func findFieldIndex(modelType reflect.Type, fieldName string) int {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		if field.Name == fieldName {
			return i
		}
	}
	return -1
}