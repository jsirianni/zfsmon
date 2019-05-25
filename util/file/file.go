package file

import (
    "os"
    "fmt"
    "bufio"

)

// manageAlertFile manages the alert by by adding or removing an zpool
// action = 0 adds to file
// action = 1 removes from file
func ManageFile(line string, action int, filePath string) error {
    // get []string of alerted pool names
    fileLines, err := readFile(filePath)
    if err != nil {
        return err
    }

    found := -1
    for i, z := range fileLines {
        if z == line {
            found = i
        }
    }

    // if found and action is to add, do nothing
    if found != -1 && action == 0 {
        return nil
    }
    // if not found and action is to add, add to array
    if found == -1 && action == 0 {
        fileLines = append(fileLines, line)
    }

    // if not found and action to remove, do nothing
    if found == -1 && action == 1 {
        return nil
    }
    // if found and action is to remove, remove from array
    if found != -1 && action == 1 {
        fileLines = removeFromArray(fileLines, found)
    }

    // write array to file here and return the error / nil
    return writeFile(fileLines, filePath)
}

// poolAlerted returns true if a zpool has already been alered (found inside
// the alert file)
func PoolAlerted(name string, filePath string) bool {
    alertedZpools, err := readFile(filePath)
    if err != nil {
        fmt.Println("failed to read alert file")
        return false
    }

    for _, z := range alertedZpools {
        if z == name {
            return true
        }
    }
    return false
}

// readFile returns an array of strings from a text file
func readFile(filePath string) ([]string, error) {
    var a []string

    f, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0600)
    if err != nil {
        return a, err
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        a = append(a, scanner.Text())
    }
    return a, nil
}

// writeFile writes an array of strings to a text file
func writeFile(a []string, filePath string) error {
    f, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer f.Close()

    w := bufio.NewWriter(f)
    for _, alert := range a {
        fmt.Fprintln(w, alert)
    }
    return w.Flush()
}

func removeFromArray(s []string, i int) []string {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}
