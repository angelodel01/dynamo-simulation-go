package main


import(
  "fmt"
  "time"
)

func main() {
  // simulate_dynamo()
   simulate_con_hash()
  // num_nodes = 8
  // simulate_gossip()
}

// func simulate_gossip(){
//   hash_num_nodes = num_nodes
//   member_ch := make(chan map[int]map[int]NodeHB)
//   fmt.Println("before spawning nodes")
//   array := make([]map[int]NodeHB, num_nodes)
//   for me := 0; me < num_nodes; me++ {
//       my_HB_Table := make(map[int]NodeHB)
//       n := chooseNeighbors(me)
//       for i := 0; i < num_neighbors ; i++{
//         my_HB_Table[n[i]] = NodeHB{id: n[i], Hbcounter: 0, time: 0, dead: false}
//       }
// 	  array[me] = my_HB_Table
//   }
//   for me := 0; me < num_nodes; me++ {
// 	spawnNodeHB(NodeHB{id: me, Hbcounter: 0, time: 0, dead: false}, array[me], member_ch)
//   }
//   wg_gossip.Wait()
//   fmt.Println("WAITING ON NODES")
// }

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
    //fmt.Printf("Ring %+v\n\n", ring)
    num_nodes = 3
    hash_num_nodes = num_nodes


    simulate_gossip()

    DeleteNodeHash(1)
    DeleteNodeHash(2)


    wg_hash.Wait()
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
    for i := 0; i < 5; i++{
        AddNodeHash(i)
    }
    num_nodes = 5
    hash_num_nodes = num_nodes


    //fmt.Printf("Ring %+v\n\n", ring)
    go simulate_gossip()

    for i := 0; ; i++{
        put("Maria", 100)
        time.Sleep(4* time.Second)
        put("John", 20)
        time.Sleep(10* time.Second)

        get("Maria")
        time.Sleep(4* time.Second)
        put("Anna", 40)
        time.Sleep(2* time.Second)
        get("Anna")
        time.Sleep(4* time.Second)
        put("Tim", 100)
        time.Sleep(4* time.Second)
        put("Alex", 10)
        time.Sleep(4* time.Second)

        get("Tim")
        time.Sleep(4* time.Second)
        get("Alex")
        time.Sleep(4* time.Second)
        get("Anna")
        time.Sleep(4* time.Second)
        get("Maria")
        time.Sleep(4* time.Second)
        get("John")
        time.Sleep(4* time.Second)

        get("Maria")
        time.Sleep(4* time.Second)
        get("John")
        time.Sleep(4* time.Second)
    }
    fmt.Println("before deletes")
    //wg.Wait()
    for{}
}
