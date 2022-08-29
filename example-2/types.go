/**
 * @author Jose Nidhin
 */
package main

import (
	"fmt"
	"reflect"
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

func (fmr *FinancialMessageRequest) PrettyPrint() string {
	var builder strings.Builder
	tw := tabwriter.NewWriter(&builder, 2, 2, 1, ' ', 0)

	cases := []struct {
		Item   field.Field
		Format string
	}{
		{
			Item:   fmr.MTI,
			Format: "MTI\t%s",
		},
		{
			Item:   fmr.PrimaryAccountNumber,
			Format: "PrimaryAccountNumber\t%d",
		},
		{
			Item:   fmr.ProcessingCode,
			Format: "ProcessingCode\t%s",
		},
		{
			Item:   fmr.TransactionAmount,
			Format: "TransactionAmount\t%d",
		},
		{
			Item:   fmr.TransmissionDateTime,
			Format: "TransmissionDateTime\t%s",
		},
		{
			Item:   fmr.STAN,
			Format: "STAN\t%d",
		},
		{
			Item:   fmr.LocalTransactionTime,
			Format: "LocalTransactionTime\t%s",
		},
		{
			Item:   fmr.LocalTransactionDate,
			Format: "LocalTransactionDate\t%s",
		},
		{
			Item:   fmr.CaptureDate,
			Format: "CaptureDate\t%s",
		},
		{
			Item:   fmr.PointOfServiceConditionCode,
			Format: "PointOfServiceConditionCode\t%s",
		},
		{
			Item:   fmr.AcquiringInstitutionIdentificationCode,
			Format: "AcquiringInstitutionIdentificationCode\t%s",
		},
		{
			Item:   fmr.RetrievalReferenceNumber,
			Format: "RetrievalReferenceNumber\t%s",
		},
		{
			Item:   fmr.CardAcceptorTerminalIdentification,
			Format: "CardAcceptorTerminalIdentification\t%s",
		},
		{
			Item:   fmr.CardAcceptorNameLocation,
			Format: "CardAcceptorNameLocation\t%s",
		},
		{
			Item:   fmr.AdditionalData,
			Format: "AdditionalData\t%s",
		},
		{
			Item:   fmr.LoyaltyData,
			Format: "LoyaltyData\t%s",
		},
		{
			Item:   fmr.POSAdditionalData,
			Format: "POSAdditionalData\t%s",
		},
	}

	for _, c := range cases {
		if c.Item == nil || reflect.ValueOf(c.Item).IsNil() {
			fmt.Fprintln(tw, c.Format, "Field not found")
			continue
		}

		switch item := c.Item.(type) {
		case *field.String:
			fmt.Fprintf(tw, c.Format, item.Value)
			fmt.Fprintln(tw)
			continue
		case *field.Numeric:
			fmt.Fprintf(tw, c.Format, item.Value)
			fmt.Fprintln(tw)
			continue
		default:
			fmt.Fprintln(tw, c.Format, "Unknown field type")
			continue
		}
	}

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
