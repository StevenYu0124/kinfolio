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
	a, err := domain.NewMoney("10.50", domain.SGD)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, err := domain.NewMoney("2.25", domain.SGD)
	if err != nil {
		fmt.Println(err)
		return
	}

	sum, err := a.Add(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sum)

	c, err := domain.NewMoney("5", domain.HKD)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = a.Add(c)
	fmt.Println(err)

	// Output:
	// 12.75 SGD
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
