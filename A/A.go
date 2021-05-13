package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)


//A server application calls the Upgrader.Upgrade method from an HTTP request handler to get a *Conn:
var Aaddr = flag.String("Aaddr", "localhost:8000", "http service address")
var Baddr = flag.String("Baddr", "localhost:8001", "http service address")
var Caddr = flag.String("Caddr", "localhost:8002", "http service address")
var Daddr = flag.String("Daddr", "localhost:8003", "http service address")
var Controller = flag.String("Controller", "localhost:8005", "http service address")

var controllerAddr Client
var clientvalue ="A"
var myAttributevalue =90

var routingTable=map[string]*string{
	"firsthop":Baddr,
	"nexthop":Caddr,
	"leader":nil,
	"controller":Controller,
}
var Leader string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//the client id which is auto generated id using autoid function and the connection variable stores the connection value
type Client struct {
	Id         string
	Connection *websocket.Conn
}
//the main Leader client value

//autoid generates a unique id to differentiate
func autoId() string {

	return uuid.Must(uuid.NewV4()).String()
}
//the sent method of client type is to sent a message to the websoccect connection
func (client *Client) Send(message [] byte) error {

	return client.Connection.WriteMessage(1, message)

}



//----------------------for passing data to channels creating struct (start).....
type QueueData struct{
	clinetdetails Client
	payloadmessage []byte
}
type DataPasser struct {
	queue chan QueueData
}

//this Function is used for websocket connections
func (p *DataPasser) ConnectNodes(w http.ResponseWriter, r *http.Request) {
	//Upgrader check orgin is returned true so that any websocket connections arent refsed
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true

	}

	fmt.Println("starting connec.......")
	//The Conn type represents a WebSocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//the client details are added to the client Val with struct Client
	ClientVal := Client{
		Id:         autoId(),
		Connection: conn,
	}


	//reading incoming data from the client in and infinte loop until the connection from the client is closed
	for {
		_, payload, err := conn.ReadMessage()
		fmt.Println("payload msg",payload)
		//the data where payload is the messafe and client val is the client details are passed to the channel queue
		p.queue<-QueueData{ClientVal,payload}


		if err != nil {
			log.Println(err)
			return
		}
	}


}
//......................... struct  to unmarshal JSON tupe----------------
type Message struct {
	Action  string          `json:"action"`
	AttributeValue string   `json:"attributevalue"`
	Client string 			`json:"client"`
}

//the log method is run as a go routine and check for data in queue channel and does the publish subscribe methods accordingly
//checking the message sent by the client
func ConnectToAddress(adressval string)(*websocket.Conn,error){
	u := url.URL{Scheme: "ws", Host: *routingTable[adressval], Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial to address failed trying address", err)
		return c,err

	}
	return c,nil
}
func (p *DataPasser) log() {
	//for range in the channel queue
	for item := range p.queue {
		//m of tpe Message used for unmarshalling the JSON
		m := Message{}
		//the message is unmarshaled to m
		err := json.Unmarshal(item.payloadmessage, &m)
		if err != nil {
			fmt.Println("This is not correct message payload",err,"payload : ",item.payloadmessage)
		}
		fmt.Println("Client Id", item.clinetdetails.Id)
		fmt.Println("Messsage :",m.AttributeValue,m.Client,m.Action)

		if m.Action=="controller"{
			controllerAddr=item.clinetdetails

			fmt.Println("------------------------------------------------")
			fmt.Println("controller connection establisheddd!!!!!!!!!!!")
			fmt.Println("------------------------------------------------")


			//if err = controllerAddr.Connection.WriteMessage(1, []byte("amango")); err != nil {
			//	fmt.Println(err)
			//}

			continue


		}



		c,err:=ConnectToAddress("firsthop")
		if err!=nil{
			routingTable["firsthop"]=Caddr
			routingTable["nexthop"]=Daddr
			c,err=ConnectToAddress("firsthop")
		}



		//for action mentioned in the action parameter publish and subscribe methods are done
		if m.Action =="election" {
			intattributevalue, err := strconv.Atoi(m.AttributeValue)
			if (m.Client==clientvalue) && (intattributevalue==myAttributevalue){
				fmt.Println(" IM THE LEADER")
				routingTable["leader"]=Aaddr
				Leader=m.Client
				if err = controllerAddr.Connection.WriteMessage(1, []byte("Leader")); err != nil {
					fmt.Println(err)
				}
			}else if intattributevalue>myAttributevalue{
				fmt.Println("attibute value is lower my atrribute value: ",myAttributevalue,"its attribute value:  ",intattributevalue)
				if err = c.WriteMessage(1, item.payloadmessage); err != nil {
					fmt.Println(err)
				}
				switch m.Client {
				case "A":
					routingTable["leader"]=Aaddr
					break
				case "B":
					routingTable["leader"]=Baddr
					break
				case "C":
					routingTable["leader"]=Caddr
					break
				case "D":
					routingTable["leader"]=Daddr
					break

				}
				controllerMessage:="passing message attribute  "+strconv.Itoa(intattributevalue)+"client :  "+m.Client
				if err = controllerAddr.Connection.WriteMessage(1, []byte(controllerMessage)); err != nil {
					fmt.Println(err)
				}

				Leader=m.Client

			}else if intattributevalue<myAttributevalue{
				fmt.Println("attibute value is lower my higher value: ",myAttributevalue,"its attribute value:  ",intattributevalue)
				msg:="{\"action\":\"election\",\"attributevalue\":\""+strconv.Itoa(myAttributevalue)+"\",\"client\":\""+clientvalue+"\"}"
				fmt.Println(msg)
				if err = c.WriteMessage(1, []byte(msg)); err != nil {
					fmt.Println(err)
				}
				routingTable["leader"]=Aaddr
				Leader=m.Client

				controllerMessage:="replacing message attribute with my attribute because mine is higher  "+strconv.Itoa(myAttributevalue)
				if err = controllerAddr.Connection.WriteMessage(1, []byte(controllerMessage)); err != nil {
					fmt.Println(err)
				}

			}
		} else if m.Action=="initiate"{
				fmt.Println("initiation happening...")
				msg:="{\"action\":\"election\",\"attributevalue\":\""+strconv.Itoa(myAttributevalue)+"\",\"client\":\""+clientvalue+"\"}"
				fmt.Println(msg)
				if err = c.WriteMessage(1, []byte(msg)); err != nil {
					fmt.Println(err)
				}
				routingTable["leader"]=Aaddr
				Leader=m.Client

				if err = controllerAddr.Connection.WriteMessage(1, []byte("initiating value A 90")); err != nil {
					fmt.Println(err)
				}


		}else if m.Action=="kill"{
			os.Exit(0)


		}





		}


	}






func main(){
	//the channel queue is used to pass data to the log funtion when new messages are recieved from the client
	queue:= make(chan QueueData)

	passer := &DataPasser{queue: queue}
	//log function sorts the request to publish and subscribe
	go passer.log()

	http.HandleFunc("/ws",passer.ConnectNodes)
	//addr=rand.Intn(10000 - 8000) + 8000
	err:=http.ListenAndServe(":8000",nil)

	if err !=nil{
		fmt.Println("error in listen and serve",err)
	}
}
