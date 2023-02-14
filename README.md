# Private Comms
This is a personal project to create a private communication system between 2 remote PCs through the terminal. It is not intended for public use.

## How to use
Note: You must have the responder's generated public encoded key (in string) to use this program. This is to ensure message authenticity.

1. Clone the repository
```bash
git clone https://github.com/galadd/private-comms.git
```
3. Run the program
```bash
go run main.go
```
If you generate a new keypair, create an additional key, "RESPONDER_PUB_KEY", in the .env file and set it to the encoded public key of the responder. 

4. Follow the program prompts

## How it works
The program uses the RSA algorithm to generate a keypair and sign/verify messages. The program uses AES-256 to encrypt/decrypt messages. The program uses the Diffie-Hellman key exchange algorithm to generate a shared secret key. This is the key used by AES-256 to encrypt/decrypt messages. 

## Future plans
- [ ] Add a way to send messages to multiple people
