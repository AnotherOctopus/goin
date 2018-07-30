import goin
NETDIR = "networkfiles"
goin.loadWal("genesisWallet".format(NETDIR),"172.18.0.2")
for i in range(3,9):
    goin.makeWal("newwallet{}".format(i),"172.18.0.{}".format(i))
    goin.loadWal("newwallet{}".format(i),"173.18.0.{}".format(i))