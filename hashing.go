package main
import(
  "fmt"
  "time"
  "crypto/md5"
  "encoding/binary"
  "strconv"
  "bufio"
  "os"
  "strings"
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
  response_get_ch[id] = make(chan Message)
  response_put_ch[id] = make(chan Message, 3)
  wg.Add(1)
  go NodeHashRoutine(n)
}


func NodeHashRoutine(me NodeHash){
  // defer wg.Done()
  flag := true
  file_name := "MEM" + strconv.Itoa(me.id) + ".txt"
  file, _ := os.Create(file_name)
  file.Close()
  for flag{
    request := <-request_ch[me.id]
    if request.command == "KILL" && request.node_id == me.id{
      fmt.Printf("Node: %d received KILL command\n", me.id)
      ring_mutex.Lock()
      ring[me.hash_1] = -1
      ring[me.hash_2] = -1
      ring_mutex.Unlock()
      time.Sleep(time.Second)
      fmt.Printf("NodeHash : %d dies\n\n", me.id)
      fmt.Printf("Ring %+v\n\n", ring)
      close(response_get_ch[me.id])
      close(response_put_ch[me.id])
      close(request_ch[me.id])
      response_get_ch[me.id] = nil
      response_put_ch[me.id] = nil
      request_ch[me.id] = nil
      flag = false
      wg.Done()
    } else if (request.command == "GET" && request.node_id == me.id) {
      // fmt.Println("GET REQUEST IN NodeHashRoutine")
      read_val := readFromFile(request.node_id, request.key)
      response_get_ch[me.id] <- Message{command: "CONF_GET", node_id: me.id, key: request.key, val: read_val}
    } else if (request.command == "PUT" && request.node_id == me.id){
      writeToFile(request.node_id, request.val, request.key)
      response_put_ch[me.id] <- Message{command: "CONF_WR", node_id: me.id, key: "N/A", val: -1}
    }
  }

}


func DeleteNodeHash(id int){
  fmt.Printf("requesting node death in DeleteNodeHash id: %d\n", id)
  // request_ch[id] <- Message{command: "KILL", node_id: id, key: "N/A", val: -1}
}


func get(key string) int{
  fmt.Println("Inside GET")
  hash := GetMD5HashString(key)%mem_size
  original := hash
  ring_mutex.Lock()
  // fmt.Printf("Ring in get %+v\n\n", ring)
  for ring[hash] == -1{
    if hash >= mem_size - 1{
      hash = 0
    }
    hash++
  }
  node_id := ring[hash]
  ring_mutex.Unlock()
  time.Sleep(100 * time.Millisecond)
  request_ch[node_id] <- Message{command: "GET", node_id: node_id, key: key, val: -1}
  resp := <-response_get_ch[node_id]
  if resp.command == "CONF_GET"{
    fmt.Printf("Performed get(key: '%s') hashes to: %d, got NodeHash id: %d\nGot Value: %d, from Key: %s\n\n", key, original, resp.node_id, resp.val, resp.key)
    return resp.node_id
  }
  return -1
}


func put(key string, value int) int{
  hash := GetMD5HashString(key)%mem_size
  original := hash
  var ids [3]int
  for i := 0; i < 3; i++{
    ring_mutex.Lock()
    for ring[hash] == -1{
      if hash >= mem_size - 1{
        hash = 0
      }
      hash++
    }
    node_id := ring[hash]
    ring_mutex.Unlock()
    ids[i] = node_id
    hash++
  }
  for j, _ := range ids{
    request_ch[ids[j]] <- Message{command: "PUT", node_id: ids[j], key: key, val: value}
  }
  for k, _ := range ids{
    resp := <-response_put_ch[ids[k]]
    if resp.command == "CONF_WR"{
      fmt.Printf("Performing put(key: '%s', value: %d) hashes to: %d, got NodeHash id: %d\n\n", key, value, original, resp.node_id)
      return resp.node_id
    }
  }
  return -1
}



func writeToFile(node_id int, value int, key string){
  file_name := "MEM" + strconv.Itoa(node_id) + ".txt"
  file, err := os.OpenFile(file_name, os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    panic(err)
  }
  _, err = file.WriteString(fmt.Sprintf("%s:%d\n", key, value))
  defer file.Close()
}


func readFromFile(node_id int, key string) int{
  file_name := "MEM" + strconv.Itoa(node_id) + ".txt"
  file, err := os.Open(file_name)
  if err != nil {
    panic(err)
  }
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
  var txtlines []string

  for scanner.Scan() {
    txtlines = append(txtlines, scanner.Text())
  }

  file.Close()

  for _, eachline := range txtlines {
    mem_lst := strings.Split(eachline, ":")
    if mem_lst[0] == key{
      num, _ := strconv.Atoi(mem_lst[1])
      return num
    }
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
