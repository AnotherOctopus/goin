//Constants is just all the constants that the currency uses
package constants

const (
	MAXTRANSNETSIZE = 10000 // Maximum size of a serialized transaction
	TRANSBROADPORT = "1917" // Port that transactions are sent on
	TRANSRXPORT    = "1943" // Port that transactions are recievd from
	BLOCKRXPORT    = "1918" // Port that blocks are recived on
	CONN_TYPE      = "tcp" // We use tcp
	NETWORK_INT    = "0.0.0.0" // The ip address of this computer
	ADDRESSSIZE    = 8 // Number of bytes in an address
	PRIVKEYSIZE    = 2000 // How big a private key can get in bytes
	HASHSIZE       = 32 //Size of all the hashes
)