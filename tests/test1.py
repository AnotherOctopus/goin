from goin import Node
import requests
NETDIR = "networkfiles"
print("Loading Wallets")
rootnode = Node("172.18.0.2")
rootnode.sendCmd("genesisWallet")
nodes = []
for i in range(3,11):
    n = Node("172.18.0.{}".format(i))
    n.makeWal("newwallet{}".format(i))
    nodes.append(n)
    wallet = requests.get("http://172.18.0.{}/newwallet{}".format(i,i)) 
    print(wallet.content)