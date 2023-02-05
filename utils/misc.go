package utils

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func FlattenTree(indexTree map[string][]string, parent string) []string {
	array := indexTree[parent]

	result := []string{}
	for _, child := range array { // for each child
		_, ok := indexTree[child]
		result = append(result, child)
		if ok {
			result = append(result, FlattenTree(indexTree, child)...)
		}
	}

	return result
}

func CoalesceError(args ...error) error {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}
	return nil
}

func IsWritable(location string) bool {
	fileName := fmt.Sprintf("smb-protect-%s-%s", uuid.New().String(), uuid.New().String())

	filePath := path.Join(location, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		errStr := err.Error()
		if !strings.Contains(errStr, "Access is denied.") {
			fmt.Println("Read-write check fails, with the following error: ", errStr)
		}
		return false
	}

	file.Close()
	os.Remove(filePath)

	return true
}
