package main

import (
	"fmt"
	"sync"
	"time"
)

// Account interface defines methods that any type of account should implement
type Account interface {
	Deposit(amount float64)
	Withdraw(amount float64)
	GetBalance() float64
	GetAccountNumber() string
}

// BankAccount struct encapsulates account number and balance
type BankAccount struct {
	accountNumber string
	balance       float64
	mu            sync.Mutex
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
	account.mu.Lock()
	defer account.mu.Unlock()
	if amount <= 0 {
		fmt.Println("Deposit amount must be positive")
		return
	}
	account.balance += amount
	fmt.Printf("Deposited $%.2f to account %s. New balance: $%.2f\n", amount, account.accountNumber, account.balance)
}

// Withdraw method subtracts a specified amount from the balance
func (account *BankAccount) Withdraw(amount float64) {
	account.mu.Lock()
	defer account.mu.Unlock()
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
	account.mu.Lock()
	defer account.mu.Unlock()
	return account.balance
}

// GetAccountNumber method returns the account number
func (account *BankAccount) GetAccountNumber() string {
	account.mu.Lock()
	defer account.mu.Unlock()
	return account.accountNumber
}

// SavingsAccount struct that also implements the Account interface
type SavingsAccount struct {
	accountNumber string
	balance       float64
	interestRate  float64 // in percentage
	mu            sync.Mutex
}

// NewSavingsAccount is a constructor function to create a new SavingsAccount
func NewSavingsAccount(accountNumber string, initialBalance float64, interestRate float64) *SavingsAccount {
	if initialBalance < 0 {
		panic("Initial balance cannot be negative")
	}
	return &SavingsAccount{
		accountNumber: accountNumber,
		balance:       initialBalance,
		interestRate:  interestRate,
	}
}

// Deposit method for SavingsAccount
func (account *SavingsAccount) Deposit(amount float64) {
	account.mu.Lock()
	defer account.mu.Unlock()
	if amount <= 0 {
		fmt.Println("Deposit amount must be positive")
		return
	}
	account.balance += amount
	fmt.Printf("Deposited $%.2f to savings account %s. New balance: $%.2f\n", amount, account.accountNumber, account.balance)
}

// Withdraw method for SavingsAccount
func (account *SavingsAccount) Withdraw(amount float64) {
	account.mu.Lock()
	defer account.mu.Unlock()
	if amount <= 0 {
		fmt.Println("Withdrawal amount must be positive")
		return
	}
	if amount > account.balance {
		fmt.Println("Insufficient balance in savings account", account.accountNumber)
		return
	}
	account.balance -= amount
	fmt.Printf("Withdrew $%.2f from savings account %s. New balance: $%.2f\n", amount, account.accountNumber, account.balance)
}

// GetBalance method for SavingsAccount
func (account *SavingsAccount) GetBalance() float64 {
	account.mu.Lock()
	defer account.mu.Unlock()
	return account.balance
}

// GetAccountNumber method for SavingsAccount
func (account *SavingsAccount) GetAccountNumber() string {
	account.mu.Lock()
	defer account.mu.Unlock()
	return account.accountNumber
}

// Customer struct encapsulates customer details and multiple accounts
type Customer struct {
	name     string
	accounts map[string]Account // Changed to Account interface
}

// NewCustomer is a constructor function to create a new Customer
func NewCustomer(name string) *Customer {
	return &Customer{
		name:     name,
		accounts: make(map[string]Account),
	}
}

// OpenAccount method allows a customer to open a new account
func (c *Customer) OpenAccount(accountNumber string, initialBalance float64, accountType string, interestRate float64) {
	if _, exists := c.accounts[accountNumber]; exists {
		fmt.Println("Account already exists:", accountNumber)
		return
	}

	var account Account
	if accountType == "savings" {
		account = NewSavingsAccount(accountNumber, initialBalance, interestRate)
	} else {
		account = NewBankAccount(accountNumber, initialBalance)
	}
	c.accounts[accountNumber] = account
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

func depositAndWithdrawConcurrently(customer *Customer, accountNumber string, numDeposits, numWithdrawals int, depositAmount, withdrawalAmount float64) {
	for i := 0; i < numDeposits; i++ {
		customer.DepositToAccount(accountNumber, depositAmount)
	}
	for i := 0; i < numWithdrawals; i++ {
		customer.WithdrawFromAccount(accountNumber, withdrawalAmount)
	}
}

func main() {
	// Creating a new customer
	customer := NewCustomer("Alice")

	// Customer opens accounts
	customer.OpenAccount("123456789", 1000.00, "bank", 0)
	customer.OpenAccount("987654321", 500.00, "savings", 5.0)

	// Perform concurrent deposits and withdrawals
	go depositAndWithdrawConcurrently(customer, "123456789", 50, 30, 50.0, 30.0)
	go depositAndWithdrawConcurrently(customer, "123456789", 20, 40, 20.0, 50.0)
	go depositAndWithdrawConcurrently(customer, "987654321", 40, 20, 40.0, 20.0)

	time.Sleep(2 * time.Second) // Allow goroutines to complete

	// Checking balances
	fmt.Printf("Balance in account 123456789: $%.2f\n", customer.GetAccountBalance("123456789"))
	fmt.Printf("Balance in account 987654321: $%.2f\n", customer.GetAccountBalance("987654321"))
}
