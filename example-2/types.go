/**
 * @author Jose Nidhin
 */
package main

import (
	"fmt"
	"strings"
	"text/tabwriter"

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
	PointOfServiceConsitionCode            *field.String  `index:"25"`
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

func (fmr *FinancialMessageRequest) PrettyPrint() string {
	var builder strings.Builder
	tw := tabwriter.NewWriter(&builder, 2, 2, 1, ' ', 0)

	fmt.Fprintf(tw, "MTI\t%s\n", fmr.MTI.Value)
	fmt.Fprintf(tw, "PrimaryAccountNumber\t%d\n", fmr.PrimaryAccountNumber.Value)
	fmt.Fprintf(tw, "ProcessingCode\t%s\n", fmr.ProcessingCode.Value)
	fmt.Fprintf(tw, "TransactionAmount\t%d\n", fmr.TransactionAmount.Value)
	fmt.Fprintf(tw, "TransmissionDateTime\t%s\n", fmr.TransmissionDateTime.Value)

	tw.Flush()

	return builder.String()
}

type FinancialMessageResponse struct {
	MTI                                *field.String  `index:"0"`
	PrimaryAccountNumber               *field.String  `index:"2"`
	ProcessingCode                     *field.String  `index:"3"`
	TransactionAmount                  *field.Numeric `index:"4"`
	TransmissionDateTime               *field.String  `index:"7"`
	STAN                               *field.String  `index:"11"`
	LocalTransactionTime               *field.String  `index:"12"`
	CaptureDate                        *field.String  `index:"17"`
	RetrievalReferenceNumber           *field.String  `index:"37"`
	CardAcceptorTerminalIdentification *field.String  `index:"41"`
	TransactionCurrencyCode            *field.String  `index:"49"`
}
