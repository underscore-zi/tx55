package main

import (
	"flag"
	"fmt"
	"os"
	"tx55/pkg/metalgearonline1/handlers"
	"tx55/pkg/metalgearonline1/testclient"
	"tx55/pkg/metalgearonline1/types"
)

func main() {
	addr := flag.String("address", "", "An ip:port string reflecting the TCP address of the server")
	flag.Parse()

	c := testclient.TestClient{Key: types.XORKEY}

	if err := c.Connect(*addr); err != nil {
		fmt.Println("FAIL1")
		os.Exit(1)

	}

	// I choose to use LoginWithSession because it is a simple request, just a session ID. To fulfil the server needs
	// to do a lookup. So with one request we can ensure the game is still accepting connections and can tell if it is
	// having any database issues without the healthchecker needing to test the database specifically

	notFoundErr := fmt.Sprintf("code(%d)", handlers.ErrNotFound.Code)
	databaseErr := fmt.Sprintf("code(%d)", handlers.ErrDatabase.Code)
	if err := c.LoginWithSession("NotALegitSession"); err != nil {
		switch err.Error() {
		case notFoundErr:
			fmt.Println("OK")
			os.Exit(0)
		case databaseErr:
			fmt.Println("FAIL2")
			os.Exit(2)
		default:
			fmt.Println("FAIL3")
			os.Exit(3)
		}
	}
}
