package com.github.logsend

import org.http4s.server.blaze.BlazeBuilder
import akka.actor.ActorSystem
import akka.actor.Props

object Main extends App {
	println("starting socket server")
	val socketPort = 2930
  val system = ActorSystem("Main")
  val ac = system.actorOf(Props(new SocketServer(socketPort)))

print("starting web server")
  BlazeBuilder.bindHttp(8000, "0.0.0.0")
    .mountService(Web.service, "/")
    .run
    .awaitShutdown()
}
