**Requirements**:

go get "github.com/gorilla/websocket"

go get "github.com/satori/go.uuid"

-run A.go

-run B.go

-run C.go

-run D.go

-run controller.go (open [http://localhost:8005/](http://localhost:8081/))

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.001.png)


Attribute values are the values given to the network and initially the network is connected using the ring formation: 

1. >B->C->D->A



Click on initiate leader election message : the clients will compare others attribute values will each other and select a leader among them transferring the messages in a ring fashion,

Intiate button will sent a “initiate message to the fuction and and election is started by A”

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.002.png)

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.003.png)

Each clients compares its attributed values with the values in the message and replaces its values if thier values are greater ;else pass on the message to the next client.

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.004.png)

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.005.png)

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.006.png)



each logs are also shown in the display card:

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.007.png)

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.008.png)

If kill leader is pressed os.exit is triggered in the leader connection and the leader exits

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.009.png)

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.010.png)

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.011.png)

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.012.png)


In client C since trying to connect to D fails it should connect to A by changing its first hop to A but

When error is detected te program exits . its not comparing its error and moving on so I tried using panic and recover but still after panic the program exits I am trying to fix that error

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.013.png)


Panic error :

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.014.png)

The fuction code that should be called after the leader is killed :(that reroutes the first hop to A)

![](Aspose.Words.bf06d843-260b-4a30-ae02-6efbb3f2dae7.015.png)

Communication is established using websockets importing 

"github.com/gorilla/websocket"

Unique id to the clients are given using uuid importing

"github.com/satori/go.uuid"

Producer uses Dial function to connect to broker 

websocket.DefaultDialer.Dial(u.String(), nil)

Subscriber uses javascript 

conn = new WebSocket("ws://localhost:8000/ws"); 



Improvements that can be made :

After fixing the panic I should also include a probable leader value in the message so the while selecting the initial leader itself a probable leader can be found .

And in the event where the leader fails the probable leader will take charge.

The selected probable leader and current leader should exchange heartbeat to know if the any of them fails .if either of them fails a election should be performed to determing the next probable or current leader.



