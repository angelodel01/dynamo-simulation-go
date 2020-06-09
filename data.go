package main


import(
  "sync"
)


type Message struct{
  command string
  node_id int
  key string
  val int
}


type NodeHB struct{
  id int
  Hbcounter int
  time int
  dead bool
}


type NodeHash struct{
  id int
  hash_1 int
  hash_2 int
}

///////////////////////////////

var wg_gossip sync.WaitGroup
var wg_hash sync.WaitGroup
var HB_mutex sync.Mutex
// const num_nodes = 8
const num_neighbors = 2
const max_cycles = 1000
const cycle_time = 2

///////////////////////////////

const mem_size = 50
var num_nodes = 0
var hash_num_nodes = 0
var ring [mem_size]int
var ring_mutex = &sync.Mutex{}
var request_ch [mem_size]chan Message
var response_get_ch [mem_size]chan Message
var response_put_ch [mem_size]chan Message
var member_ch map[int]map[int]chan NodeHB
