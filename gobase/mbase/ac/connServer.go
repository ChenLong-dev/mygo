/*
 * @Description:
 * @Author: Chen Long
 * @Date: 2020-12-16 14:10:36
 * @LastEditTime: 2020-12-16 14:10:36
 * @LastEditors: Chen Long
 * @Reference:
 */

package ac

import (
	"net"

	"mbase/mcom"
	"mlog"
)

func DialConn(addr string) (conn *Conn, err error) {
	mlog.Tracef("addr=%v", addr)
	defer func() { mlog.Tracef("addr=%s,err=%v", addr, err) }()

	c, cerr := net.Dial("tcp", addr)
	if cerr != nil {
		return nil, cerr
	}

	tcpConn, ok := c.(*net.TCPConn)
	if ok {
		/*tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(time.Second * 20)*/
		mcom.SetTcpKeepAlive(tcpConn, 60, 30, 6)
		tcpConn.SetNoDelay(true)
	}

	return NewConn(c, true), nil
}
