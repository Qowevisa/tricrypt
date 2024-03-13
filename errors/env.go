package errors

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ENV_EMPTY    = errors.New("Environment variable was empty")
	NOT_INIT     = errors.New("Is not initialised")
	NOT_HANDLED  = errors.New("Is not handled")
	OUT_OF_BOUND = errors.New("Out of bound")
	NOT_SET      = errors.New("Is not set")
)

func WrapErr(context string, err error) error {
	return fmt.Errorf("%s: %w", context, err)
}

func CheckFieldInitialized(t interface{}, fieldName string) error {
	v := reflect.ValueOf(t)
	fieldValue := v.Elem().FieldByName(fieldName)

	if !fieldValue.IsValid() {
		return fmt.Errorf("field %s does not exist", fieldName)
	}

	if fieldValue.IsNil() {
		return WrapErr(fieldName, NOT_INIT)
	}

	return nil
}
