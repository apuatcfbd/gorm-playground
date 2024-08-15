package main

import (
	"fmt"
	"log"
)

func main() {
	// create first record
	m := DynamicKeyValue{
		StrFieldBeforeValue: "I'm before value field",
		Value: KvValue{
			Val: 1,
		},
		ValueType:          KindInt,
		StrFieldAfterValue: "I'm after value field, I'll be lost in serializer!",
	}

	tx := DB.Create(&m)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	fmt.Printf("%#v -> %#v \n", "Successfully inserted 1st record", m)

	// retrieve
	mm := new(DynamicKeyValue)
	tx2 := DB.First(mm, 1)
	if tx2.Error != nil {
		log.Fatal(tx2.Error)
	}

	fmt.Printf("%#v -> %#v \n", "Successfully retrieved 1st record", mm)
}
