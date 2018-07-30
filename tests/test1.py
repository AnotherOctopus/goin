import goin
NETDIR = "networkfiles"
print("Loading Wallets")
goin.loadWal("genesisWallet","172.18.0.2")
for i in range(3,9):
    print("Wallet")
    goin.makeWal("newwallet{}".format(i),"172.18.0.{}".format(i))
    goin.loadWal("newwallet{}".format(i),"172.18.0.{}".format(i))