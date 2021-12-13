package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "path"
)

func getLicenseJSON() (string, error) {
  licenses, err := getAllRepoModules()
  if err != nil {
    return "", err
  }
  jsonBytes, err := json.MarshalIndent(licenses, "", "  ")
  if err != nil {
    return "", err
  }
  return string(jsonBytes), nil
}

func help() {
  fmt.Println("licensebot <command> [bundlePath]")
  fmt.Println("  command     command to run (\"update\" or \"check\")")
  fmt.Println("  bundlePath  JSON file to update or compare against")
  fmt.Println("                (defaults to $CWD/licenses.json)")
}

func update(bundlePath string) {
  licenses, err := getLicenseJSON()
  if err != nil {
    log.Fatalln(err)
  }
  if err := ioutil.WriteFile(bundlePath, []byte(licenses), 0644); err != nil {
    log.Fatalln(err)
  }
}

func check(bundlePath string) {
  licenses, err := getLicenseJSON()
  if err != nil {
    log.Fatalln(err)
  }
  fromDiskBytes, err := ioutil.ReadFile(bundlePath)
  fromDisk := string(fromDiskBytes)
  if licenses != fromDisk {
    log.Fatalln("License bundle is out of date.\nRun `licensebot update` and commit the changes.")
  }
}

func main() {
  if len(os.Args) == 1 {
    log.Fatalln("expected command as argument")
  }
  command := os.Args[1]

  cwd, err := os.Getwd()
  if err != nil {
    log.Fatalln(err)
  }
  bundlePath := path.Join(cwd, "licenses.json")
  if len(os.Args) >= 3 {
    bundlePath = os.Args[2]
  }

  switch command {
  case "help":
    fallthrough
  case "--help":
    fallthrough
  case "-h":
    help()
  case "update":
    update(bundlePath)
  case "check":
    check(bundlePath)
  default:
    help()
    log.Fatalf("unexpected command '%s' as argument", command)
  }
}
