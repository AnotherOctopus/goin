package main

import (
	"wallet"
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
	"cnet"
	"io/ioutil"
	"encoding/hex"
	"time"
	"encoding/json"
	"encoding/base64"
	"constants"
)
func main(){
	//test()
	peerips := []string{"127.0.0.1"}
	nd := cnet.New(peerips)
	go nd.TxListener()
	go nd.BlListener()
	go nd.CmdListener()
}