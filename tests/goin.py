import sys
import requests
import socket
import time

#[1] Send Transaction From File"
#[2] Manually Prepare Transaction To File"
#[3] View Current Balence To File"
#[4] Make A New Wallet To File"
#[5] Load A Wallet From File"
#[6] Save A Wallet To File
#Other exits

#expects json of
#{
#	<key>: <valuetype>, <which tasks require>
#	"job": int, [1,2,3,4,5,6]
#	"filename":string, [1,2,3,4,5,6]
#	"walindex":int,[1,3,6]
#	"addridx": int, [1]
#	"inputs": [(Hash,int)], [2]
#	"outputs": [(Hash,int)] [2]
#}
#*/
class Node(object):
        def __init__(self,ip="localhost"):
                self.cmdport = 1945
                self.ip = ip
                self.walIdx = 0

        def sendTx(self,filename,walindex,addridx):
                return self.sendCmd(1,
                                filename,
                                walindex=walindex,
                                addridx=addridx
                                )
        def loadWallet(self,filename):
                self.walIdx += 1
                return self.sendCmd(5,filename)
        def prepTx(self,filename,walindex,addridx,inputs,outputs):
                return self.sendCmd(2,filename,
                                walindex=walindex,
                                addridx=addridx,
                                inputs=inputs,
                                outputs=outputs
                                )
        def newWallet(self,filename):
                self.sendCmd(4,filename)
                self.loadWallet(filename)
        def sendCmd(self,job,filename,walindex="",addridx="",inputs=[],outputs=[]):
                send = {
                        "job": job,
                        "filename":filename,
                        "addridx":addridx,
                        "walidx":walindex,
                        "inputs":inputs,
                        "outputs":outputs
                }
                r = requests.post("http://{}:{}/cmd".format(self.ip,self.cmdport),json=send)
                return r.content
        def getClaimedTxs(self):
                txByWallet = requests.get("http://{}:{}/claimedtxs".format(self.ip,self.cmdport)).content
                txByWallet = txByWallet.strip('%')
                txByWallet = txByWallet.split('%')
                txs = [tx.strip('#').split('#') for tx in txByWallet]
                return txs


def getAddresses(ip,cmdport):
        raw = requests.get("http://{}:{}/addresses".format(ip,cmdport)).content
        raw = raw.strip('%')
        raw = raw.split('%')
        addrs = [addr.strip('#').split('#') for addr in raw]
        return addrs

if __name__ == "__main__":
        NETDIR = "networkfiles"
        print("Loading Wallets")
        n1 = Node("localhost")
        n1.loadWallet("networkfiles/genesisWallet")
        addressToSendTo = getAddresses(n1.ip,1945)[0][0]
        txToSend = n1.getClaimedTxs()[0][0]
        print "Claimed Transactions", txToSend
        print "Lets send {} goins to {}".format(txToSend,addressToSendTo)
        print n1.prepTx("trax1",0,0,[{"hash":txToSend,"idx":0}],[{"hash":addressToSendTo,"amt":100}])
        print n1.sendTx("trax1",0,0)
