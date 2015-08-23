package com.github.logsend

import akka.actor.{ Actor, ActorRef, Props }
import akka.io.{ IO, Tcp }
import akka.util.ByteString
import java.net.InetSocketAddress

class SocketServer(port: Int) extends Actor {
 
  import Tcp._
  import context.system
 
  IO(Tcp) ! Bind(self, new InetSocketAddress("0.0.0.0", port))
 
  def receive = {
    case b @ Bound(localAddress) =>
      // do some logging or setup ...
 
    case CommandFailed(_: Bind) => context stop self
 
    case c @ Connected(remote, local) =>
      val handler = context.actorOf(Props[SimplisticHandler])
      val connection = sender()
      connection ! Register(handler)
  }
 
}

class SimplisticHandler extends Actor {
  import Tcp._

	def process(msg: ByteString) = {
		println(msg.decodeString("utf-8"))
	}

  def receive = {
    case Received(data) => process(data)
    case PeerClosed     => context stop self
  }
}
