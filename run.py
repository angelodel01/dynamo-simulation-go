import os
import sys

def tester():
    if len(sys.argv) < 1:
        os.system("go run dynamo.go membership.go hashing.go data.go")
    elif len(sys.argv) == 1 and sys.argv[0] == "clean":
        os.system("rm MEM*")

if __name__ == '__main__':
    tester()
