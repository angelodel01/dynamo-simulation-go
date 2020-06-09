import os
import sys

def tester():
    if len(sys.argv) < 2:
        os.system("go run dynamo.go gossip_membership.go hashing.go data.go")
        os.system("rm MEM*")
    elif len(sys.argv) == 2 and sys.argv[1] == "noclean":
        os.system("go run dynamo.go membership.go hashing.go data.go")


if __name__ == '__main__':
    tester()
