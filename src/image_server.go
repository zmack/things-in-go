package main
import(
  "./mogile"
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
    },
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  fmt.Printf("Up, listening on %s, connecting to %s.\n", settings["listen_address"], settings["tracker"])
  fmt.Printf("Domain is %s. Everything is horrible.\n", settings["domain"])
  log.Fatal(server.ListenAndServe())
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
  client, err := mogile.Connect(handler.tracker)
  key, _ := filepath.Split(request.URL.Path)

  if err != nil {
    writer.Header().Set("Content-Type", "image/jpeg")
    writer.Write(handler.missing_image)
    return
  }

  paths := client.GetPaths(key[1:len(key)-1], handler.domain)

  if len(paths) == 0 {
    writer.Header().Set("Content-Type", "image/jpeg")
    writer.Write(handler.missing_image)
  } else {
    spitFile(paths[0], writer)
  }
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
