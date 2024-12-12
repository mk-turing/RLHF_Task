package main

import (
	"fmt"
)

// BankAccount struct encapsulates account number and balance
type BankAccount struct {
	accountNumber string
	balance       float64
}

// NewBankAccount is a constructor function to create a new BankAccount
func NewBankAccount(accountNumber string, initialBalance float64) *BankAccount {
	if initialBalance < 0 {
		panic("Initial balance cannot be negative")
	}
	return &BankAccount{
		accountNumber: accountNumber,
		balance:       initialBalance,
	}
}

// Deposit method adds a specified amount to the balance
func (account *BankAccount) Deposit(amount float64) {
	if amount <= 0 {
		fmt.Println("Deposit amount must be positive")
		return
	}
	account.balance += amount
	fmt.Printf("Deposited $%.2f. New balance: $%.2f\n", amount, account.balance)
}

// Withdraw method subtracts a specified amount from the balance
func (account *BankAccount) Withdraw(amount float64) {
	if amount <= 0 {
		fmt.Println("Withdrawal amount must be positive")
		return
	}
	if amount > account.balance {
		fmt.Println("Insufficient balance")
		return
	}
	account.balance -= amount
	fmt.Printf("Withdrew $%.2f. New balance: $%.2f\n", amount, account.balance)
}

// GetBalance method returns the current balance
func (account *BankAccount) GetBalance() float64 {
	return account.balance
}

// GetAccountNumber method returns the account number
func (account *BankAccount) GetAccountNumber() string {
	return account.accountNumber
}

func main() {
	// Creating a new bank account
	account := NewBankAccount("123456789", 1000.00)

	// Accessing methods to manipulate the data
	account.Deposit(500.00)
	account.Withdraw(300.00)
	fmt.Printf("Current balance: $%.2f\n", account.GetBalance())
	fmt.Printf("Account Number: %s\n", account.GetAccountNumber())
}
