package main


import(
  "fmt"
)

func main() {
  simulate_con_hash()
  // num_nodes = 8
  // simulate_gossip()
}

func simulate_gossip(){
  hash_num_nodes = num_nodes
  member_ch := make(chan map[int]map[int]NodeHB)
  for me := 0; me < num_nodes; me++ {
      my_HB_Table := make(map[int]NodeHB)
      n := chooseNeighbors(me)
      for i := 0; i < num_neighbors ; i++{
        my_HB_Table[n[i]] = NodeHB{id: n[i], Hbcounter: 0, time: 0, dead: false}
      }
      spawnNodeHB(NodeHB{id: me, Hbcounter: 0, time: 0, dead: false}, my_HB_Table, member_ch)
  }
  wg.Wait()
}

func simulate_dynamo(){
    for r := range ring{
      ring[r] = -1
    }
    for r := range request_ch{
      request_ch[r] = nil
      response_ch[r] = nil
    }
    AddNodeHash(0)
    AddNodeHash(1)
    AddNodeHash(2)
    AddNodeHash(3)
    AddNodeHash(4)
    AddNodeHash(5)
    AddNodeHash(6)
    AddNodeHash(7)
    fmt.Printf("Ring %+v\n\n", ring)

    hash_num_nodes = num_nodes

    simulate_gossip()
    put("Maria", 100)
    put("John", 20)
    put("Anna", 40)
    put("Alex", 10)

    get("Alex")
    get("Anna")
    get("Maria")
    get("John")

    DeleteNodeHash(0)

    get("Maria")
    get("John")

    DeleteNodeHash(1)
    DeleteNodeHash(2)
    DeleteNodeHash(3)
    DeleteNodeHash(4)
    DeleteNodeHash(5)
    DeleteNodeHash(6)
    DeleteNodeHash(7)
}


func simulate_con_hash(){
    for r := range ring{
      ring[r] = -1
    }
    for r := range request_ch{
      request_ch[r] = nil
      response_ch[r] = nil
    }
    AddNodeHash(0)
    AddNodeHash(1)
    AddNodeHash(2)
    AddNodeHash(3)
    fmt.Printf("Ring %+v\n\n", ring)


    put("Maria", 100)
    put("John", 20)
    put("Anna", 40)
    put("Tim", 100)
    put("Alex", 10)

    get("Tim")
    get("Alex")

    DeleteNodeHash(0)

    get("Anna")
    get("Maria")
    get("John")
    get("Maria")
    get("John")

    DeleteNodeHash(1)
    DeleteNodeHash(2)
    DeleteNodeHash(3)
}
