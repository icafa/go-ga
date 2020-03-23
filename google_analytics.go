package google_analytics

import (
   "io/ioutil"
   "net/http"

   "fmt"
   "time"
   "golang.org/x/oauth2"
   "golang.org/x/oauth2/google"
   ga "google.golang.org/api/analyticsreporting/v4"
   "github.com/influxdata/telegraf"
   "github.com/influxdata/telegraf/plugins/inputs"
)


type GoogleAnlayticsReport struct {
   KeyFile string `toml:"key_file"`
   ViewID string `toml:"view_id"`
}

func (_ *GoogleAnlayticsReport) Description() string {
   return "Read metrics about google analytics"
}

var sampleConfig = `  
## path to keyfile.
key_file = "./plugins/inputs/google_analytics/clrty-secret.json"
## viewid
view_id = "ga:155849743"
`

func (_ *GoogleAnlayticsReport) SampleConfig() string {
   return sampleConfig
}

var (
   keyfile = "./plugins/inputs/google_analytics/clrty-secret.json"
   viewID = "ga:155849743"
)

func makeReportSvc() (*ga.Service, error) {
   // defer TimeTrack(time.Now(), "Make reporting service")

   data, err := ioutil.ReadFile(keyfile)

   if err != nil {
      fmt.Println("Failed to load credentials for Google Analytics")
      return nil, err
   }

   fmt.Println("keyfile", keyfile)

   conf, err := google.JWTConfigFromJSON(data, ga.AnalyticsReadonlyScope)

   if err != nil {
      fmt.Println("Failed to create JWT config from JSON creds")
      return nil, err
   }

   fmt.Println("Created jwt config")

   // Initiate an http.Client. The following GET request will be
   // authorized and authenticated on the behalf of
   // your service account.
   var netClient *http.Client
   netClient = conf.Client(oauth2.NoContext)

   fmt.Println("Created authentication capable HTTP client")

   // Construct the Analytics Reporting service object.
   svc, err := ga.New(netClient)

   if err != nil {
      fmt.Println("Failed to create Google Analytics Reporting Service")
      return nil, err
   }

   fmt.Println("Created Google Analytics Reporting Service object")

   return svc, nil
}

func getReport(svc *ga.Service) (*ga.GetReportsResponse, error) {
//   defer TimeTrack(time.Now(), "GET Analytics Report")
   req := &ga.GetReportsRequest{
      ReportRequests: []*ga.ReportRequest{
         // Create the ReportRequest object.
         {
            ViewId: viewID,
            DateRanges: []*ga.DateRange{
               // Create the DateRange object.
               {StartDate: "7daysAgo", EndDate: "today"},
            },
            Metrics: []*ga.Metric{
               // Create the Metrics object.
               {Expression: "ga:sessions"},
            },
            Dimensions: []*ga.Dimension{
               {Name: "ga:country"},
            },
         },
      },
   }

   fmt.Println("Doing GET request from analytics reporting")

   return svc.Reports.BatchGet(req).Do()
}

func (s *GoogleAnlayticsReport) Gather(acc telegraf.Accumulator) error {

   keyfile = s.KeyFile
   viewID = s.ViewID
   svc, err := makeReportSvc()

   if err != nil {
     return fmt.Errorf("Error while creating Google Analytics Reporting Service: %s", err)
   }

   res, err := getReport(svc)

   if err != nil {
     return fmt.Errorf("GET request to analyticsreporting/v4 returned error: %s", err)
   }

   if res.HTTPStatusCode != 200 {
     return fmt.Errorf("Did not get expected HTTP response code.\n HTTPStatusCode: %d", res.HTTPStatusCode)
   }
   now := time.Now()

   for _, report := range res.Reports {
      header := report.ColumnHeader
      // dimHdrs := header.Dimensions
      metricHdrs := header.MetricHeader.MetricHeaderEntries
      rows := report.Data.Rows

      tags := map[string]string{
         "gatype": header.MetricHeader.MetricHeaderEntries[0].Type,
      }

      var rowMetrics = []string{}

      for _, row := range rows {

         metrics := row.Metrics

         for _, metric := range metrics {
            for j := 0; j < len(metricHdrs) && j < len(metric.Values); j++ {
               rowMetrics = append(rowMetrics, metric.Values[j])
            }
         }
      }
      // Add metrics
      fieldsR := map[string]interface{}{
         "name":       header.MetricHeader.MetricHeaderEntries[0].Name,
         "gatype":     header.MetricHeader.MetricHeaderEntries[0].Type,
         "is_golden":  report.Data.IsDataGolden,
         "max_val":    report.Data.Maximums[0].Values[0],
         "min_val":    report.Data.Minimums[0].Values[0],
         "row_count":  report.Data.RowCount,
         "totals":     report.Data.Totals[0].Values[0],
         "row_metrics":rowMetrics,
      }

      acc.AddCounter("gatype", fieldsR, tags, now)
   }

   return err
}

func init() {
   inputs.Add("google_analytics", func() telegraf.Input {
      return &GoogleAnlayticsReport{
         KeyFile: "./plugins/inputs/google_analytics/clrty-secret.json",
         ViewID: "ga:155849743",
      }
   })
}