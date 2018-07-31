import goin
import requests
NETDIR = "networkfiles"
print("Loading Wallets")
goin.loadWal(NETDIR + "/genesisWallet","172.18.0.2")
for i in range(3,9):
    goin.makeWal("newwallet{}".format(i),"172.18.0.{}".format(i))
    wallet = requests.get("172.18.0.{}".format(i)) 
    print(wallet)
    goin.loadWal("newwallet{}".format(i),"172.18.0.{}".format(i))
goin.prepTx()