import sys
import socket
import time
CMDPORT=1945

class Wallet(object):
        def __init__(self,n,idx,filename=None):
                if not filename:
                        if n.wallIdx = 0:
                                raise ValueError("This node has no wallets")
                        if idx >= n.walIdx:
                                raise ValueError("The node does not have wallets that high")
                        n.saveWal("walletu")
                        walletraw = requests.get("http://{}/wallet{}".format(i,i)) 
                

class Node(object):
        def __init__(self,ip):
                self.ip = ip
                self.walIdx = 0

        def sendString(self,toSend):
                sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

                # Connect the socket to the port where the server is listening
                server_address = (self.ip, CMDPORT)
                sock.connect(server_address) 
                # Create a TCP/IP socket
                sock.send(toSend)
                sock.close()
                
                time.sleep(1)

        
        def sendTx(self,wallIdx, addressIdx,txTx):
                self.sendString("1")
                self.sendString(wallIdx)
                self.sendString(addressIdx)
                self.sendString(txTx)

        def prepTx(self,txFilename, inputs, outputs):
                self.sendString("2")
                self.sendString(txFilename)
                for inp in inputs:
                        self.sendString("i")
                        self.sendString(inp.Hash)
                        self.sendString(inp.Idx)
                for oup in outputs:
                        self.sendString("o")
                        self.sendString(oup.Hash)
                        self.sendString(oup.Amount)

        def viewBal(self):
                self.sendString("3")

        def makeWal(self,walFilename):
                self.sendString("4")
                self.sendString(walFilename)
                self.walIdx += 1

        def loadWal(self,walFilename):
                self.sendString("5")
                self.sendString(walFilename)
        def saveWal(self,walFilename,idx):
                self.sendString("6")
                self.sendString(walFilename)
                self.sendString(string(idx))

        def done(self):
                self.sendString("10")

if __name__ == "__main__":
        n1 = Node("172.18.0.2")
        n1.makeWal("newwallet")
        n1.done()
       
