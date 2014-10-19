import os
import time

lines_count = 250000
log_dir = './bench_logs'
log_file_name = log_dir + '/test.log'
config_file = log_dir + '/config.json'
logsend_binary = 'logsend'
msg = "test string one\n"
binary="go run main.go"
run_params = "-config=%s -dry-run -read-whole-log -read-once %s &>/dev/null" % (config_file, log_dir)

config="""
{
  "influxdb": {
    "host": "localhost:4444",
    "user": "root",
    "password": "root",
    "database": "logers",
    "udp": true,
    "send_buffer": 8
  },
  "groups": [
    {
        "mask": "test.log",
        "rules": [
            {
                "regexp": "test string (?P<word_STRING>\\\\w+)",
                "influxdb": {
                    "name": "test"
                }
            }
        ]
    }
  ]
}
"""

config2="""
{
  "influxdb": {
    "host": "localhost:4444",
    "user": "root",
    "password": "root",
    "database": "logers",
    "udp": true,
    "send_buffer": 8
  },
  "groups": [
    {
        "mask": "test.log",
        "rules": [
            {
                "regexp": "test string (?P<word_STRING>\\\\w+)",
                "influxdb": {
                    "name": "test"
                }
            },
            {
                "regexp": "(?P<word_STRING>\\\\w+) string one",
                "influxdb": {
                    "name": "test"
                }
            },
            {
                "regexp": "string (?P<word_STRING>\\\\w+) one",
                "influxdb": {
                    "name": "test"
                }
            }
        ]
    }
  ]
}
"""

def bench(logs_count=1, config=config):
    os.system("rm -rf %s; mkdir -p %s" % (log_dir, log_dir))
    with open(config_file, "w") as myfile:
        myfile.write(config)

    for x in range(0, logs_count):
        with open(log_file_name + str(x), "a") as myfile:
            myfile.write(msg*lines_count)

    os.system("/bin/bash -c 'time %s %s'" % (binary, run_params))
    os.system("echo '\n'")
    os.system("rm -f %s %s" % (config_file, log_file_name + '*'))


if __name__ == '__main__':
    print("with 1 file containing %s matching lines each" % (lines_count))
    bench()
    print("with 5 file containing %s matching lines each" % (lines_count))
    bench(5)
    print("with 10 file containing %s matching lines each" % (lines_count))
    bench(10)

    print("with 1 file containing %s matching lines each with 3 rules" % (lines_count))
    bench(1,config2)
    print("with 5 file containing %s matching lines each with 3 rules" % (lines_count))
    bench(5,config2)
    print("with 10 file containing %s matching lines each with 3 rules" % (lines_count))
    bench(10,config2)