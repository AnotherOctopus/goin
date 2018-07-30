import goin

for i in range(2,9):
    goin.makeWal("newwallet{}".format(i),"172.18.0.{}".format(i))