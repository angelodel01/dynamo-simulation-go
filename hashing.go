package main
import(
  "fmt"
  "time"
  "crypto/md5"
  "encoding/binary"
  "strconv"
)

func AddNodeHash(id int){
  num_nodes++
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
    if request.command == "KILL" && request.node_id == me.id{
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
    } else if (request.command == "GET" || request.command == "PUT") && (request.node_id == me.id) {
      response_ch[me.id] <- Message{command: "CONFIRM", node_id: me.id}
    }
  }

}


func DeleteNodeHash(id int){
  request_ch[id] <- Message{command: "KILL", node_id: id}
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
  node_id := ring[hash]
  ring_mutex.Unlock()
  time.Sleep(100 * time.Millisecond)
  request_ch[node_id] <- Message{command: "GET", node_id: node_id}
  resp := <-response_ch[node_id]
  if resp.command == "CONFIRM"{
    fmt.Printf("Performing get(key: '%s') hashes to: %d, got NodeHash id: %d\n\n", key, original, resp.node_id)
    return resp.node_id
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
  node_id := ring[hash]
  ring_mutex.Unlock()
  time.Sleep(100 * time.Millisecond)
  request_ch[node_id] <- Message{command: "PUT", node_id: node_id}
  resp := <-response_ch[node_id]
  if resp.command == "CONFIRM"{
    fmt.Printf("Performing put(key: '%s', value: %d) hashes to: %d, got NodeHash id: %d\n\n", key, value, original, resp.node_id)
    return resp.node_id
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
