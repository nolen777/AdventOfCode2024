package Day1.Part1

import scala.io.Source
import scala.util.Using

@main def Day1Part1 = {
  val (leftList, rightList) =
    Using(Source.fromFile("resources/Day1/input.txt")) { src =>
      val lists = src
        .getLines()
        .map(_.split(" {3}").map(_.toInt).toArray)
        .toArray
        .transpose
      (lists(0), lists(1))
    }.get

  println(
    leftList.sorted
      .zip(rightList.sorted)
      .map { case (left, right) => Math.abs(left - right) }
      .sum
  )
}
