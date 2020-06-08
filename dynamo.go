package main


import(
  "fmt"
  "time"
)

func main() {
  simulate_dynamo()
  // simulate_con_hash()
  // num_nodes = 8
  // simulate_gossip()
}

func simulate_gossip(){
  hash_num_nodes = num_nodes
  member_ch := make(chan map[int]map[int]NodeHB)
  fmt.Println("before spawning nodes")
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
      response_put_ch[r] = nil
      response_get_ch[r] = nil
    }
    AddNodeHash(0)
    AddNodeHash(1)
    AddNodeHash(2)
    fmt.Printf("Ring %+v\n\n", ring)
    num_nodes = 3
    hash_num_nodes = num_nodes


    simulate_gossip()
	
    DeleteNodeHash(1)
    DeleteNodeHash(2)

    wg.Wait()
}


func simulate_con_hash(){
    for r := range ring{
      ring[r] = -1
    }
    for r := range request_ch{
      request_ch[r] = nil
      response_put_ch[r] = nil
      response_get_ch[r] = nil
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
    get("Anna")
    get("Maria")
    get("John")
    DeleteNodeHash(1)

    get("Maria")
     get("John")

     fmt.Println("before deletes")
    DeleteNodeHash(0)

    DeleteNodeHash(2)
    DeleteNodeHash(3)
    wg.Wait()
    time.Sleep(2* time.Second)
}
