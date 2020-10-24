package main

import (
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type hop struct {
	number  string
	name    string
	address string
	latency float64
}

func trace(host string) ([]hop, error) {
	var hops []hop
	// TODO: implement traceroute in pure go, rather than using exec
	cmd := exec.Command("traceroute", "-q 1", host)
	// run command
	if output, err := cmd.CombinedOutput(); err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Error("Error running traceroute")
		return nil, err
	} else {
		strOutput := strings.TrimSpace(string(output)) //trim to handle trailing \n
		for _, line := range strings.Split(strOutput, "\n")[1:] {
			if h, err := parseHop(line); err != nil {
				log.WithFields(log.Fields{"error": err.Error()}).Error("Error parsing hop")
				return hops, err
			} else {
				hops = append(hops, h)
			}
		}
	}
	return hops, nil
}

func parseHop(line string) (hop, error) {
	var h hop
	var err error

	values := strings.Fields(line)
	h.number = values[0]

	if h.name = values[1]; h.name == "*" {
		// no data for this hop, return
		return h, nil
	}

	h.address = values[2]

	if h.latency, err = strconv.ParseFloat(values[3], 64); err != nil {
		return h, err
	}

	return h, nil
}
