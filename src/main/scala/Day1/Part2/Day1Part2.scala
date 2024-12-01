package Day1.Part2

import scala.io.Source
import scala.util.Using

@main def Day1Part1 = {
  val (leftList, rightList) =
    Using(Source.fromFile("src/main/scala/Day1/input.txt")) { src =>
      val lists = src
        .getLines()
        .map(_.split(" {3}").map(_.toInt).toArray)
        .toArray
        .transpose
      (lists(0), lists(1))
    }.get

  println(
    leftList.map { entry =>
      entry * rightList.count(_ == entry)
    }.sum
  )
}
