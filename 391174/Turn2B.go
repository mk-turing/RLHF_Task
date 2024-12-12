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

// Bank struct manages a collection of BankAccount structs
type Bank struct {
	accounts map[string]*BankAccount
}

// NewBank is a constructor function to create a new Bank
func NewBank() *Bank {
	return &Bank{
		accounts: make(map[string]*BankAccount),
	}
}

// OpenAccount method adds a new account to the bank
func (bank *Bank) OpenAccount(accountNumber string, initialBalance float64) {
	if _, exists := bank.accounts[accountNumber]; exists {
		fmt.Println("Account already exists")
		return
	}
	bank.accounts[accountNumber] = NewBankAccount(accountNumber, initialBalance)
	fmt.Printf("Account %s opened with balance $%.2f\n", accountNumber, initialBalance)
}

// Transfer method transfers money between two accounts
func (bank *Bank) Transfer(fromAccountNumber string, toAccountNumber string, amount float64) {
	fromAccount, existsFrom := bank.accounts[fromAccountNumber]
	toAccount, existsTo := bank.accounts[toAccountNumber]

	if !existsFrom {
		fmt.Println("From account does not exist")
		return
	}
	if !existsTo {
		fmt.Println("To account does not exist")
		return
	}

	fromAccount.Withdraw(amount)
	toAccount.Deposit(amount)
}

func main() {
	// Creating a new bank
	bank := NewBank()

	// Opening accounts
	bank.OpenAccount("123456789", 1000.00)
	bank.OpenAccount("987654321", 500.00)

	// Transferring money between accounts
	bank.Transfer("123456789", "987654321", 200.00)

	// Accessing account balances
	fmt.Printf("Balance of account 123456789: $%.2f\n", bank.accounts["123456789"].GetBalance())
	fmt.Printf("Balance of account 987654321: $%.2f\n", bank.accounts["987654321"].GetBalance())
}
