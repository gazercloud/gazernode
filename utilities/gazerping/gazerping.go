package gazerping

import (
	"errors"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"sync"
	"time"
)

var sequenceNumber uint16
var mtx sync.Mutex
var defaultData map[int][]byte

func init() {
	sequenceNumber = 1
	defaultData = make(map[int][]byte)
	for s := 0; s < 1500; s++ {
		d := make([]byte, s)
		for i := 0; i < s; i++ {
			d[i] = byte(i%26) + 0x41
		}
		defaultData[s] = d
	}
}

func Ping(addr string, dataSize int, timeoutMs int) (result int, peer net.Addr, err error) {

	var data []byte
	data, _ = defaultData[dataSize]
	if data == nil {
		data = make([]byte, dataSize)
	}

	mtx.Lock()
	seqIndex := sequenceNumber
	srcIndex := uint16(os.Getpid() & 0xFFFF)
	sequenceNumber++
	if sequenceNumber > 65534 {
		sequenceNumber = 1
	}
	mtx.Unlock()

	return PingHost(addr, data, timeoutMs, srcIndex, seqIndex)
}

func PingHost(addr string, dataFrame []byte, timeoutMs int, source uint16, sequenceNum uint16) (result int, peer net.Addr, err error) {
	if len(dataFrame) < 1 || len(dataFrame) > 1400 {
		err = errors.New("wrong data frame length")
		return
	}

	if len(addr) < 1 {
		err = errors.New("wrong address")
		return
	}

	if timeoutMs < 1 {
		err = errors.New("wrong timeout")
		return
	}

	if timeoutMs > 10000 {
		err = errors.New("wrong timeout")
		return
	}

	var IPs []net.IP
	IPs, err = net.LookupIP(addr)
	if err != nil {
		return
	}

	var ipAddr net.IP

	for _, ip := range IPs {
		ipv4 := ip.To4()
		if ipv4 != nil {
			ipAddr = ipv4
		}
		//logger.Println("IPs:", ip.String())
	}

	if len(ipAddr) == 0 {
		err = errors.New("cannot lookup address")
	}

	var srv *icmp.PacketConn
	srv, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return
	}
	defer srv.Close()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   int(source),
			Seq:  int(sequenceNum),
			Data: dataFrame,
		},
	}
	var wb []byte
	wb, err = wm.Marshal(nil)
	if err != nil {
		return
	}
	if _, err = srv.WriteTo(wb, &net.IPAddr{IP: ipAddr}); err != nil {
		return
	}

	timeout := time.Millisecond * time.Duration(timeoutMs)

	err = srv.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return
	}

	rb := make([]byte, 1500)
	var n int
	t1 := time.Now()
	n, peer, err = srv.ReadFrom(rb)
	t2 := time.Now()
	if err != nil {
		return
	}
	var rm *icmp.Message
	rm, err = icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
	if err != nil {
		return
	}

	if rm.Type == ipv4.ICMPTypeDestinationUnreachable {
		err = errors.New("destination unreachable")
		return
	}

	if rm.Type != ipv4.ICMPTypeEchoReply {
		err = errors.New("error")
		return
	}

	result = int(t2.Sub(t1).Milliseconds())

	if result > int(timeout.Milliseconds()) {
		err = errors.New("timeout")
		return
	}

	return
}
