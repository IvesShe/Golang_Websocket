package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
)

// 定義flag參數，這邊會返回一個相應的指針
var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	// 開啟websocket服務
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Debug("upgrade:", err)
		return
	}
	logger.Debug("服務器啟動，監聽端口8080中......")

	// 預先關閉，此行在離開echo時會執行
	defer ws.Close()

	// 一直待命收資料
	for {
		// 讀取資料
		mt, message, err := ws.ReadMessage()
		if err != nil {
			logger.Debug("read:", err)
			break
		}

		logger.Debug("recv: %s", message)
		// 如果是ping
		if string(message) == "ping" {
			// 就回pong
			mes := "pong"
			logger.Debug("心跳")
			// 將讀到的資料送出
			err = ws.WriteMessage(mt, []byte(mes))
			logger.Debug("write:", string(mes))
		} else {
			// 將讀到的資料送出
			err = ws.WriteMessage(mt, message)
			logger.Debug("write:", string(message))
		}
		// 將讀到的資料送出
		if err != nil {
			logger.Debug("write:", err)
			break
		}
	}
}

// 這邊會產生一個實體網頁
func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	logger.Debug("程序啟動")

	// 調用flag.Parse()解析命令行參數到定義的flag
	flag.Parse()

	// SetFlags(flag int)可以用來自定義log的輸出格式
	log.SetFlags(0)

	// 處理/echo的路由
	http.HandleFunc("/echo", echo)

	// 處理/的路由
	http.HandleFunc("/", home)

	// 啟動服務，addr是指針，所以加*取值，若有錯會回傳錯誤log
	log.Fatal(http.ListenAndServe(*addr, nil))

}

// 網頁html
var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
