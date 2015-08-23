lazy val commonSettings = Seq(
  organization := "com.github",
  version := "0.1.0",
  scalaVersion := "2.11.7"
)

lazy val root = (project in file(".")).
  settings(commonSettings: _*).
  settings(
    name := "LogsendServer",
    resolvers += "Scalaz Bintray Repo" at "http://dl.bintray.com/scalaz/releases",
    resolvers += "Akka Snapshot Repository" at "http://repo.akka.io/snapshots/",
    libraryDependencies ++= Seq(
	  "org.http4s" %% "http4s-blazeserver" % "0.8.4",
	  "org.http4s" %% "http4s-dsl"         % "0.8.4",
	  "org.http4s" %% "http4s-argonaut"    % "0.8.4",
	  "com.typesafe.akka" %% "akka-actor" % "2.4-SNAPSHOT",
	  "com.typesafe.akka" %% "akka-stream-experimental" % "1.0"
   )
  )
