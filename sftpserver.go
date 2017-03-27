package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"testing"

	"golang.org/x/crypto/ssh"
)

var sshServerDebugStream = ioutil.Discard

var (
	hostPrivateKeySigner ssh.Signer
	privKey              = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAsSxU4LdH6sG/S5dlblfDntulWRI6m+weak1uBNNqriSSkG4g
KohMCbHmJJY6a1cAv53GFCpWJoq7XYvDgPFb3XKqb7riAr65RFtFzM8sQsVQoqV8
fknbIH16E72yMT5FceIbSM4/O3itOSlUY2vJ7BsEUjx7nDWaTP+OMFkwremo72Tq
tUNgins3X8m00Y5sKRA5qLVHwR04ZbFy9nGN3hKEPH6CVMejqJJLOeV81TX/XBni
CT6qMCmxYCdbgmGQROnDHUHSWNv17F0wb/0RMlQl0ymFHOSegtmdG2Xt9oYWYiDN
Gy0K8GBoCtHVLZt3fL1vafksS/sM5RYRUAFSCQIDAQABAoIBAQCeDrCB8MBF3Cau
ZxfkAoPP2p9+ANcsds8DgqQdxgYr6RCfrL8hcopzM7Pe++6OCAXw6+3j24kTxTw1
zhPRmoCb5EnMd2pdjIx3QP3aIxCXWLQBBaU0fOrx5z7bEaZAbA9D87TnlKewhI30
qrxQHb771XZbbv3Pc7p96paM52SYIJlvyUwQKCQNAZNqRDanR+efUWUlGWAXdY+x
/cG2U/EKA/8HQNXCF2TDNGb97MN8/Bl0VRYpxdBqHQl2ozXkpgPa1vXNQEQ0cHVT
JNYkI2vOQ3Iya3T1b6zx3opfiP7Y64ydpZu/mrU5oWMyikmdyX6wCpoMa0oyzRDZ
ZTPYm1gJAoGBAN5Uo69O5Up/C9qIp8UtfnH9XFTRsTnfh6SlU9xpuKjbN5Kw1sHb
DeXSt+Jzb3qeUMdg+98QoHAIy4oxsWjx9yfSjcm+EPpPtgJ8Wko+DJKQKmulo0/g
epKFo7FSoWGqtH8fr3zf7/S8+bphDQHYlWRmH/mIsLpyBovLTvaJ1eBLAoGBAMwB
BLu/iF+0l1rlnl9MN1rB7+5xxXGPoH4LUjI8PHCMimPW4WVwqCLCgNtlqx1zVK//
L92MOuAn9PjLvw8sspC56Tdzba7HAIPLHK8AtU4IP8+be0Oexctsj57IGKIcXI+y
YqFTptpVvEkPej+P9riX23qPh0Zuwtf2K69Fbup7AoGAW1FcYdb/6pdAISRb9Gr5
Moyj7dqq9mBPcFrPlQp/ZCuWKdQkgT8d+DWSfZp4QV7hQuMc0MQdgaa7IynB+p7X
qy2aOzCr/IPc+CxnUXMm6tP3+HryFw7WiXQGhgCwdFMPC9/RznKUNmugDuNp2kZB
JhmkLHPuUsYe1jBNYInApP0CgYEAwWrVygw2iEb4mb3LAh+I/AuUKEbGJH1AdUDW
lbp2s18Mdsxst3iwcQRol5s1OZ73VEZmY29pAs3ffWPvqbt/MaiSbXiLLYKQAmS4
tVO+klVP6s5HeD042z36jVi5wjmRqMxApyRgtfFDqyF5jno4OZwBA5rBbw3kvk0v
7eWu27ECgYBbQWJs4Z2w0HyrstJsInXlbRlR4J7L9JoE5h6sPZ3lOd+SN1ABgC5z
VnCc5jd++chnB4rEN9wRCr5NBSc/yqYx18kCI8eiN+OjhcmwwYq3glur7mHgQ/YI
33O4je+nAivvpHKW+rlB+K0NaESt547Iio41CD/03Q7P0Q5ufkuhQA==
-----END RSA PRIVATE KEY-----
`)
)

func init() {
	var err error
	hostPrivateKeySigner, err = ssh.ParsePrivateKey(privKey)
	if err != nil {
		panic(err)
	}
}

func getTestServer(t *testing.T) (net.Listener, string, int) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	host, portStr, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Fprintf(sshServerDebugStream, "ssh server socket closed: %v\n", err)
				break
			}

			go func() {
				defer conn.Close()
				sshSvr, err := sshServerFromConn(conn, useSubsystem, basicServerConfig())
				if err != nil {
					t.Error(err)
					return
				}
				err = sshSvr.Wait()
				fmt.Fprintf(sshServerDebugStream, "ssh server finished, err: %v\n", err)
			}()
		}
	}()

	return listener, host, port
}
