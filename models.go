package main

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strconv"
)

type DynamicKeyValue struct {
	gorm.Model

	StrFieldBeforeValue string  `json:"str_field_bv" gorm:"type:varchar(255); not null;"`
	Value               KvValue `json:"value" gorm:"type:text;"`
	//ValueType Kind   `json:"kind" gorm:"type:enum('str', 'bool', 'int'); not null;"`
	ValueType          Kind   `json:"kind" gorm:"type:varchar(255); not null;"`
	StrFieldAfterValue string `json:"str_field_av" gorm:"type:varchar(255); not null;"`
}

// Arrangement to have enum in "ValueType" start //

type Kind string

func (n *Kind) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	switch val := dbValue.(type) {
	case []byte:
		*n = Kind(val)
	case string:
		*n = Kind(val)
	default:
		return fmt.Errorf("invalid kind: %v", dbValue)
	}
	return nil
}

func (*Kind) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return string(fieldValue.(Kind)), nil
}

const (
	KindStr  Kind = "str"
	KindBool Kind = "bool"
	KindInt  Kind = "int"
)

// Arrangement to have enum in "ValueType" end //

type KvValue struct {
	Val any `json:"val"`
}

// Scan implements schema.SerializerInterface (db out)
func (v *KvValue) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {

	fmt.Printf("vvv -> StrFieldBeforeValue's Value: %#v <<<<[OK] \n", dst.FieldByName("StrFieldBeforeValue"))
	fmt.Printf("xxx -> StrFieldAfterValue's Value: %#v <<<<[THIS IS THE ISSUE, THIS FIELD IS EMPTY, ALL THE FIELDS AFTER current field (Value) have ZERO value !!!] \n", dst.FieldByName("StrFieldAfterValue"))
	fmt.Printf("%#v \n", "I want to return the value as 'int' if 'ValueType' == 'int', bool -> bool, etc. Now this is not possible as other fields are empty")
	// find out the type
	// decode value base on the type
	switch val := dbValue.(type) {
	case []byte:
		v.Val = string(val)
	case string:
		v.Val = val
	default:
		err = fmt.Errorf("unsupported data type (%T) in value", val)
	}

	return
}

// Value implements schema.SerializerInterface (db in)
func (*KvValue) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	sv, ok := fieldValue.(KvValue)
	if !ok {
		return nil, errors.New("setting value field should be type of SettingValue")
	}

	var r any

	switch v := sv.Val.(type) {
	case string, bool:
		r = v
	case int:
		r = strconv.Itoa(v)
	// case ...
	default:
		return nil, fmt.Errorf("invalid value type %T in field %s", fieldValue, field.Name)
	}

	return r, nil
}
