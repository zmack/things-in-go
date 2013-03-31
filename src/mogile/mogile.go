package mogile

import(
  "net"
  "fmt"
  "net/url"
  "strings"
  "net/http"
  "io/ioutil"
  "strconv"
  "errors"
)

type MogileClient struct {
  connection net.Conn
}

func Connect(address string) (*MogileClient, error) {
  conn, err := net.Dial("tcp", address)
  if err != nil {
    return &MogileClient{}, err
  }
  return &MogileClient{connection: conn}, nil
}

func (client MogileClient) RunCommand(command string, params url.Values) (url.Values, error) {
  client_command := []byte(command + " " + params.Encode() + "\r\n")
  _, err := client.connection.Write(client_command)

  if err != nil {
    fmt.Print(err)
  }

  read_bytes := make([]byte, 1024 * 1024)
  bytes_read, err := client.connection.Read(read_bytes)

  if err != nil {
    fmt.Print(err)
    return nil, nil
  }

  index := strings.Index(string(read_bytes), " ")

  status := string(read_bytes[0:index])

  if status == "ERR" {
    return nil, errors.New("ERR returned")
  }

  return url.ParseQuery(string(read_bytes[index:bytes_read-index]))
}

func (client MogileClient) GetPaths(key string, domain string) []string {
  response, err := client.RunCommand("get_paths", url.Values{ "key": []string{key},
                                                              "domain": []string{domain},
                                                              "noverify": []string{"1"} })
  if err != nil {
    return []string{}
  }

  fmt.Print("\n")
  fmt.Print(response)
  fmt.Print("\n")
  fmt.Print(response["unknown_key"])
  num_paths, _ := strconv.ParseInt(response["paths"][0], 10, 8)
  paths := make([]string, num_paths)

  index := 0
  for key, value := range response {
    if key == "paths" {
      continue
    }
    paths[index] = value[0]
    index += 1
  }

  return paths
}

func GetUrl(url string) {
  resp, err := http.Get(url)

  if err == nil {
    buffer, err := ioutil.ReadAll(resp.Body)
    if err == nil {
      fmt.Print(buffer)
    } else {
      fmt.Print(err)
    }
  } else {
    fmt.Print(err)
  }
}
