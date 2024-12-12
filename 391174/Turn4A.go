package main

import (
	"fmt"
	"sync"
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
	mu            sync.Mutex // Mutex to protect access to balance
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

// Deposit method adds a specified amount to the balance safely
func (account *BankAccount) Deposit(amount float64) {
	account.mu.Lock()         // Lock the mutex
	defer account.mu.Unlock() // Unlock when function exits

	if amount <= 0 {
		fmt.Println("Deposit amount must be positive")
		return
	}
	account.balance += amount
	fmt.Printf("Deposited $%.2f to account %s. New balance: $%.2f\n", amount, account.accountNumber, account.balance)
}

// Withdraw method subtracts a specified amount from the balance safely
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

// GetBalance method returns the current balance safely
func (account *BankAccount) GetBalance() float64 {
	account.mu.Lock()
	defer account.mu.Unlock()
	return account.balance
}

// GetAccountNumber method returns the account number
func (account *BankAccount) GetAccountNumber() string {
	return account.accountNumber
}

// SavingsAccount struct that also implements the Account interface
type SavingsAccount struct {
	accountNumber string
	balance       float64
	interestRate  float64    // in percentage
	mu            sync.Mutex // Mutex to protect access to balance
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
	account.mu.Lock()         // Lock the mutex
	defer account.mu.Unlock() // Unlock when function exits

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
	return account.accountNumber
}

// Customer struct encapsulates customer details and multiple accounts
type Customer struct {
	name     string
	accounts map[string]Account // Changed to Account interface
	mu       sync.Mutex         // Mutex to protect access to accounts map
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
	c.mu.Lock()         // Lock the mutex
	defer c.mu.Unlock() // Unlock when function exits

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
	c.mu.Lock()         // Lock the mutex
	defer c.mu.Unlock() // Unlock when function exits

	if account, exists := c.accounts[accountNumber]; exists {
		account.Deposit(amount)
	} else {
		fmt.Println("Account not found:", accountNumber)
	}
}

// WithdrawFromAccount method allows withdrawing money from a specific account
func (c *Customer) WithdrawFromAccount(accountNumber string, amount float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if account, exists := c.accounts[accountNumber]; exists {
		account.Withdraw(amount)
	} else {
		fmt.Println("Account not found:", accountNumber)
	}
}

// GetAccountBalance method returns the balance of a specific account
func (c *Customer) GetAccountBalance(accountNumber string) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	if account, exists := c.accounts[accountNumber]; exists {
		return account.GetBalance()
	}
	fmt.Println("Account not found:", accountNumber)
	return 0
}

func main() {
	// Creating a new customer
	customer := NewCustomer("Alice")

	// Open the accounts concurrently
	go customer.OpenAccount("123456789", 1000.00, "bank", 0)
	go customer.OpenAccount("987654321", 500.00, "savings", 5.0)

	// Allow time for accounts to be created before performing operations
	// In a production application, avoid sleep and use synchronization
	// This is only for demonstration purposes
	// (You might consider using sync.WaitGroup for better control)
	// For illustration, we'll execute a sleep
	fmt.Println("Waiting for account setup...")
	// Sleep for a short duration to allow goroutines some time to finish their execution
	// Please note: in actual code, never include sleeps for synchronization like this.
	// This is a simple way to illustrate concurrency without proper management.
	_, _ = fmt.Scanln()

	// Customer performs deposits and withdrawals concurrently
	go customer.DepositToAccount("123456789", 200.00)
	go customer.WithdrawFromAccount("987654321", 100.00)
	go customer.WithdrawFromAccount("123456789", 300.00)

	// Sleep for a moment to see operations in action
	fmt.Println("Waiting for transactions to complete...")
	_, _ = fmt.Scanln()

	// Checking balances
	fmt.Printf("Balance in account 123456789: $%.2f\n", customer.GetAccountBalance("123456789"))
	fmt.Printf("Balance in account 987654321: $%.2f\n", customer.GetAccountBalance("987654321"))
}
