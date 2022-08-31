/**
 * @author Jose Nidhin
 */
package main

import (
	"log"
)

var logger *log.Logger

func init() {
	logger = log.Default()
}
