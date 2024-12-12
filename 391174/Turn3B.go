package main

import (
	"fmt"
)

// Account interface defines common methods for all accounts
type Account interface {
	Deposit(float64)
	Withdraw(float64)
	GetBalance() float64
	GetAccountNumber() string
}

// BankAccount struct implements the Account interface
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
	fmt.Printf("Deposited $%.2f to account %s. New balance: $%.2f\n", amount, account.accountNumber, account.balance)
}

// Withdraw method subtracts a specified amount from the balance
func (account *BankAccount) Withdraw(amount float64) {
	if amount <= 0 {
		fmt.Println("Withdrawal amount must be positive")
		return
	}
	if amount > account.balance {
		fmt.Println("Insufficient balance in account", account.accountNumber)
		return
	}
	account.balance -= amount
	fmt.Printf("Withdrew $%.2f from account %s. New balance: $%.2f\n", amount, account.accountNumber, account.balance)
}

// GetBalance method returns the current balance
func (account *BankAccount) GetBalance() float64 {
	return account.balance
}

// GetAccountNumber method returns the account number
func (account *BankAccount) GetAccountNumber() string {
	return account.accountNumber
}

// SavingsAccount struct implements the Account interface with a different interest logic
type SavingsAccount struct {
	accountNumber string
	balance       float64
	interestRate  float64
}

// NewSavingsAccount is a constructor function to create a new SavingsAccount
func NewSavingsAccount(accountNumber string, initialBalance float64, interestRate float64) *SavingsAccount {
	if initialBalance < 0 {
		panic("Initial balance cannot be negative")
	}
	if interestRate < 0 {
		panic("Interest rate cannot be negative")
	}
	return &SavingsAccount{
		accountNumber: accountNumber,
		balance:       initialBalance,
		interestRate:  interestRate,
	}
}

// Deposit method adds a specified amount to the balance
func (account *SavingsAccount) Deposit(amount float64) {
	if amount <= 0 {
		fmt.Println("Deposit amount must be positive")
		return
	}
	account.balance += amount
	fmt.Printf("Deposited $%.2f to account %s. New balance: $%.2f\n", amount, account.accountNumber, account.balance)
}

// Withdraw method subtracts a specified amount from the balance
func (account *SavingsAccount) Withdraw(amount float64) {
	if amount <= 0 {
		fmt.Println("Withdrawal amount must be positive")
		return
	}
	if amount > account.balance {
		fmt.Println("Insufficient balance in account", account.accountNumber)
		return
	}
	account.balance -= amount
	fmt.Printf("Withdrew $%.2f from account %s. New balance: $%.2f\n", amount, account.accountNumber, account.balance)
}

// GetBalance method returns the current balance
func (account *SavingsAccount) GetBalance() float64 {
	return account.balance
}

// GetAccountNumber method returns the account number
func (account *SavingsAccount) GetAccountNumber() string {
	return account.accountNumber
}

// Customer struct encapsulates customer details and multiple accounts
type Customer struct {
	name     string
	accounts map[string]Account
}

// NewCustomer is a constructor function to create a new Customer
func NewCustomer(name string) *Customer {
	return &Customer{
		name:     name,
		accounts: make(map[string]Account),
	}
}

// OpenAccount method allows a customer to open a new bank account
func (c *Customer) OpenAccount(accountType string, accountNumber string, initialBalance float64, interestRate float64) {
	if _, exists := c.accounts[accountNumber]; exists {
		fmt.Println("Account already exists:", accountNumber)
		return
	}
	switch accountType {
	case "checking":
		c.accounts[accountNumber] = NewBankAccount(accountNumber, initialBalance)
	case "savings":
		c.accounts[accountNumber] = NewSavingsAccount(accountNumber, initialBalance, interestRate)
	default:
		fmt.Println("Invalid account type")
		return
	}
	fmt.Printf("Customer %s opened account %s with balance: $%.2f\n", c.name, accountNumber, initialBalance)
}

// DepositToAccount method allows depositing money into a specific account
func (c *Customer) DepositToAccount(accountNumber string, amount float64) {
	if account, exists := c.accounts[accountNumber]; exists {
		account.Deposit(amount)
	} else {
		fmt.Println("Account not found:", accountNumber)
	}
}

// WithdrawFromAccount method allows withdrawing money from a specific account
func (c *Customer) WithdrawFromAccount(accountNumber string, amount float64) {
	if account, exists := c.accounts[accountNumber]; exists {
		account.Withdraw(amount)
	} else {
		fmt.Println("Account not found:", accountNumber)
	}
}

// GetAccountBalance method returns the balance of a specific account
func (c *Customer) GetAccountBalance(accountNumber string) float64 {
	if account, exists := c.accounts[accountNumber]; exists {
		return account.GetBalance()
	}
	fmt.Println("Account not found:", accountNumber)
	return 0
}

func main() {
	// Creating a new customer
	customer := NewCustomer("Alice")

	// Customer opens accounts
	customer.OpenAccount("checking", "123456789", 1000.00, 0)
	customer.OpenAccount("savings", "987654321", 500.00, 0.02)

	// Customer deposits and withdraws from accounts
	customer.DepositToAccount("123456789", 200.00)
	customer.WithdrawFromAccount("987654321", 100.00)
	customer.WithdrawFromAccount("123456789", 300.00)

	// Checking balances
	fmt.Printf("Balance in account 123456789: $%.2f\n", customer.GetAccountBalance("123456789"))
	fmt.Printf("Balance in account 987654321: $%.2f\n", customer.GetAccountBalance("987654321"))
}
