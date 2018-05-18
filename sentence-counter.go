package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	var (
		show         bool
		rm           bool
		showReverse  bool
		destFilePath string
	)

	flag.BoolVar(&show, "show", false, "show mode")
	flag.BoolVar(&rm, "rm", false, "remove the sentnce")
	flag.BoolVar(&showReverse, "show-reverse", false, "show mode")
	flag.StringVar(&destFilePath, "dest", "", "log file path ")

	flag.Parse()

	if destFilePath == "" {
		fmt.Fprintf(os.Stderr, "need log file path -dest [path/to/log]\n")
		return
	}

	if _, err := os.Stat(destFilePath); os.IsNotExist(err) {
		if _, err = os.Stat(filepath.Dir(destFilePath)); os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(destFilePath), 0744)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create log file %s \n", destFilePath)
				return
			}
		}
	}

	if show || showReverse {
		err := showLog(destFilePath, showReverse)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse log file %s \n", destFilePath)
			return
		}
	} else if rm {
		// remove
		sentence, err := readFromStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get input %+v \n", err)
			return
		}

		err = rmFromLog(sentence, destFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write file %+v \n", err)
			return
		}

	} else {
		sentence, err := readFromStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get input %+v \n", err)
			return
		}

		err = incLog(sentence, destFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write file %+v \n", err)
			return
		}
	}
}

func showLog(destFilePath string, reverse bool) error {
	log, err := readLog(destFilePath)
	if err != nil {
		return err
	}

	kvs := toKV(log)
	sort.Slice(kvs, func(i, j int) bool {
		if reverse {
			if kvs[i].Value <= kvs[i].Value {
				return true
			}
		} else {
			if kvs[i].Value >= kvs[i].Value {
				return true
			}
		}

		return false
	})

	for _, v := range kvs {
		fmt.Fprintf(os.Stdout, "%s\n", v.Key)
	}

	return nil
}

type KV struct {
	Key   string
	Value int
}

func toKV(m map[string]int) (result []KV) {
	for k, v := range m {
		result = append(result, KV{Key: k, Value: v})
	}
	return
}

func incLog(s string, destFilePath string) error {
	log, err := readLog(destFilePath)
	if err != nil {
		return err
	}
	if existsNumber, ok := log[s]; !ok {
		log[s] = 1
	} else {
		existsNumber += 1
		log[s] = existsNumber
	}
	return write(log, destFilePath)
}

func rmFromLog(s string, destFilePath string) error {
	log, err := readLog(destFilePath)
	if err != nil {
		return err
	}
	if _, ok := log[s]; ok {
		delete(log, s)
	}
	return write(log, destFilePath)
}

func readLog(destFilePath string) (map[string]int, error) {
	var log = map[string]int{}

	data, err := ioutil.ReadFile(destFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]int{}, nil
		}

		return nil, err
	}
	if len(data) == 0 {
		return map[string]int{}, nil
	}
	err = json.Unmarshal(data, &log)
	if err != nil {
		return nil, err
	}
	return log, nil

}

func write(log map[string]int, destFilePath string) error {
	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destFilePath, data, 0744)
	if err != nil {
		return err
	}

	return nil
}

func readFromStdin() (string, error) {
	r := bufio.NewReader(os.Stdin)
	input, _, err := r.ReadLine()
	if err != nil {
		return "", err
	}
	return string(input), nil
}
