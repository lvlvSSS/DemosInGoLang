package main

import (
	"github.com/sirupsen/logrus"
	_ "go_tools/log"
)

func main() {
	/*
		myserver := server.TCPServer{}
		myserver.SetHeartBeat(time.Second * 3600)

		myserver.AddLast("192.168.50.251", func(msg []byte) {
			fmt.Printf("Server receive msg from 192.168.50.251 : \r\n%#X \r\n", msg)
		})
		myserver.AddLast("192.168.50.7", func(msg []byte) {
			fmt.Printf("Server receive msg from 192.168.50.7 : \r\n%#X \r\n", msg)
		})
		myserver.AddLast("192.168.50.17", func(msg []byte) {
			fmt.Printf("Server receive msg from 192.168.50.17 : \r\n%#X \r\n", msg)
		})
		myserver.AddLast("192.168.50.37", func(msg []byte) {
			fmt.Printf("Server receive msg from 192.168.50.37 : \r\n%#X \r\n", msg)
		})
		myserver.Start(":50000")

		input := bufio.NewReader(os.Stdin)

		line, _, _ := input.ReadLine()
		for string(line) != "quit" {
			mybytes := []byte{0x02,0x00,0x00,0x10,0x03,0x70,0x73,0x03}
			myserver.Broadcast(mybytes)
			line, _, _ = input.ReadLine()
		}
	*/
	/*
		var exampleFormatter = &zt_formatter.ZtFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
			Formatter: nested.Formatter{
				//HideKeys: true,
				FieldsOrder: []string{"component", "category"},
			},
		}
		printDemo(exampleFormatter, "hello world")
	*/

	logrus.Infoln("testing")
}
func printDemo(f logrus.Formatter, title string) {
	l := logrus.New()

	l.SetLevel(logrus.DebugLevel)
	l.SetReportCaller(true)

	if f != nil {
		l.SetFormatter(f)
	}

	l.Infof("this is %v demo", title)

	lWebServer := l.WithField("component", "web-server")
	lWebServer.Info("starting...")

	lWebServerReq := lWebServer.WithFields(logrus.Fields{
		"req":   "GET /api/stats",
		"reqId": "#1",
	})

	lWebServerReq.Info("params: startYear=2048")
	lWebServerReq.Error("response: 400 Bad Request")

	lDbConnector := l.WithField("category", "db-connector")
	lDbConnector.Info("connecting to db on 10.10.10.13...")
	lDbConnector.Warn("connection took 10s")

	l.Info("demo end.")
}
