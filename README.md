# Telegraf plugin: GoogleAnalytics


#### Plugin arguments:
- **key_file** string: Link to json file that is created using steps described in https://developers.google.com/analytics/devguides/reporting/core/v4/quickstart/installed-py.  
Create a private key clrty-secret.json and download it to your file system. Also add the new service account to the Google Analytics Account with Read & Analyze permission.
- **view_id** string: view id of what you are going to track

##### Configuration:
```
# Read metrics about google_analytics usage
[[inputs.google_analytics]]
  ## path to keyfile.
  key_file = "./plugins/inputs/google_analytics/clrty-secret.json"
  ## viewid
  view_id = "ga:155849743"
```

Result for above configuration.
```
{
	"gatype":"INTEGER", 
	"is_golden":false, 
	"max_val":"118", 
	"min_val":"1", 
	"row_count":20, 
	"totals":"239", 
	"row_metrics":{
		"1", "4", "7", "1", "62", "16", "1", "1", "1", "1", "3", "118", "1", "2", "2", "2", "2", "1", "10", "3"
	}, 
	"name":"ga:sessions"
} 
```



#### Description

The GoogleAnalytics input plugin collects data from google analytics api and convert to store on influxDB.

#### Testing Package

1. copy this plugin to  
`~/go/src/github.com/influxdata/telegraf/plugins/inputs/google_analytics`.   
2. add `_ "github.com/influxdata/telegraf/plugins/inputs/google_analytics"` at the end of `~/go/src/github.com/influxdata/telegraf/plugins/inputs/all/all.go`  
3. cd ~/go/src/github.com/influxdata/telegraf  
4. `make`  
5. copy `clrty-secret.json` to `~/go/src/github.com/influxdata/telegraf/plugins/inputs/google_analytics` folder.  
6. run below command.  
   `./telegraf --config telegraf.conf --test`   
7. check all goes well. :)
   In my case, I got below result.  
   `gatype,gatype=INTEGER,host=Luks-iMac.local gatype="INTEGER",is_golden=false,max_val="121",min_val="1",name="ga:sessions",row_count=20i,totals="244" 1534186610000000000`
