package main


import(
  "fmt"
)

func main() {
  simulate_con_hash()
  simulate_gossip()
}

func simulate_gossip(){
  member_ch := make(chan map[int]map[int]NodeHB)
  for me := 0; me < num_nodes; me++ {
      my_HB_Table := make(map[int]NodeHB)
      n := chooseNeighbors(me)
      for i := 0; i < num_neighbors; i++{
        my_HB_Table[n[i]] = NodeHB{id: n[i], Hbcounter: 0, time: 0, dead: false}
      }
      spawnNodeHB(NodeHB{id: me, Hbcounter: 0, time: 0, dead: false}, my_HB_Table, member_ch)
  }
  wg.Wait()
}

func simulate_con_hash(){
    for r := range ring{
      ring[r] = -1
    }
    for r := range request_ch{
      request_ch[r] = nil
      response_ch[r] = nil
    }
    AddNodeHash(num_nodes)
    AddNodeHash(num_nodes)
    AddNodeHash(num_nodes)
    AddNodeHash(num_nodes)
    AddNodeHash(num_nodes)
    fmt.Printf("Ring %+v\n\n", ring)

    put("Maria", 100)
    put("John", 20)
    put("Anna", 40)
    put("Tim", 100)
    put("Alex", 10)

    get("Tim")
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
}
