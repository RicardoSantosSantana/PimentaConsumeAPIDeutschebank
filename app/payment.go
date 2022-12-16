package main

func is_payment(bankTransaction FINANCE_BankTransactions) bool {
	return bankTransaction.Amount > 0
}
