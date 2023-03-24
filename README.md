Install Browser WebSocket Client chrome extension to test 

To run the go program in server URL
1. Go to CMD
2. Run command - go run *.go

Testing from extension
1. Set this URL - ws://127.0.0.1:8083/ws on the server URL in extension and hit connect
2. Set this dummy JSON data in message and hit the send
    {
      "action": "abc",
      "username": "aqw",
      "message": "123"
    }
