/**
 * @author Jose Nidhin
 */
package main

import (
	"fmt"

	"github.com/moov-io/iso8583"
)

const HEADER_SIZE = 12

var rawMessages []string = []string{
	"ISO02110005502007238800008808000101234567890000000000000010000092618372419060118510009260926000000005928MON50EDIOX     N484",
	"ISO0211000550420723880000E80800010123456789000000000000001000009261837241906011851000926092600000000592812345600MON50EDIOX     N484",
	"ISO0211000550800822000000000000004000000000000000821083216015795301",
}

func main() {
	for _, rawMsg := range rawMessages {
		fmt.Printf("Raw Message = %s\n", rawMsg)

		msg := iso8583.NewMessage(Spec1)
		msg.Unpack([]byte(rawMsg[HEADER_SIZE:]))

		// for pos := 0; pos < 128; pos++ {
		// 	value, err := msg.GetString(pos)

		// 	if err != nil {
		// 		continue
		// 	}

		// 	if value == "" {
		// 		continue
		// 	}

		// 	field := msg.GetField(pos)

		// 	fmt.Printf("%3d\t%s\t%s\n", pos, field.Spec().Description, value)
		// }

		mti, err := msg.GetMTI()
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}

		switch mti {
		case "0200":
			financeMsgHandler(msg)
			break
		case "0420":
			reverseMsgHandler(msg)
			break
		case "0800":
			echoMsgHandler(msg)
			break
		default:
			fmt.Println("Unknown message type")
		}

		fmt.Println()
	}
}

func financeMsgHandler(msg *iso8583.Message) {
	req := FinancialMessageRequest{}

	err := msg.Unmarshal(&req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(req.PrettyPrint())
}

func reverseMsgHandler(msg *iso8583.Message) {
	req := ReversalMessageRequest{}

	err := msg.Unmarshal(&req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(req.PrettyPrint())
}

func echoMsgHandler(msg *iso8583.Message) {
	req := EchoMessageRequest{}

	err := msg.Unmarshal(&req)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(req.PrettyPrint())
}
