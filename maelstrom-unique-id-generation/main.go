package main
// a go program consists of packages where one of the packages are the main
import (
	"encoding/json"
	"log"
	"os"
	"context"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)
func main() {
	n := maelstrom.NewNode();

	//create a handle to be added to the network
	//this is like registering a handle callback
	//use a lock
	//and a number
	//start at 0
	

	n.Handle("generate", func(msg maelstrom.Message) error {
		body := make(map[string]any)
		var request map[string]any
		if err := json.Unmarshal(msg.Body, &request); err != nil {
				return err
		}
		//lock
	
		//call lin-tso service 
		//get the ts back
		unique_id_body := make(map[string]any)
		unique_id_body["type"] = "ts" // TYPE SHIT!!!!
		ts_msg, ts_err := n.SyncRPC(context.Background(), "lin-tso", unique_id_body)
		if ts_err != nil {
			return ts_err
		}
		var ts map[string]any
		if err := json.Unmarshal(ts_msg.Body, &ts); err != nil {
			return err
		}	
		body["id"] = ts["ts"]
		body["in_reply_to"] = request["msg_id"]
		body["type"] = "generate_ok"


		//increment
	
		//unlock
		
		//send out message
		
		

		return n.Reply(msg, body)
	})
	
	if err := n.Run(); err != nil {
		log.Printf("ERROR %s", err);
		os.Exit(1)
	}

	
}