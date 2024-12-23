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

// Customer struct encapsulates customer details and multiple bank accounts
type Customer struct {
	name     string
	accounts map[string]*BankAccount
}

// NewCustomer is a constructor function to create a new Customer
func NewCustomer(name string) *Customer {
	return &Customer{
		name:     name,
		accounts: make(map[string]*BankAccount),
	}
}

// OpenAccount method allows a customer to open a new bank account
func (c *Customer) OpenAccount(accountNumber string, initialBalance float64) {
	if _, exists := c.accounts[accountNumber]; exists {
		fmt.Println("Account already exists:", accountNumber)
		return
	}
	c.accounts[accountNumber] = NewBankAccount(accountNumber, initialBalance)
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
	customer.OpenAccount("123456789", 1000.00)
	customer.OpenAccount("987654321", 500.00)

	// Customer deposits and withdraws from accounts
	customer.DepositToAccount("123456789", 200.00)
	customer.WithdrawFromAccount("987654321", 100.00)
	customer.WithdrawFromAccount("123456789", 300.00)

	// Checking balances
	fmt.Printf("Balance in account 123456789: $%.2f\n", customer.GetAccountBalance("123456789"))
	fmt.Printf("Balance in account 987654321: $%.2f\n", customer.GetAccountBalance("987654321"))
}
