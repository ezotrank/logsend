Getstarted:

```
make
./vendor/bin/logsend -log-dir=./tmp -send-buffer=1 -db-host="localhost:4444" -debug -database="test1" -udp -config=./config.json
```
Examples:

example.config.json

```
{
    "groups": [
        {
            "mask": "bee.[0-9]+.log$",
            "rules": [
                {
                    "columns": [
                        ["gate"],
                        ["exec_time","float"]
                    ],
                    "name": "gate_response",
                    "regexp": "chain_processor.+ .+MIRROR] ([a-zA-Z0-9_]+)_gate -> (.+)"
                },
                {
                    "columns": [
                        ["gate_id", "int"],
                        ["host"]
                    ],
                    "name": "clicks",
                    "regexp": "deeplink.+] {.+ \"gate_id\": ([0-9]+), .+ \"host\": \"([a-z.]*)\""
                },
                {
                    "columns": [
                        ["code","int"],
                        ["adaptor"],
                        ["resp_time","float"],
                        ["host","GetHostname"]
                    ],
                    "name": "adaptors_response_time",
                    "regexp": "web.+] ([0-9]+) .+ /adaptors/([a-zA-Z0-9_/]+) \\(.+\\) ([0-9]+.[0-9]+)ms$"
                }
            ]
        }
    ]
}
```

series

```
[
  {
    "name" : "response_times",
    "columns" : ["gate", "exec_time"],
    "points" : [
      [10, 123.33]
    ]
  }
]
```

```
[
  {
    "name" : "clicks",
    "columns" : ["gate_id", "adaptor", "resp_time", "host"],
    "points" : [
      [200, "check_something", 23.313, "test1.example.com"]
    ]
  }
]
```

```
[
  {
    "name" : "adaptors_response_time",
    "columns" : ["code", "adaptor"],
    "points" : [
      [10, "test1.example.com"]
    ]
  }
]
```



todo

- config checker
- add new configs