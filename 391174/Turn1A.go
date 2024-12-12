package main

import (
	"errors"
	"fmt"
)

// BankAccount defines a struct with an unexported balance field
type BankAccount struct {
	owner   string
	balance float64 // unexported field
}

// NewBankAccount is a constructor function to create a new BankAccount
func NewBankAccount(owner string) *BankAccount {
	return &BankAccount{
		owner:   owner,
		balance: 0, // starting balance is 0
	}
}

// Deposit allows money to be added to the account balance
func (account *BankAccount) Deposit(amount float64) {
	if amount > 0 {
		account.balance += amount
	} else {
		fmt.Println("Deposit amount must be positive")
	}
}

// Withdraw allows money to be subtracted from the account balance
func (account *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if amount > account.balance {
		return errors.New("insufficient funds")
	}
	account.balance -= amount
	return nil
}

// GetBalance provides access to the current balance
func (account *BankAccount) GetBalance() float64 {
	return account.balance
}

func main() {
	// Create a new bank account
	account := NewBankAccount("John Doe")

	// Deposit some money
	account.Deposit(100)
	fmt.Println("Current Balance:", account.GetBalance()) // Current Balance: 100

	// Attempt to withdraw too much money
	err := account.Withdraw(150)
	if err != nil {
		fmt.Println("Withdrawal Error:", err) // Withdrawal Error: insufficient funds
	}

	// Withdraw money successfully
	err = account.Withdraw(50)
	if err != nil {
		fmt.Println("Withdrawal Error:", err)
	} else {
		fmt.Println("Withdraw successful, new balance:", account.GetBalance()) // new balance: 50
	}
}
