package main

import (
  "fmt"
  "time"
  "math/rand"
  "sync"
)

type NodeHB struct{
  id int
  Hbcounter int
  time int
  dead bool
}


var wg sync.WaitGroup
var HB_mutex sync.Mutex
const num_NodeHBs = 8
const num_neighbors = 2
const max_cycles = 20
const cycle_time = 2

func main() {
  member_ch := make(chan map[int]map[int]NodeHB)
  for me := 0; me < num_NodeHBs; me++ {
      my_HB_Table := make(map[int]NodeHB)
      n := chooseNeighbors(me)
      for i := 0; i < num_neighbors; i++{
        my_HB_Table[n[i]] = NodeHB{id: n[i], Hbcounter: 0, time: 0, dead: false}
      }
      spawnNodeHB(NodeHB{id: me, Hbcounter: 0, time: 0, dead: false}, my_HB_Table, member_ch)
  }
  wg.Wait()
}

func chooseNeighbors(me int) [num_neighbors]int {
  var n [num_neighbors]int
  for i := 0; i < num_neighbors; i++{
    var curr = rand.Intn(num_NodeHBs)
    for curr == me || curr == n[0]{
      curr = rand.Intn(num_NodeHBs)
    }
    n[i] = curr
  }
  fmt.Printf("Chose neighbors %v, for NodeHB: %d\n", n, me)
  return n
}

func spawnNodeHB(my_NodeHB NodeHB, my_HB_Table map[int]NodeHB, member_ch chan map[int]map[int]NodeHB){
  wg.Add(2)
  go updateHeartBeats(my_NodeHB, my_HB_Table, member_ch)
  go listenForTraffic(my_NodeHB, my_HB_Table, member_ch)
}

func listenForTraffic(my_NodeHB NodeHB, my_HB_Table map[int]NodeHB,
                      member_ch chan map[int]map[int]NodeHB){
  defer wg.Done()
  for i := 0; i < max_cycles; i++{//listening on channel
    var mp = <-member_ch
    for k, v := range mp{//should only give us one iteration
      HB_mutex.Lock()
      _, found := my_HB_Table[k]
      HB_mutex.Unlock()
      if found {//if the information coming in is from a neighbor
        updateTable(k, my_NodeHB, v, my_HB_Table)
      }
    }
  }
}

func updateTable(sender_NodeHB_id int, my_NodeHB NodeHB, new_values map[int]NodeHB, my_HB_Table map[int]NodeHB){
  for k, v := range new_values{//for all the information coming in
    value, found := my_HB_Table[k]
    if found && !value.dead{//if the stuff in the incoming table is in the neighborhood
      if v.time > value.time && v.Hbcounter <= value.Hbcounter{
        // fmt.Printf("NodeHB %d, has killed NodeHB %d\n", my_NodeHB.id, v.id)
        HB_mutex.Lock()
        my_HB_Table[k] = NodeHB{id: v.id, Hbcounter: v.Hbcounter, time: v.time, dead: true}
        HB_mutex.Unlock()
        fmt.Printf("NodeHB %d, has killed NodeHB %d\n" + "-found %d in table from NodeHB %d\n-updating: %+v to: %+v\n"+ "-NEW NodeHB %d TABLE: %+v\n\n", my_NodeHB.id, v.id, k, sender_NodeHB_id, value, my_HB_Table[k], my_NodeHB.id,my_HB_Table)
      }else if v.time > value.time {//if the information is more recent
        HB_mutex.Lock()
        my_HB_Table[k] = v
        HB_mutex.Unlock()
        fmt.Printf("For NodeHB : %d\n" + "-found %d in table from NodeHB %d\n-updating: %+v to: %+v\n"+ "-NEW NodeHB %d TABLE: %+v\n\n", my_NodeHB.id, k, sender_NodeHB_id, value, v, my_NodeHB.id,my_HB_Table)
      }
    }
  }
}

// mark NodeHB as failed when failed
func updateHeartBeats(my_NodeHB NodeHB, my_HB_Table map[int]NodeHB,
                      member_ch chan map[int]map[int]NodeHB){
  timer2 := time.NewTimer(time.Second*cycle_time)
  defer wg.Done()
  var sender_map = make(map[int]map[int]NodeHB)
  my_HB_Table[my_NodeHB.id] = my_NodeHB
  fmt.Printf("NodeHB: %+v, Initial table: %+v\n", my_NodeHB, my_HB_Table)
  sender_map[my_NodeHB.id] = my_HB_Table
  for i := 0; i < max_cycles; i++{
    <-timer2.C
    if !(i%2 == 0 && my_NodeHB.id%3 == 0){//generating failures for certain NodeHBs
      my_NodeHB.Hbcounter += 1
    }
    my_NodeHB.time += 1
    HB_mutex.Lock()
    my_HB_Table[my_NodeHB.id] = my_NodeHB
    HB_mutex.Unlock()
    fmt.Printf("Update for NodeHB %d => time: %d, HB: %d\n", my_NodeHB.id, my_NodeHB.time, my_NodeHB.Hbcounter)
    sender_map[my_NodeHB.id] = my_HB_Table
    member_ch <- sender_map
    timer2 = time.NewTimer(time.Second)
  }
}