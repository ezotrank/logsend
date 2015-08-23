package com.github.logsend

import org.http4s._
import org.http4s.server._
import org.http4s.dsl._

import _root_.argonaut._, Argonaut._
import org.http4s.argonaut._


class Service (
	val _name: String,
	val _filesMask: List[String],
	val _chunkSize: Int,
	val _compress: Boolean
){

	def name: String = _name

}

class Server (
	val nameMask: String,
	val services: List[Service]
){}

object Config {

}

object Web {

  def getConfig(req: Request): String = {
		println(req.params("hostname"))
		return """
		{
			"PushAddr": "pxl2.int.avs.io:2930",
			"FilesMask": ["/home/aviasales/fuzzy/shared/logs/bee*.log"],
			"ChunkSize": 100,
			"Compress": true	
		}
		"""
	}

  val service = HttpService {
		case req@GET -> Root / "config" => Ok(getConfig(req))
  }
}
