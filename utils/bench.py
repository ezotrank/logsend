import os
import time

line_per_second = 25000
cycles = 100
log_file_name = 'test.log'
msg = "test string one\n"

try:
    os.remove(log_file_name)
except:
    pass

with open(log_file_name, "a") as myfile:
    for x in range(0, cycles):
        myfile.write(msg*line_per_second)
        time.sleep(1)
        print("%s/%s" % (x, cycles))