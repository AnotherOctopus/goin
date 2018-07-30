import sys
import socket
import time
CMDPORT=1945

def sender(ip):
        def sendIP(toSend):
                sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

                # Connect the socket to the port where the server is listening
                server_address = (ip, CMDPORT)
                print(server_address)
                sock.connect(server_address) 
                # Create a TCP/IP socket
                sock.send(toSend)
                sock.close()
                
                time.sleep(1)
        return sendIP
        
def sendTx(wallIdx, addressIdx,txTx,ip):
        sendString = sender(ip)
        sendString("1")
        sendString(wallIdx)
        sendString(addressIdx)
        sendString(txTx)

def prepTx(txFilename, inputs, outputs,ip):
        sendString = sender(ip)
        sendString("2")
        sendString(txFilename)
        sendString(addressIdx)
        sendString(txTx)

def viewBal(ip):
        sendString = sender(ip)
        sendString("3")

def makeWal(walFilename,ip):
        sendString = sender(ip)
        sendString("4")
        sendString(walFilename)

def loadWal(walFilename,ip):
        sendString = sender(ip)
        sendString("5")
        sendString(walFilename)

def done(ip):
        sendString = sender(ip)
        sendString("10")

if __name__ == "__main__":
        options = { "1": sendTx,
                    "2": prepTx,
                    "3": viewBal,
                    "4": makeWal,
                    "5": loadWal,
                   "10": done
                  }
        makeWal("newwallet","172.18.0.2")
        done("172.18.0.2")
       
