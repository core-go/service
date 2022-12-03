package core

import (
	"context"
	"strings"
	"unicode"
)

type ResultInfo struct {
	Status  int            `yaml:"status" mapstructure:"status" json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status"`
	Errors  []ErrorMessage `yaml:"errors" mapstructure:"errors" json:"errors,omitempty" gorm:"column:errors" bson:"errors,omitempty" dynamodbav:"errors,omitempty" firestore:"errors,omitempty"`
	Value   interface{}    `yaml:"value" mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	Message string         `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

type ErrorMessage struct {
	Field   string `yaml:"field" mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `yaml:"code" mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `yaml:"param" mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}
type ErrorDetail struct {
	ErrorField string `yaml:"error_field" mapstructure:"error_field" json:"errorField,omitempty" gorm:"column:error_field" bson:"errorField,omitempty" dynamodbav:"errorField,omitempty" firestore:"errorField,omitempty"`
	ErrorCode  string `yaml:"error_code" mapstructure:"error_code" json:"errorCode,omitempty" gorm:"column:error_code" bson:"errorCode,omitempty" dynamodbav:"errorCode,omitempty" firestore:"errorCode,omitempty"`
	ErrorDesc  string `yaml:"error_desc" mapstructure:"error_desc" json:"errorDesc,omitempty" gorm:"column:error_desc" bson:"errorDesc,omitempty" dynamodbav:"errorDesc,omitempty" firestore:"errorDesc,omitempty"`
}
type ErrorDetails struct {
	ErrorDetails []ErrorDetail `yaml:"error_details" mapstructure:"error_details" json:"errorDetails,omitempty" gorm:"column:error_details" bson:"errorDetails,omitempty" dynamodbav:"errorDetails,omitempty" firestore:"errorDetails,omitempty"`
}
type Validator interface {
	Validate(ctx context.Context, model interface{}) ([]ErrorMessage, error)
}
type MapValidator interface {
	Validate(ctx context.Context, model map[string]interface{}) ([]ErrorMessage, error)
}

func RemoveRequiredError(errors []ErrorMessage) []ErrorMessage {
	if errors == nil || len(errors) == 0 {
		return errors
	}
	errs := make([]ErrorMessage, 0)
	for _, s := range errors {
		if s.Code != "required" && !strings.HasPrefix(s.Code, "minlength") {
			errs = append(errs, s)
		} else if strings.Index(s.Field, ".") >= 0 {
			errs = append(errs, s)
		}
	}
	return errs
}
func FormatErrorField(s string) string {
	splitField := strings.Split(s, ".")
	length := len(splitField)
	if length == 1 {
		return lcFirstChar(splitField[0])
	} else if length > 1 {
		var tmp []string
		for _, v := range splitField[1:] {
			tmp = append(tmp, lcFirstChar(v))
		}
		return strings.Join(tmp, ".")
	}
	return s
}
func lcFirstChar(s string) string {
	if len(s) > 0 {
		runes := []rune(s)
		runes[0] = unicode.ToLower(runes[0])
		return string(runes)
	}
	return s
}
func BuildErrorDetails(errors []ErrorMessage) []ErrorDetail {
	errs := make([]ErrorDetail, 0)
	if errors == nil || len(errors) == 0 {
		return errs
	}
	for _, s := range errors {
		d := ErrorDetail{ErrorCode: s.Code, ErrorField: s.Field, ErrorDesc: s.Message}
		errs = append(errs, d)
	}
	return errs
}
