/**
 * @author Jose Nidhin
 */
package main

import (
	"github.com/moov-io/iso8583/field"
)

type FinancialMessageRequest struct {
	MTI                                    *field.String  `index:"0"`
	PrimaryAccountNumber                   *field.Numeric `index:"2"`
	ProcessingCode                         *field.String  `index:"3"`
	TransactionAmount                      *field.Numeric `index:"4"`
	TransmissionDateTime                   *field.String  `index:"7"`
	STAN                                   *field.Numeric `index:"11"`
	LocalTransactionTime                   *field.String  `index:"12"`
	LocalTransactionDate                   *field.String  `index:"13"`
	CaptureDate                            *field.String  `index:"17"`
	PointOfServiceConditionCode            *field.String  `index:"25"`
	AcquiringInstitutionIdentificationCode *field.String  `index:"32"`
	RetrievalReferenceNumber               *field.String  `index:"37"`
	CardAcceptorTerminalIdentification     *field.String  `index:"41"`
	CardAcceptorNameLocation               *field.String  `index:"43"`
	AdditionalData                         *field.String  `index:"48"`
	TransactionCurrencyCode                *field.String  `index:"49"`
	AdditionalAmounts                      *field.String  `index:"54"`
	LoyaltyData                            *field.String  `index:"58"`
	POSAdditionalData                      *field.String  `index:"63"`
}

type FinancialMessageResponse struct {
	MTI                                    *field.String  `index:"0"`
	PrimaryAccountNumber                   *field.Numeric `index:"2"`
	ProcessingCode                         *field.String  `index:"3"`
	TransactionAmount                      *field.Numeric `index:"4"`
	TransmissionDateTime                   *field.String  `index:"7"`
	STAN                                   *field.Numeric `index:"11"`
	LocalTransactionTime                   *field.String  `index:"12"`
	LocalTransactionDate                   *field.String  `index:"13"`
	SettlementDate                         *field.String  `index:"15"`
	CaptureDate                            *field.String  `index:"17"`
	PointOfServiceConditionCode            *field.String  `index:"25"`
	AcquiringInstitutionIdentificationCode *field.String  `index:"32"`
	RetrievalReferenceNumber               *field.String  `index:"37"`
	AuthorizationIdentificationResponse    *field.String  `index:"38"`
	ResponseCode                           *field.String  `index:"39"`
	CardAcceptorTerminalIdentification     *field.String  `index:"41"`
	AdditionalData                         *field.String  `index:"48"`
	TransactionCurrencyCode                *field.String  `index:"49"`
	AdditionalAmounts                      *field.String  `index:"54"`
	LoyaltyData                            *field.String  `index:"58"`
	POSAdditionalData                      *field.String  `index:"63"`
}
