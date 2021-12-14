package main

import (
    "fmt"
    "strings"
)

// go:embed blacklist.txt
var blacklistFileContent string

var blacklist []string

func init() {
    blacklistLines := splitOnNewlines(blacklistFileContent)
    for _, line := range blacklistLines {
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }
        if line[0] == '#' {
            continue
        }
        key := strings.ToLower(line)
        blacklist = append(blacklist, key)
    }
}

func checkBlacklist(licenseId string) error {
    lowercased := strings.ToLower(licenseId)
    for _, key := range blacklist {
        if strings.Contains(lowercased, key) {
            return fmt.Errorf("license %s is blacklisted", licenseId)
        }
    }
    return nil
}
