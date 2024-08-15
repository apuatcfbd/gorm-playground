package main

import (
	"reflect"
	"testing"
)

// GORM_REPO: https://github.com/go-gorm/gorm.git
// GORM_BRANCH: master
// TEST_DRIVERS: sqlite, mysql, postgres, sqlserver

//	func TestGORM(t *testing.T) {
//		user := User{Name: "jinzhu"}
//
//		DB.Create(&user)
//
//		var result User
//		if err := DB.First(&result, user.ID).Error; err != nil {
//			t.Errorf("Failed, got error: %v", err)
//		}
//	}
func TestSerializer(t *testing.T) {
	m := DynamicKeyValue{
		StrFieldBeforeValue: "I'm before value field",
		Value: KvValue{
			Val: 1,
		},
		ValueType:          KindStr,
		StrFieldAfterValue: "I'm after value field, I'll be lost in serializer!",
	}

	tx := DB.Create(&m)
	if tx.Error != nil {
		t.Errorf("%c", tx.Error)
	}

	_, ok := m.Value.Val.(string)
	if !ok {
		// expecting string as value type is 'str' but that conversion never happened in the serializer. :(
		// checkout: models.go:59 for the issue
		t.Errorf("Expected Value type 'string' got '%s'. Checkout models.go:59 for the issue", reflect.TypeOf(m.Value.Val).Name())
	}
}
