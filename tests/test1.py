from goin import Node,getAddresses
import requests
NETDIR = "networkfiles"
print("Loading Wallets")
n1 = Node("172.18.0.2")
n2 = Node("172.18.0.3")
n1.loadWallet("networkfiles/genesisWallet")
n2.newWallet("n2wall")

addressToSendTo = getAddresses(n2.ip,1945)[0][0]
txToSend = n1.getClaimedTxs()[0][0]
print "Claimed Transactions", txToSend
print "Lets send {} goins to {}".format(txToSend,addressToSendTo)
n1.prepTx("trax1",0,0,[{"hash":txToSend,"idx":0}],[{"hash":addressToSendTo,"amt":100}])
print n1.sendTx("trax1",0,0)
print "TXT SENT"

addressToSendTo = getAddresses(n1.ip,1945)[0][0]
txToSend = n2.getClaimedTxs()[0][0]

print "Claimed Transactions", txToSend
'''
print "Lets send {} goins to {}".format(txToSend,addressToSendTo)
print n2.prepTx("trax1",0,0,[{"hash":txToSend,"idx":0}],[{"hash":addressToSendTo,"amt":100}])
print n2.sendTx("trax1",0,0)
'''
