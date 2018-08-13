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

        def sendCmd(self,job,filename,walindex="",addridx="",inputs=[],outputs=[]):
                send = {
                        "job": job,
                        "filename":filename,
                        "addridx":addridx,
                        "walindex":walindex,
                        "inputs":inputs,
                        "outputs":outputs
                }
                r = requests.post("http://{}:{}/cmd".format(self.ip,self.cmdport),json=send)
                return r.content

if __name__ == "__main__":
        n1 = Node()
        print n1.sendCmd(4,"newwallet")
