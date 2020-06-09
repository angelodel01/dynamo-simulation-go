package main

import (
  "fmt"
  "time"
  "math/rand"
)

func chooseNeighbors(me int) [num_neighbors]int {
  var n [num_neighbors]int
  for i := 0; i < num_neighbors; i++{
    var curr = rand.Intn(hash_num_nodes)
    for curr == me || curr == n[0]{
      curr = rand.Intn(hash_num_nodes)
      fmt.Println("After int generation 2")
    }
    n[i] = curr
  }
  fmt.Printf("Chose neighbors %v, for NodeHB: %d\n", n, me)
  return n
}

func spawnNodeHB(my_NodeHB NodeHB, my_HB_Table map[int]NodeHB, member_ch chan map[int]map[int]NodeHB){
  wg_gossip.Add(2)
  go updateHeartBeats(my_NodeHB, my_HB_Table, member_ch)
  go listenForTraffic(my_NodeHB, my_HB_Table, member_ch)
}

func listenForTraffic(my_NodeHB NodeHB, my_HB_Table map[int]NodeHB,
                      member_ch chan map[int]map[int]NodeHB){
  var done bool = false
  for ; done == false; {//listening on channel
    fmt.Println("Waiting on message at node: ", my_NodeHB.id)
    var mp = <-member_ch
    for k, v := range mp{//should only give us one iteration
      HB_mutex.Lock()
      _, found := my_HB_Table[k]
      HB_mutex.Unlock()
      if found {//if the information coming in is from a neighbor
        updateTable(k, my_NodeHB, v, my_HB_Table)
      }
    }
    HB_mutex.Lock()
    me, _ := my_HB_Table[my_NodeHB.id]
    if me.Hbcounter == -1 {
        done = true
    }
    HB_mutex.Unlock()
  }
  fmt.Println("Cleanly exiting listenForTraffic node: ", my_NodeHB.id)
  wg_gossip.Done()
}

func updateTable(sender_NodeHB_id int, my_NodeHB NodeHB, new_values map[int]NodeHB, my_HB_Table map[int]NodeHB){
  for k, v := range new_values{//for all the information coming in
    value, found := my_HB_Table[k]
    if found && !value.dead{//if the stuff in the incoming table is in my table and I didn't mark it dead already
        fmt.Printf("I AM NODE: %d, for node: %d t = %d hb = %d, my t = %d my hb = %d\n", my_NodeHB.id, v.id, v.time, v.Hbcounter, value.time, value.Hbcounter)
      if v.time > value.time && v.Hbcounter <= value.Hbcounter{//If the incoming table's time value has been updated but its heartbeat counter is the same
        fmt.Printf("NodeHB %d, has killed NodeHB %d\n" + "-found %d in table from NodeHB %d\n-updating: %+v to: %+v\n"+ "-NEW NodeHB %d TABLE: %+v\n\n", my_NodeHB.id, v.id, k, sender_NodeHB_id, value, v, my_NodeHB.id,my_HB_Table)
        fmt.Printf("my_HB_Table[k].dead %v\n", my_HB_Table[k].dead)
        HB_mutex.Lock()
		my_HB_Table[k] = NodeHB{id: v.id, Hbcounter: v.Hbcounter, time: v.time, dead: true}
		HB_mutex.Unlock()
        if my_HB_Table[k].dead{
          DeleteNodeHash(v.id)
        }
      }else if v.dead{
        HB_mutex.Lock()
        my_HB_Table[k] = v
        HB_mutex.Unlock()
        fmt.Printf("For NodeHB : %d\n" + "-found %d in table from NodeHB %d\n-updating: %+v to: %+v\n"+ "-NEW NodeHB %d TABLE: %+v\n\n", my_NodeHB.id, k, sender_NodeHB_id, value, v, my_NodeHB.id,my_HB_Table)
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
  dead_time := 0
  var sender_map = make(map[int]map[int]NodeHB)
  my_HB_Table[my_NodeHB.id] = my_NodeHB
  fmt.Printf("NodeHB: %+v, Initial table: %+v\n", my_NodeHB, my_HB_Table)
  sender_map[my_NodeHB.id] = my_HB_Table
  for i := 0; i < max_cycles; i++{
    <-timer2.C
    my_NodeHB.time += 1
    if i == 2 && my_NodeHB.id == 0{
        dead_time = 5
    }
    if (dead_time == 0) {
        my_NodeHB.Hbcounter += 1
    } else {
        dead_time--
    }
    HB_mutex.Lock()
    my_HB_Table[my_NodeHB.id] = my_NodeHB
    HB_mutex.Unlock()
    fmt.Printf("Update for NodeHB %d => time: %d, HB: %d\n", my_NodeHB.id, my_NodeHB.time, my_NodeHB.Hbcounter)
    HB_mutex.Lock()
    sender_map[my_NodeHB.id] = my_HB_Table
    HB_mutex.Unlock()
    member_ch <- sender_map
    timer2 = time.NewTimer(time.Second)
  }
  fmt.Println("Cleanly exiting updateHeartBeats")
  HB_mutex.Lock()
  my_NodeHB.Hbcounter = -1
  my_HB_Table[my_NodeHB.id] = my_NodeHB
  HB_mutex.Unlock()
  wg_gossip.Done()
}
