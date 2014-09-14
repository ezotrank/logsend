#
# start python tornado_srv.py --log_file_max_size=1000 --log_file_num_backups=100 --logging=debug --log_file_prefix=bee.10.log
import tornado.ioloop
import tornado.web

from tornado.options import define, options

class MainHandler(tornado.web.RequestHandler):
    def get(self):
        self.write("Hello, world")

application = tornado.web.Application([
    (r"/", MainHandler),
])

if __name__ == "__main__":
    tornado.options.parse_command_line()

    application.listen(8888)
    tornado.ioloop.IOLoop.instance().start()
