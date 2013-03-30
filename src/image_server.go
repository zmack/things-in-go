package main
import(
  "./mogile"
  "fmt"
  "net/http"
  "time"
  "log"
  "io/ioutil"
  "path/filepath"
)

type myHandler struct {
  missing_image []byte
  handler http.Handler
}

func main() {
  fmt.Print("Hello")
  missing, err := ioutil.ReadFile("/home/zmack/Pictures/smallgoat.jpg")
  if err != nil {
    return
  }

  server := &http.Server{
    Addr: ":3002",
    Handler: myHandler{
      missing_image: missing,
    },
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  log.Fatal(server.ListenAndServe())
}

func (handler myHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
  client, err := mogile.Connect("127.0.0.1:7001")
  key, _ := filepath.Split(request.URL.Path)

  if err != nil {
    writer.Header().Set("Content-Type", "image/jpeg")
    writer.Write(handler.missing_image)
    return
  }

  paths := client.GetPaths(key[1:len(key)-1], "foobar")

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
