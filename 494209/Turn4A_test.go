package _94209

import (
	"fmt"
	"sort"
	"testing"
)

type Transaction struct {
	CustomerID string
	ProductID  string
	Category   string
	Amount     float64
}

func summarizeTransactions(transactions []Transaction) (totalRevenue, avgTransaction float64, count int, topCustomers map[string]float64, mostPopularCategory string) {
	topCustomers = make(map[string]float64)
	mostSales := 0.0

	for _, trans := range transactions {
		totalRevenue += trans.Amount
		count++
		avgTransaction += trans.Amount
		topCustomers[trans.CustomerID] += trans.Amount

		categorySales := mostPopularCategory
		categorySalesAmount := mostSales

		if trans.Category != "" {
			categorySales = trans.Category
			categorySalesAmount += trans.Amount
		}

		if categorySalesAmount > mostSales {
			mostPopularCategory = categorySales
			mostSales = categorySalesAmount
		}
	}

	if count > 0 {
		avgTransaction = totalRevenue / float64(count)
	}

	// Sort top customers by spend and limit to top 10
	var topCustomerSlices []struct {
		CustomerID string
		Spend      float64
	}
	for customer, spend := range topCustomers {
		topCustomerSlices = append(topCustomerSlices, struct {
			CustomerID string
			Spend      float64
		}{CustomerID: customer, Spend: spend})
	}

	sort.Slice(topCustomerSlices, func(i, j int) bool {
		return topCustomerSlices[i].Spend > topCustomerSlices[j].Spend
	})

	topCustomers = make(map[string]float64)
	for i := 0; i < min(10, len(topCustomerSlices)); i++ {
		topCustomers[topCustomerSlices[i].CustomerID] = topCustomerSlices[i].Spend
	}

	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func BenchmarkSummarizeTransactions(b *testing.B) {
	const numTransactions = 1000000
	transactions := make([]Transaction, numTransactions)

	for i := range transactions {
		transactions[i].CustomerID = fmt.Sprintf("customer%06d", i)
		transactions[i].ProductID = fmt.Sprintf("product%06d", i)
		transactions[i].Category = fmt.Sprintf("category%06d", i%10)
		transactions[i].Amount = float64(i) * 10.0
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		summarizeTransactions(transactions)
	}
}
