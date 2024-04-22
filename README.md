Example:
curl -X POST http://localhost:8080/create -H "X-Server: http://mein.prometheus.server:9090" -H "X-Header-Content-Type: application/json" -d '{"calcFunc":"if queryResults[0]+queryResults[1]/5 < 5 { return \"green\" } else if queryResults[0]+queryResults[1]/5 >= 5 { return \"orange\" } else if queryResults[0]+queryResults[1]/5 > 5 { return \"red\" } else if queryResults[2] == 1 { return \"maintenance\" } else { return \"unknown\" }", "queries":["up", "process_cpu_seconds_total"]}'