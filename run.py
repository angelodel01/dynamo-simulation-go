import os
import sys

def tester():
    if len(sys.argv) < 2:
        os.system("go run dynamo.go membership.go hashing.go data.go")
    elif len(sys.argv) == 2 and sys.argv[1] == "clean":
        os.system("rm MEM*")

if __name__ == '__main__':
    tester()
