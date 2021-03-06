package main

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "sort"
    "strings"
)

type License struct {
    Name string `json:"name"`
    LicenseId string `json:"licenseId"`
    LicenseText string `json:"license"`
    SPDX bool `json:"spdx"` // for consistency with licensebot-client-js
}

type LicenseWithoutText struct {
    Name string
    LicenseId string
    LicenseUrl string
}

func init() {
    // install go-licenses
    if _, err := runCommand("go", "install", "github.com/google/go-licenses@latest"); err != nil {
        log.Fatalln(err)
    }
}

func getModuleLicenses(relativePath string) ([]LicenseWithoutText, error) {
    stdout, err := runCommand("go-licenses", "csv", relativePath)
    if err != nil {
        return nil, err
    }

    r := csv.NewReader(strings.NewReader(stdout))
    var entries []LicenseWithoutText
    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
        dependency := record[0]
        if strings.HasPrefix(dependency, "mosaicmfg.com") {
            // skip our own subpackages
            continue
        }
        licenseUrl := record[1]
        licenseId := record[2]
        entries = append(entries, LicenseWithoutText{
            Name:       dependency,
            LicenseId:  licenseId,
            LicenseUrl: licenseUrl,
        })
    }
    return entries, nil
}

func getAllRepoModules() ([]License, error) {
    licensesMap := make(map[string]LicenseWithoutText)

    moduleLicenses, err := getModuleLicenses(".")
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    // add licenses to map where dependency is not already present
    for _, license := range moduleLicenses {
        if _, exists := licensesMap[license.Name]; !exists {
            licensesMap[license.Name] = license
        }
    }

    // check all entries against the blacklist before resolving license text
    blacklistErrors := 0
    for _, license := range licensesMap {
        if err := checkBlacklist(license.LicenseId); err != nil {
            fmt.Println(err)
            blacklistErrors++
        }
    }
    if blacklistErrors > 0 {
        return nil, errors.New("failed blacklist check")
    }

    // convert licenses map to a slice including text content
    licenses := make([]License, 0, len(licensesMap))
    for _, license := range licensesMap {
       licenseUrl := license.LicenseUrl
       if strings.HasPrefix(licenseUrl, "https://github.com") {
           licenseUrl = getGitHubRawUrl(licenseUrl)
       }

       // retrieve license text from URL
       resp, err := http.Get(licenseUrl)
       if err != nil {
           return nil, err
       }
       bodyBytes, err := io.ReadAll(resp.Body)
       if err != nil {
           return nil, err
       }
       licenseText := normalizeNewlines(string(bodyBytes))
       if err := resp.Body.Close(); err != nil {
           return nil, err
       }

       // normalize license name from GitHub URL to repository name
       licenseName := license.Name
       if strings.HasPrefix(licenseName, "github.com") {
           // github.com/<user-or-organization>/<repository>/<...path>
           nameParts := strings.Split(licenseName, "/")
           licenseName = nameParts[2]
       }

       licenses = append(licenses, License{
           Name:        licenseName,
           LicenseId:   license.LicenseId,
           LicenseText: licenseText,
       })
    }
    // sort the slice alphabetically
    sort.Slice(licenses, func(i, j int) bool {
       return licenses[i].Name < licenses[j].Name
    })

    return licenses, nil
}
