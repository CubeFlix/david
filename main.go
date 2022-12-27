// david virtual desktop

package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"net"
	"os"

	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kbinani/screenshot"
	"github.com/lxn/win"
)

type writer struct {
	Conn io.Writer
}

// write to the connection
func (c *writer) Write(b []byte) (int, error) {
	start := 0
	for {
		n, err := c.Conn.Write(b[start:])
		if err != nil {
			return start + n, err
		}
		start += n
		if start == len(b) {
			break
		}
	}
	return start, nil
}

type reader struct {
	Conn io.Reader
}

// read to the connection
func (c *reader) Read(b []byte) (int, error) {
	start := 0
	for {
		n, err := c.Conn.Read(b[start:])
		if err != nil {
			return start + n, err
		}
		start += n
		if start == len(b) {
			break
		}
	}
	return start, nil
}

// david server
func server(host string) {
	l, err := net.Listen("tcp", host+":42069")
	if err != nil {
		fmt.Println("failed to listen:", err.Error())
		return
	}

	defer l.Close()
	fmt.Println("waiting for client on", host+":42069")

	l2, err := net.Listen("tcp", host+":42070")
	if err != nil {
		fmt.Println("failed to listen:", err.Error())
		return
	}

	defer l2.Close()
	fmt.Println("waiting for client on", host+":42070")

	netConn, err := l.Accept()
	conn := writer{Conn: netConn}
	if err != nil {
		fmt.Println("failed to accept connection:", err.Error())
		return
	}
	fmt.Println("connected to", netConn.LocalAddr().String())

	netConn2, err := l2.Accept()
	conn2 := reader{Conn: netConn2}
	if err != nil {
		fmt.Println("failed to accept connection:", err.Error())
		return
	}
	fmt.Println("connected to", netConn2.LocalAddr().String())

	// start the screen capture
	sizeInfoBuf := make([]byte, 8)
	w := uint32(win.GetSystemMetrics(win.SM_CXSCREEN))
	h := uint32(win.GetSystemMetrics(win.SM_CYSCREEN))
	binary.LittleEndian.PutUint32(sizeInfoBuf, w)
	binary.LittleEndian.PutUint32(sizeInfoBuf[4:], h)
	if _, err := conn.Write(sizeInfoBuf); err != nil {
		fmt.Println("failed to write:", err.Error())
		return
	}
	//lastCapture := image.NewRGBA(image.Rect(0, 0, int(win.GetSystemMetrics(win.SM_CXSCREEN)), int(win.GetSystemMetrics(win.SM_CYSCREEN))))
	go func() {
		for {
			//now := time.Now()
			img, err := screenshot.Capture(0, 0, int(w), int(h))
			if err != nil {
				fmt.Println("failed to capture screen:", err.Error())
				return
			}
			//if !bytes.Equal(img.Pix, lastCapture.Pix) {
			//fmt.Println("capture", time.Now().Sub(now))
			//now = time.Now()
			//w := zlib.NewWriter(netConn)
			_, err = conn.Write(img.Pix)
			if err != nil {
				fmt.Println("failed to write:", err.Error())
				return
			}
			//w.Close()
			//fmt.Println("write", time.Now().Sub(now))
			//copy(lastCapture.Pix, img.Pix)
			//}
		}
	}()
	for {
		fmt.Println("waiting...")
		reqLenBuf := make([]byte, 4)
		_, err := conn2.Read(reqLenBuf)
		if err != nil {
			fmt.Println("failed to read:", err.Error())
			return
		}
		reqLen := binary.LittleEndian.Uint32(reqLenBuf)
		fmt.Println(reqLen)
		req := make([]byte, reqLen)
		_, err = conn2.Read(req)
		if err != nil {
			fmt.Println("failed to read:", err.Error())
			return
		}
		fmt.Println(req)
		reqType := string(req[:2])
		if reqType == "kd" {
			// Keyboard down.
			fmt.Println("kd", req[2:])
		} else if reqType == "ku" {
			// keyboard up
			fmt.Println("ku", req[2:])
		} else if reqType == "mv" {
			// Set mouse pos.
			fmt.Println("mv", req[2:])
		}
	}
}

var keys = []ebiten.Key{
	ebiten.KeyA,
	ebiten.KeyB,
	ebiten.KeyC,
	ebiten.KeyD,
	ebiten.KeyE,
	ebiten.KeyF,
	ebiten.KeyG,
	ebiten.KeyH,
	ebiten.KeyI,
	ebiten.KeyJ,
	ebiten.KeyK,
	ebiten.KeyL,
	ebiten.KeyM,
	ebiten.KeyN,
	ebiten.KeyO,
	ebiten.KeyP,
	ebiten.KeyQ,
	ebiten.KeyR,
	ebiten.KeyS,
	ebiten.KeyT,
	ebiten.KeyU,
	ebiten.KeyV,
	ebiten.KeyW,
	ebiten.KeyX,
	ebiten.KeyY,
	ebiten.KeyZ,
	ebiten.KeyAltLeft,
	ebiten.KeyAltRight,
	ebiten.KeyArrowDown,
	ebiten.KeyArrowLeft,
	ebiten.KeyArrowRight,
	ebiten.KeyArrowUp,
	ebiten.KeyBackquote,
	ebiten.KeyBackslash,
	ebiten.KeyBackspace,
	ebiten.KeyBracketLeft,
	ebiten.KeyBracketRight,
	ebiten.KeyCapsLock,
	ebiten.KeyComma,
	ebiten.KeyContextMenu,
	ebiten.KeyControlLeft,
	ebiten.KeyControlRight,
	ebiten.KeyDelete,
	ebiten.KeyDigit0,
	ebiten.KeyDigit1,
	ebiten.KeyDigit2,
	ebiten.KeyDigit3,
	ebiten.KeyDigit4,
	ebiten.KeyDigit5,
	ebiten.KeyDigit6,
	ebiten.KeyDigit7,
	ebiten.KeyDigit8,
	ebiten.KeyDigit9,
	ebiten.KeyEnd,
	ebiten.KeyEnter,
	ebiten.KeyEqual,
	ebiten.KeyEscape,
	ebiten.KeyF1,
	ebiten.KeyF2,
	ebiten.KeyF3,
	ebiten.KeyF4,
	ebiten.KeyF5,
	ebiten.KeyF6,
	ebiten.KeyF7,
	ebiten.KeyF8,
	ebiten.KeyF9,
	ebiten.KeyF10,
	ebiten.KeyF11,
	ebiten.KeyF12,
	ebiten.KeyHome,
	ebiten.KeyInsert,
	ebiten.KeyMetaLeft,
	ebiten.KeyMetaRight,
	ebiten.KeyMinus,
	ebiten.KeyNumLock,
	ebiten.KeyNumpad0,
	ebiten.KeyNumpad1,
	ebiten.KeyNumpad2,
	ebiten.KeyNumpad3,
	ebiten.KeyNumpad4,
	ebiten.KeyNumpad5,
	ebiten.KeyNumpad6,
	ebiten.KeyNumpad7,
	ebiten.KeyNumpad8,
	ebiten.KeyNumpad9,
	ebiten.KeyNumpadAdd,
	ebiten.KeyNumpadDecimal,
	ebiten.KeyNumpadDivide,
	ebiten.KeyNumpadEnter,
	ebiten.KeyNumpadEqual,
	ebiten.KeyNumpadMultiply,
	ebiten.KeyNumpadSubtract,
	ebiten.KeyPageDown,
	ebiten.KeyPageUp,
	ebiten.KeyPause,
	ebiten.KeyPeriod,
	ebiten.KeyPrintScreen,
	ebiten.KeyQuote,
	ebiten.KeyScrollLock,
	ebiten.KeySemicolon,
	ebiten.KeyShiftLeft,
	ebiten.KeyShiftRight,
	ebiten.KeySlash,
	ebiten.KeySpace,
	ebiten.KeyTab,
}

type Game struct {
	img     *image.RGBA
	netConn net.Conn
	conn2   writer
	// pressedKeys []ebiten.Key
}

func (g *Game) Update() error {
	// get key input
	// oldPressedKeys := make([]ebiten.Key, len(g.pressedKeys))
	// copy(oldPressedKeys, g.pressedKeys)
	// g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])
	// // compare the
	lenBuf := make([]byte, 4)
	for i := range keys {
		if inpututil.IsKeyJustPressed(keys[i]) {
			payload := []byte("kd" + keys[i].String())
			binary.LittleEndian.PutUint32(lenBuf, uint32(len(payload)))
			_, err := g.conn2.Write(lenBuf)
			if err != nil {
				fmt.Println("failed to send input:", err)
				return err
			}
			_, err = g.conn2.Write(payload)
			if err != nil {
				fmt.Println("failed to send input:", err)
				return err
			}
		}
		if inpututil.IsKeyJustReleased(keys[i]) {
			payload := []byte("ku" + keys[i].String())
			binary.LittleEndian.PutUint32(lenBuf, uint32(len(payload)))
			_, err := g.conn2.Write(lenBuf)
			if err != nil {
				fmt.Println("failed to send input:", err)
				return err
			}
			_, err = g.conn2.Write(payload)
			if err != nil {
				fmt.Println("failed to send input:", err)
				return err
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.img.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.img.Rect.Dx(), g.img.Rect.Dy()
}

func client(host string) {
	// connect to the client
	tcpAddr, err := net.ResolveTCPAddr("tcp", host+":42069")
	if err != nil {
		fmt.Println("failed to resolve address:", err.Error())
		return
	}

	netConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("tcp dial failed:", err.Error())
		return
	}
	conn := reader{Conn: netConn}

	tcpAddr2, err := net.ResolveTCPAddr("tcp", host+":42070")
	if err != nil {
		fmt.Println("failed to resolve address:", err.Error())
		return
	}

	netConn2, err := net.DialTCP("tcp", nil, tcpAddr2)
	if err != nil {
		fmt.Println("tcp dial failed:", err.Error())
		return
	}
	conn2 := writer{Conn: netConn2}

	// get the width and height info
	sizeInfoBuf := make([]byte, 8)
	if _, err := conn.Read(sizeInfoBuf); err != nil {
		fmt.Println("failed to read:", err.Error())
		return
	}
	sizeX := int(binary.LittleEndian.Uint32(sizeInfoBuf))
	sizeY := int(binary.LittleEndian.Uint32(sizeInfoBuf[4:]))

	fmt.Printf("connected to host: sizex=%d, sizey=%d\n", sizeX, sizeY)

	// start the screen
	img := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))
	ebiten.SetWindowSize(sizeX, sizeY)
	ebiten.SetWindowTitle("david")
	g := &Game{
		img:     img,
		netConn: netConn,
		conn2:   conn2,
	}
	go func() {
		for {
			//zlibReader, err := zlib.NewReader(g.netConn)
			//if err != nil {
			//	fmt.Println("failed to decompress:", err)
			//	os.Exit(1)
			//}
			//_, err = io.Copy(&g.tempPixBuf, zlibReader)
			//if err != nil {
			//	fmt.Println("failed to decompress:", err)
			//	os.Exit(1)
			//}
			//zlibReader.Close()
			conn.Read(g.img.Pix)
		}
	}()
	if err := ebiten.RunGame(g); err != nil {
		fmt.Println("error running window:", err.Error())
	}
}

// main
func main() {
	// check arguments
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Println("usage: david <command>")
		return
	}
	if os.Args[1] == "help" {
		// help
		fmt.Println("david is a virtual desktop server and client.")
		fmt.Println("usage: david <command>")
		fmt.Println("available commands:")
		fmt.Println("	- help: display this message")
		fmt.Println("	- server: run the david server with an optional host (defaults to 127.0.0.1)")
		fmt.Println("	- client: run the david client with an optional host (defaults to 127.0.0.1)")
	} else if os.Args[1] == "server" {
		// server
		host := "127.0.0.1"
		if len(os.Args) == 3 {
			host = os.Args[2]
		}
		server(host)
	} else if os.Args[1] == "client" {
		// client
		host := "127.0.0.1"
		if len(os.Args) == 3 {
			host = os.Args[2]
		}
		client(host)
	}
}
