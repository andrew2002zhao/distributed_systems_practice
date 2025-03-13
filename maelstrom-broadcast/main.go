// similar to namespaces in c++
package main

import (
	"encoding/json"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// need an integer store

func main() {
	n := maelstrom.NewNode();
	var self string;
	var neighbors[]interface{};
	n.Handle("topology", func(msg maelstrom.Message) error {
		request := make(map[string]any);
		if err := json.Unmarshal(msg.Body, &request); err != nil {
			return err;
		}
		self = msg.Dest;
		//access the value at topology and cast to a string array
		var topology map[string]any = request["topology"].(map[string]any);
		
		neighbors = topology[self].([]interface{});
		
		

		body := make(map[string]any);
		body["type"] = "topology_ok";
		return n.Reply(msg, body);
	});



	//make a set to hold all values seen
	values := make(map[float64]float64);

	n.Handle("broadcast_ok", func(msg maelstrom.Message) error {
		return nil;
	});
	
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		//need to unmarshall data
		var request map[string]any

		if err := json.Unmarshal(msg.Body, &request); err != nil {
			return err;
		}
		log.Printf("request message: ");
		// log.Printf("%f", request["message"].(float64));
		value := request["message"].(float64);
		// value_int := int(value_float);

		//propagate to all other nodes if not seen before		

		if _, ok := values[value]; !ok {
			//store the value for the future
			values[value] = value;
			//propagate
			for i := 0; i < len(neighbors); i++ {
				propagate_request := make(map[string]any);
				propagate_request["type"] = "broadcast";
				propagate_request["message"] = value;
				n.Send(neighbors[i].(string), propagate_request);
			} 
		} 

		body := make(map[string]any)
		body["type"] = "broadcast_ok";
		return n.Reply(msg, body);
	});

	n.Handle("read", func(msg maelstrom.Message) error {
		body := make(map[string]any);
		var values_array[]float64;
		for _, v := range values {
			// log.Printf("%d", v);
			values_array = append(values_array, v);
		}
		body["messages"] = values_array;
		body["type"] = "read_ok";
		return n.Reply(msg, body);
	});




	if err := n.Run(); err != nil {
		log.Printf("ERROR %s", err);
		os.Exit(1)
	}

}