/**
 * @author Jose Nidhin
 */
package main

import (
	"io"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/network"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/pkg/errors"
)

const Spec1HeaderSize = 12

var Spec1 *iso8583.MessageSpec = &iso8583.MessageSpec{
	Name: "ISO 8583 ASCII Test Spec 1",
	Fields: map[int]field.Field{
		0: field.NewString(&field.Spec{
			Length:      4,
			Description: "Message Type Indicator",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		1: field.NewBitmap(&field.Spec{
			Description: "Bitmap",
			Enc:         encoding.BytesToASCIIHex,
			Pref:        prefix.Hex.Fixed,
		}),
		2: field.NewNumeric(&field.Spec{
			Length:      19,
			Description: "Primary Account Number",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		3: field.NewString(&field.Spec{
			Length:      6,
			Description: "Processing Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		4: field.NewNumeric(&field.Spec{
			Length:      12,
			Description: "Transaction Amount",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
			Pad:         padding.Left('0'),
		}),
		7: field.NewString(&field.Spec{
			Length:      10,
			Description: "Transmission Date & Time",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		11: field.NewNumeric(&field.Spec{
			Length:      6,
			Description: "Systems Trace Audit Number (STAN)",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		12: field.NewString(&field.Spec{
			Length:      6,
			Description: "Local Transaction Time",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		13: field.NewString(&field.Spec{
			Length:      4,
			Description: "Local Transaction Date",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		15: field.NewString(&field.Spec{
			Length:      4,
			Description: "Settlement Date",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		17: field.NewString(&field.Spec{
			Length:      4,
			Description: "Capture Date",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		25: field.NewString(&field.Spec{
			Length:      2,
			Description: "Point of Service Condition Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		32: field.NewString(&field.Spec{
			Length:      99,
			Description: "Acquiring Institution Identification Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		37: field.NewString(&field.Spec{
			Length:      12,
			Description: "Retrieval Reference Number",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		38: field.NewString(&field.Spec{
			Length:      6,
			Description: "Authorization Identification Response",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		39: field.NewString(&field.Spec{
			Length:      2,
			Description: "Response Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		41: field.NewString(&field.Spec{
			Length:      16,
			Description: "Card Acceptor Terminal Identification",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		43: field.NewString(&field.Spec{
			Length:      40,
			Description: "Card Acceptor Name/Location",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		48: field.NewString(&field.Spec{
			Length:      99,
			Description: "Additional Data - Retailer Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		49: field.NewString(&field.Spec{
			Length:      3,
			Description: "Transaction Currency Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		54: field.NewString(&field.Spec{
			Length:      99,
			Description: "Additional Amounts",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LL,
		}),
		58: field.NewString(&field.Spec{
			Length:      999,
			Description: "Loyalty Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		63: field.NewString(&field.Spec{
			Length:      999,
			Description: "POS Additional Data",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.LLL,
		}),
		70: field.NewString(&field.Spec{
			Length:      3,
			Description: "Network Management Information Code",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
		90: field.NewString(&field.Spec{
			Length:      99,
			Description: "Original Data Elements",
			Enc:         encoding.ASCII,
			Pref:        prefix.ASCII.Fixed,
		}),
	},
}

func MsgLenReader(r io.Reader) (int, error) {
	header := network.NewBinary2BytesHeader()

	_, err := header.ReadFrom(r)
	if err != nil {
		return 0, errors.Wrap(err, "reading header failed")
	}

	return header.Length(), nil
}

func MsgLenWriter(w io.Writer, length int) (int, error) {
	header := network.NewBinary2BytesHeader()
	header.SetLength(length)

	wrote, err := header.WriteTo(w)
	if err != nil {
		return wrote, errors.Wrap(err, "writing header failed")
	}

	return wrote, nil
}
