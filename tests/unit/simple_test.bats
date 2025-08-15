#!/usr/bin/env bats

# Simple test to verify testing framework works

@test "addition test" {
    result="$((2 + 2))"
    [ "$result" -eq 4 ]
}

@test "chatmate directory exists" {
    [ -d "mates" ]
}

@test "hire script exists" {
    [ -f "hire.sh" ]
}
