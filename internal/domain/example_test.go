package domain_test

import (
	"fmt"

	"github.com/StevenYu0124/kinfolio/internal/domain"
)

func ExampleNewMoney() {
	m, err := domain.NewMoney("10.50", domain.SGD)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(m)
	// Output: 10.5 SGD
}

func ExampleMoney_Add() {
	a, _ := domain.NewMoney("10.50", domain.SGD)
	b, _ := domain.NewMoney("2.25", domain.SGD)
	sum, err := a.Add(b)
	fmt.Println(sum, err)

	c, _ := domain.NewMoney("5", domain.HKD)
	_, err = a.Add(c)
	fmt.Println(err)
	// Output:
	// 12.75 SGD <nil>
	// currency mismatch: SGD and HKD
}

func ExampleParseCurrency() {
	c, err := domain.ParseCurrency("SGD")
	fmt.Println(c, err)

	_, err = domain.ParseCurrency("USD")
	fmt.Println(err)
	// Output:
	// SGD <nil>
	// unknown currency "USD"
}
