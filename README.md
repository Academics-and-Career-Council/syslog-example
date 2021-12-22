# syslog-example


              ************        SYSLOG        ************
It is used to send logs, i.e., metadata and all info related to any relevant activity happening.

It has 3 sections to it, but essensially works by sending the logs to a RabbitMQ service, which is a service which sends 
messages which never fail, asynchronously. This is stored in a temporary database, which is then sent to mongoDB database periodically.

To send a request, make use of the code given in the example. 



              SOME GENERAL THINGS TO KEEP NOTE OF:
        1. The structure(JSON) of the data sent should be kept as such:
        type Syslog struct {
          ID			int			`bson:"id_for_ref"`       (This is used in sending process and need not be filled)
          ServiceName string		`bson:"service_name"`       (Fill in the name of service, eg, Kratos, Xenon, etc) ***** Has to be filled always!!! *****
          StatusCode	int			`bson:"status_code"`        (self explainatory)
          Severity	string		`bson:"severity"`          (Look at codes below, or refer to "syslog") 
          MsgName		string		`bson:"msg_name"`           (A heading for the message)
          Msg			string		`bson:"msg"`              (content of the message)
          InvokedBy	string		`bson:"invoked_by"`       (in case its a log of an action done by a user, then the users email ID)
          Result		string		`bson:"result"`         (success/failure)
          Batch      	int    		`bson:"batch"`        (Not needed to fill, this is for the reference of the syslog service)
          Timestamp	time.Time	`bson:"timestamp"`        (time when the error occured)
          CreatedAt 	time.Time   `bson:"createdAt,omitempty"`       (Not needed to fill)
        }
        
The severity levels should be one of the following:
      code        Name of Log         string to be filled (ensure correct case and spelling)
1.  	  0         Emergency	            emerg
2.  	  1         Alert	                alert
3.  	  2         Critical	            crit		
4.  	  3         Error	                err	
5.  	  4         Warning	              warning	
6.  	  5         Notice	              notice		
7.  	  6         Informational	        info		 
8.      7         Debug	                debug
