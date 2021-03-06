// Package traceroute is a Go wrapper around the traceroute CLI utility
package traceroute

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Hop describes information about a single host along a given traceroute.
type Hop struct {
	Number  string
	Name    string
	Address string
	Latency float64
}

// Tracer provides acommon type for trace functions.
// They should receive a hostname as a string for input, and return a slice of
// Hops ([]Hop, error) that describe the hops along the route taken.
type Tracer func(string) ([]Hop, error)

// Trace is a Go wrapper around the CLI utility traceroute implementing Tracer.
// Given a hostname or address as a string it will return a slice of Hops
// describing the route taken to the remote host.
func Trace(host string) ([]Hop, error) {
	var hops []Hop
	// TODO: implement traceroute in pure go, rather than using exec
	cmd := exec.Command("traceroute", "-q", "1", host)
	// run command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error(), "output": string(output)}).Error("Error running traceroute")
		return nil, err
	}
	strOutput := strings.TrimSpace(string(output)) //trim to handle trailing \n
	for _, line := range strings.Split(strOutput, "\n")[1:] {
		h, err := parseHop(line)
		if err != nil {
			log.WithFields(log.Fields{"error": err.Error()}).Error("Error parsing hop")
			return hops, err
		}
		hops = append(hops, h)
	}
	return hops, nil
}

func parseHop(line string) (Hop, error) {
	var h Hop
	var err error

	values := strings.Fields(line)
	h.Number = values[0]

	if h.Name = values[1]; h.Name == "*" {
		// no data for this hop, return
		return h, nil
	}

	h.Address = values[2]

	if h.Latency, err = strconv.ParseFloat(values[3], 64); err != nil {
		return h, err
	}

	//scale latency based on units provided
	if values[4] == "ms" {
		h.Latency = h.Latency / 1000 // ms -> seconds
	} else {
		// TODO: Expand of the possibility of NS, Sec...
		return h, errors.New("Error parsing latency. Unknown units from traceroute")
	}

	return h, nil
}
