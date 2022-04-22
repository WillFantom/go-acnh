# AC:NH API Client

Use the [AC:NH API](https://acnhapi.com) from Go!

---

## Supported Functions

 - **K.K.Slider Songs**: Search for and download K.K.Slider songs
 - **Background Music**: Search for and download BGM via hour, weather or both

---

```go
package main

import (
  acnh "github.com/willfantom/go-acnh"
)

const (
  downloadDir string = "~/Downloads"
)

func main() {
  client := acnh.New()

  allBGM, _ := client.BGMList()

  for _, track := range allBGM {
    client.BGMDownload(track, downloadDir)
  }
}

```