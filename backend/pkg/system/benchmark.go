package system

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DiskBenchmark struct {
	Device     string    `json:"device"`
	ReadSpeed  float64   `json:"readSpeed"`
	WriteSpeed float64   `json:"writeSpeed"`
	ReadIOPS   int       `json:"readIOPS"`
	WriteIOPS  int       `json:"writeIOPS"`
	AccessTime float64   `json:"accessTime"`
	Timestamp  time.Time `json:"timestamp"`
}

// RunBenchmark runs a simple benchmark using hdparm for read speed
func RunBenchmark(device string) (*DiskBenchmark, error) {
	devPath := device
	if !strings.HasPrefix(device, "/dev/") {
		devPath = "/dev/" + device
	}

	benchmark := &DiskBenchmark{
		Device:    device,
		Timestamp: time.Now(),
	}

	// Read Speed using hdparm
	// -t: perform timings of device reads
	// --direct: use O_DIRECT to bypass page cache
	cmdRead := exec.Command("hdparm", "-t", "--direct", devPath)
	outputRead, err := cmdRead.Output()
	if err == nil {
		// Example output: Timing buffered disk reads:  864 MB in  3.00 seconds = 287.67 MB/sec
		re := regexp.MustCompile(`=\s+([\d\.]+)\s+MB/sec`)
		m := re.FindStringSubmatch(string(outputRead))
		if len(m) >= 2 {
			speed, _ := strconv.ParseFloat(m[1], 64)
			benchmark.ReadSpeed = speed * 1024 * 1024 // convert to bytes/sec
		}
	} else {
		// fallback to a simple dd read if hdparm fails or is not available
		cmdDD := exec.Command("dd", "if="+devPath, "of=/dev/null", "bs=1M", "count=1024", "status=progress")
		// dd output is to stderr
		outputDD, _ := cmdDD.CombinedOutput()
		reDD := regexp.MustCompile(`, ([\d\.]+)\s+([KMG]B/s)`)
		mDD := reDD.FindStringSubmatch(string(outputDD))
		if len(mDD) >= 3 {
			speed, _ := strconv.ParseFloat(mDD[1], 64)
			unit := mDD[2]
			switch unit {
			case "GB/s": speed *= 1024 * 1024 * 1024
			case "MB/s": speed *= 1024 * 1024
			case "KB/s": speed *= 1024
			}
			benchmark.ReadSpeed = speed
		}
	}

	return benchmark, nil
}
