package main

type Transactions struct {
	OriginIban string `json:"originIban"`
	Amount float32 `json:"amount"`
	CounterPartyName string `json:"counterPartyName"`
	PaymentReference string  `json:"paymentReference"` 
	BookingDate string  `json:"bookingDate"`  
	CurrencyCode string  `json:"currencyCode"`  
	TransactionCode string  `json:"transactionCode"`  
	ExternalBankTransactionDomainCode string  `json:"externalBankTransactionDomainCode"`  
	ExternalBankTransactionFamilyCode string  `json:"externalBankTransactionFamilyCode"`  
	ExternalBankTransactionSubFamilyCode string  `json:"externalBankTransactionSubFamilyCode"`  
	MandateReference string  `json:"mandateReference"`  
	CreditorId string  `json:"creditorId"`  
	E2eReference string  `json:"e2eReference"`  
	PaymentIdentification string  `json:"paymentIdentification"`  
	ValueDate string  `json:"valueDate"`  
	Id string  `json:"id"`  
}

type CashAccount struct {
	TotalItems int `json:"totalItems"`
	Offset int `json:"offset"`
	Limit int `json:"limit"`
	Transactions []Transactions `json:"transactions"`
}
 