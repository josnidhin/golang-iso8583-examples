/**
 * @author Jose Nidhin
 */
package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/moov-io/iso8583"
)

const HEADER_SIZE = 12

var rawMessages []string = []string{
	"ISO02110005502007238800008808000101234567890000000000000010000092618372419060118510009260926000000005928MON50EDIOX     N484",
	"ISO0211000550210723080000E8080001012345678900000000000000100000926183724190601185100092600000000592812345600MON50EDIOX     N484",
	"ISO0211000550420723880000E80800010123456789000000000000001000009261837241906011851000926092600000000592812345600MON50EDIOX     N484",
	"ISO0211000550430722080000A8080001012345678900000000000000100000926183724190601092600000000592800MON50EDIOX     N484",
	"ISO0211000550800822000000000000004000000000000000821083216015795301",
	"ISO021100055081082200000020000000400000000000000082108321601579500301",
}

func main() {
	for _, rawMsg := range rawMessages {
		fmt.Printf("Raw Message = %s\n", rawMsg)

		tw := tabwriter.NewWriter(os.Stdout, 2, 2, 1, ' ', 0)

		msg := iso8583.NewMessage(Spec1)
		msg.Unpack([]byte(rawMsg[HEADER_SIZE:]))

		for pos := 0; pos < 128; pos++ {
			value, err := msg.GetString(pos)

			if err != nil {
				continue
			}

			if value == "" {
				continue
			}

			field := msg.GetField(pos)

			fmt.Fprintf(tw, "%3d\t%s\t%s\n", pos, field.Spec().Description, value)
		}
		tw.Flush()

		fmt.Println()
	}
}
