package intro

import (
	"fmt"
)

func Intro() {
	fmt.Println(`
		┌─┐┌─┐┬  ┌─┐┌┬┐┌┬┐   ┌─┐┌─┐┌─┐┌┬┐┌┬┐┌─┐
		│ ┬├─┤│  ├─┤ ││ ││───├─┘│  │ │││││││└─┐
		└─┘┴ ┴┴─┘┴ ┴─┴┘─┴┘   ┴  └─┘└─┘┴ ┴┴ ┴└─┘	

WELCOME TO GALADD-PCOMMS
Start the program if you have a keypair and responder's public key in .env file
if not, generate a new keypair

If you generate a new keypair, you will need to send the public key to the responder
and also receive the responder's public key, saved to .env file as RESPONDER_PUB_KEY

You need to know if you are the initiator or responder.
If you are the responder, you will also need to know the initiator's PC IP address.

If all of the above is known, start the program and follow the prompts.
	`)
}