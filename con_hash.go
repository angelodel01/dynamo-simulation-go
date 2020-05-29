package main
import(
  "fmt"
  "sync"
  "time"
  "crypto/md5"
  "encoding/binary"
  "strconv"
)

type NodeHash struct{
  id int
  hash_1 int
  hash_2 int
}

type Message struct{
  command string
  NodeHash_id int
}

const mem_size = 50
var num_NodeHashs = 0
var ring [mem_size]int
var wg sync.WaitGroup
var ring_mutex = &sync.Mutex{}
var request_ch [mem_size]chan Message
var response_ch [mem_size]chan Message

func main(){
  for r := range ring{
    ring[r] = -1
  }
  for r := range request_ch{
    request_ch[r] = nil
    response_ch[r] = nil
  }
  AddNodeHash(num_NodeHashs)
  AddNodeHash(num_NodeHashs)
  AddNodeHash(num_NodeHashs)
  AddNodeHash(num_NodeHashs)
  AddNodeHash(num_NodeHashs)
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


func AddNodeHash(id int){
  num_NodeHashs++
  hash_1 := GetMD5HashInt(id)%mem_size
  hash_2 := GetMD5HashInt(hash_1)%mem_size
  ring_mutex.Lock()
  ring[hash_1] = id
  ring[hash_2] = id
  ring_mutex.Unlock()
  n := NodeHash{id: id, hash_1: hash_1, hash_2: hash_2}
  fmt.Printf("Adding NodeHash: %+v\n\n", n)
  request_ch[id] = make(chan Message)
  response_ch[id] = make(chan Message)
  wg.Add(1)
  go NodeHashRoutine(n)
}


func NodeHashRoutine(me NodeHash){
  defer wg.Done()
  flag := true
  for flag{
    request := <-request_ch[me.id]
    if request.command == "KILL" && request.NodeHash_id == me.id{
      ring_mutex.Lock()
      ring[me.hash_1] = -1
      ring[me.hash_2] = -1
      fmt.Printf("NodeHash : %d dies\n\n", me.id)
      fmt.Printf("Ring %+v\n\n", ring)
      close(response_ch[me.id])
      close(request_ch[me.id])
      ring_mutex.Unlock()
      response_ch[me.id] = nil
      request_ch[me.id] = nil
      flag = false
    } else if (request.command == "GET" || request.command == "PUT") && (request.NodeHash_id == me.id) {
      response_ch[me.id] <- Message{command: "CONFIRM", NodeHash_id: me.id}
    }
  }

}


func DeleteNodeHash(id int){
  request_ch[id] <- Message{command: "KILL", NodeHash_id: id}
}


func get(key string) int{
  hash := GetMD5HashString(key)%mem_size
  original := hash
  ring_mutex.Lock()
  for ring[hash] == -1{
    if hash >= mem_size - 1{
      hash = 0
    }
    hash++
  }
  NodeHash_id := ring[hash]
  ring_mutex.Unlock()
  time.Sleep(100 * time.Millisecond)
  request_ch[NodeHash_id] <- Message{command: "GET", NodeHash_id: NodeHash_id}
  resp := <-response_ch[NodeHash_id]
  if resp.command == "CONFIRM"{
    fmt.Printf("Performing get(key: '%s') hashes to: %d, got NodeHash id: %d\n\n", key, original, resp.NodeHash_id)
    return resp.NodeHash_id
  }
  return -1
}


func put(key string, value int) int{
  hash := GetMD5HashString(key)%mem_size
  original := hash
  ring_mutex.Lock()
  for ring[hash] == -1{
    if hash >= mem_size - 1{
      hash = 0
    }
    hash++
  }
  NodeHash_id := ring[hash]
  ring_mutex.Unlock()
  time.Sleep(100 * time.Millisecond)
  request_ch[NodeHash_id] <- Message{command: "PUT", NodeHash_id: NodeHash_id}
  resp := <-response_ch[NodeHash_id]
  if resp.command == "CONFIRM"{
    fmt.Printf("Performing put(key: '%s', value: %d) hashes to: %d, got NodeHash id: %d\n\n", key, value, original, resp.NodeHash_id)
    return resp.NodeHash_id
  }
  return -1
}


func GetMD5HashString(text string) int {
   hash := md5.Sum([]byte(text))
   data := int(binary.BigEndian.Uint64(hash[:8]))
   if data < 0{
     data *= -1
   }
   return data
}


func GetMD5HashInt(text int) int {
   hash := md5.Sum([]byte(strconv.Itoa(text)))
   data := int(binary.BigEndian.Uint64(hash[:8]))
   if data < 0{
     data *= -1
   }
   return data
}
