package realtime

import (
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	log "github.com/sirupsen/logrus"
)

var SocketIO *socketio.Server

func StartSocketIO() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	//server.OnConnect("/", func(s socketio.Conn) error {
	//	fmt.Println("connected:", s.ID())
	//	return nil
	//})
	server.OnEvent("/", "join", func(s socketio.Conn, room string) {
		log.Infof("%s joining room %s", s.ID(), room)
		s.Join(room)

		if d, found := Current.GetMetaById(room); found {
			dd, _ := json.Marshal(d)
			s.Emit(fmt.Sprintf("%s:data", room), string(dd))
		}
	})
	//server.OnError("/", func(s socketio.Conn, e error) {
	//	//fmt.Println("meet error: ", e)
	//})
	//server.OnDisconnect("/", func(s socketio.Conn, reason string) {
	//	//fmt.Println("closed", reason)
	//})

	SocketIO = server
	//go server.Serve()
	//if err != nil {
	//	log.Fatal(err)
	//}

}
