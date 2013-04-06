package main
import(
  "./mogile"
  "./lru_cache"
  "fmt"
  "net/http"
  "time"
  "log"
  "io/ioutil"
  "path/filepath"
  "launchpad.net/goyaml"
)

type myHandler struct {
  missing_image []byte
  domain string
  tracker string
  handler http.Handler
  cache *cache.LRUCache
}

type lruValue struct {
  value string
}

func main() {
  missing, err := ioutil.ReadFile("/home/zmack/Pictures/smallgoat.jpg")
  if err != nil {
    return
  }

  settings := getSettings("settings.yaml")

  server := &http.Server{
    Addr: settings["listen_address"],
    Handler: myHandler{
      missing_image: missing,
      domain: settings["domain"],
      tracker: settings["tracker"],
      cache: cache.NewLRUCache(1000),
    },
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  fmt.Printf("Up, listening on %s, connecting to %s.\n", settings["listen_address"], settings["tracker"])
  fmt.Printf("Domain is %s. Everything is horrible.\n", settings["domain"])
  log.Fatal(server.ListenAndServe())
}

func (value lruValue) Size() int {
  return len(value.value)
}

func (value lruValue) Value() string {
  return value.value
}

func getSettings(path string) map[string]string {
  bytes, err := ioutil.ReadFile(path)
  if err != nil {
    return map[string]string{ "domain": "defaultDomain" }
  }
  var settings map[string]string
  goyaml.Unmarshal(bytes, &settings)

  return settings
}

func (handler myHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
  path, _ := filepath.Split(request.URL.Path)
  mogile_key := path[1:len(path)-1]

  path, err := handler.getPathForKey(mogile_key)
  fmt.Printf("Got path %s\n", path)

  if len(path) == 0 || err != nil {
    writer.Header().Set("Content-Type", "image/jpeg")
    writer.Write(handler.missing_image)
  } else {
    spitFile(path, writer)
  }
}

func (handler myHandler) getPathForKey(mogile_key string) (value string, err error) {
  cached_value, ok := handler.cache.Get(mogile_key)

  if ok {
    lru_value := *cached_value.(*lruValue)
    return lru_value.value, nil
  }

  fmt.Printf("Cache miss for %s", mogile_key)

  client, err := mogile.Connect(handler.tracker)
  if err != nil {
    return "", err
  }

  paths := client.GetPaths(mogile_key, handler.domain)

  if len(paths) == 0 {
    handler.cache.Set(mogile_key, &lruValue{value: ""})
    return "", nil
  }

  handler.cache.Set(mogile_key, &lruValue{value: paths[0]})
  return paths[0], nil

}

func spitFile(url string, writer http.ResponseWriter) {
  resp, err := http.Get(url)
  if err != nil {
    return
  }

  buffer, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return
  }

  writer.Header().Set("Content-Type", http.DetectContentType(buffer))
  writer.Write(buffer)
}
