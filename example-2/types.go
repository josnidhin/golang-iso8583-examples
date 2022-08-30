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

type ReversalMessageRequest struct {
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
	CardAcceptorNameLocation               *field.String  `index:"43"`
	AdditionalData                         *field.String  `index:"48"`
	TransactionCurrencyCode                *field.String  `index:"49"`
	AdditionalAmounts                      *field.String  `index:"54"`
	POSAdditionalData                      *field.String  `index:"63"`
	OriginalDataElements                   *field.String  `index:"90"`
}

func (rmr *ReversalMessageRequest) PrettyPrint() string {
	var builder strings.Builder
	tw := tabwriter.NewWriter(&builder, 2, 2, 1, ' ', 0)

	cases := []struct {
		Item   field.Field
		Format string
	}{
		{
			Item:   rmr.MTI,
			Format: "MTI\t%s",
		},
		{
			Item:   rmr.PrimaryAccountNumber,
			Format: "PrimaryAccountNumber\t%d",
		},
		{
			Item:   rmr.ProcessingCode,
			Format: "ProcessingCode\t%s",
		},
		{
			Item:   rmr.TransactionAmount,
			Format: "TransactionAmount\t%d",
		},
		{
			Item:   rmr.TransmissionDateTime,
			Format: "TransmissionDateTime\t%s",
		},
		{
			Item:   rmr.STAN,
			Format: "STAN\t%d",
		},
		{
			Item:   rmr.LocalTransactionTime,
			Format: "LocalTransactionTime\t%s",
		},
		{
			Item:   rmr.LocalTransactionDate,
			Format: "LocalTransactionDate\t%s",
		},
		{
			Item:   rmr.SettlementDate,
			Format: "SettlementDate\t%s",
		},
		{
			Item:   rmr.CaptureDate,
			Format: "CaptureDate\t%s",
		},
		{
			Item:   rmr.PointOfServiceConditionCode,
			Format: "PointOfServiceConditionCode\t%s",
		},
		{
			Item:   rmr.AcquiringInstitutionIdentificationCode,
			Format: "AcquiringInstitutionIdentificationCode\t%s",
		},
		{
			Item:   rmr.RetrievalReferenceNumber,
			Format: "RetrievalReferenceNumber\t%s",
		},
		{
			Item:   rmr.AuthorizationIdentificationResponse,
			Format: "AuthorizationIdentificationResponse\t%s",
		},
		{
			Item:   rmr.ResponseCode,
			Format: "ResponseCode\t%s",
		},
		{
			Item:   rmr.CardAcceptorTerminalIdentification,
			Format: "CardAcceptorTerminalIdentification\t%s",
		},
		{
			Item:   rmr.CardAcceptorNameLocation,
			Format: "CardAcceptorNameLocation\t%s",
		},
		{
			Item:   rmr.AdditionalData,
			Format: "AdditionalData\t%s",
		},
		{
			Item:   rmr.TransactionCurrencyCode,
			Format: "TransactionCurrencyCode\t%s",
		},
		{
			Item:   rmr.AdditionalAmounts,
			Format: "AdditionalAmounts\t%s",
		},
		{
			Item:   rmr.POSAdditionalData,
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

type ReversalMessageResponse struct {
	MTI                                    *field.String  `index:"0"`
	PrimaryAccountNumber                   *field.Numeric `index:"2"`
	ProcessingCode                         *field.String  `index:"3"`
	TransactionAmount                      *field.Numeric `index:"4"`
	TransmissionDateTime                   *field.String  `index:"7"`
	STAN                                   *field.Numeric `index:"11"`
	SettlementDate                         *field.String  `index:"15"`
	CaptureDate                            *field.String  `index:"17"`
	PointOfServiceConditionCode            *field.String  `index:"25"`
	AcquiringInstitutionIdentificationCode *field.String  `index:"32"`
	RetrievalReferenceNumber               *field.String  `index:"37"`
	ResponseCode                           *field.String  `index:"39"`
	CardAcceptorTerminalIdentification     *field.String  `index:"41"`
	TransactionCurrencyCode                *field.String  `index:"49"`
	AdditionalAmounts                      *field.String  `index:"54"`
	POSAdditionalData                      *field.String  `index:"63"`
	OriginalDataElements                   *field.String  `index:"90"`
}

type EchoMessageRequest struct {
	MTI                              *field.String  `index:"0"`
	Bitmap                           *field.Bitmap  `index:"1"`
	TransmissionDateTime             *field.String  `index:"7"`
	STAN                             *field.Numeric `index:"11"`
	SettlementDate                   *field.String  `index:"15"`
	AdditionalData                   *field.String  `index:"48"`
	NetworkManagementInformationCode *field.String  `index:"70"`
}

func (emr *EchoMessageRequest) PrettyPrint() string {
	var builder strings.Builder
	tw := tabwriter.NewWriter(&builder, 2, 2, 1, ' ', 0)

	cases := []struct {
		Item   field.Field
		Format string
	}{
		{
			Item:   emr.MTI,
			Format: "MTI\t%s",
		},
		{
			Item:   emr.Bitmap,
			Format: "Bitmap\t%s",
		},
		{
			Item:   emr.TransmissionDateTime,
			Format: "TransmissionDateTime\t%s",
		},
		{
			Item:   emr.STAN,
			Format: "STAN\t%d",
		},
		{
			Item:   emr.SettlementDate,
			Format: "SettlementDate\t%s",
		},
		{
			Item:   emr.AdditionalData,
			Format: "AdditionalData\t%s",
		},
		{
			Item:   emr.NetworkManagementInformationCode,
			Format: "NetworkManagementInformationCode\t%s",
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
		case *field.Bitmap:
			strVal, err := item.String()
			if err != nil {
				strVal = err.Error()
			}
			fmt.Fprintf(tw, c.Format, strVal)
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

type EchoResponse struct {
	MTI                              *field.String  `index:"0"`
	TransmissionDateTime             *field.String  `index:"7"`
	STAN                             *field.Numeric `index:"11"`
	SettlementDate                   *field.String  `index:"15"`
	ResponseCode                     *field.String  `index:"39"`
	NetworkManagementInformationCode *field.String  `index:"70"`
}
