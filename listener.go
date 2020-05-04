package ping

import (
	"context"
	"fmt"
	"net"
	"syscall"
)

func ListenPacket(network, address string, vrf string) (net.PacketConn, error) {
	control := func(network, address string, c syscall.RawConn) (err error) {
		defer func() {
			if err != nil {
				fmt.Printf("ping Listener.control for vrf %s returning %s", vrf, err)
			}
		}()
		var syscallErr error
		controlErr := c.Control(func(fd uintptr) {
			syscallErr = syscall.BindToDevice(int(fd), vrf)
		})
		if syscallErr != nil {
			log.Errorf("syscall error setting sockopt TCP_USER_TIMEOUT: %v", syscallErr)
			return syscallErr
		}
		if controlErr != nil {
			log.Errorf("control error setting sockopt TCP_USER_TIMEOUT: %v", syscallErr)
			return controlErr
		}
		return nil
	}

	var lc net.ListenConfig
	if vrf != "" {
		lc.Control = control
	}
	return lc.ListenPacket(context.Background(), network, address)
}
