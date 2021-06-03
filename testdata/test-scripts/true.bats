#!/usr/bin/env bats

@test "true" {
  run true
  [ $status = "0" ]
}
