package main

import (
	"os/exec"
	"testing"
	"time"
)

// BenchmarkStartupTime measures the startup time of the chatmate binary
func BenchmarkStartupTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("./chatmate", "--version")
		start := time.Now()
		err := cmd.Run()
		elapsed := time.Since(start)

		if err != nil {
			b.Fatalf("Command failed: %v", err)
		}

		// Log the startup time for analysis
		if i == 0 {
			b.Logf("First startup time: %v", elapsed)
		}
	}
}

// BenchmarkStartupTimeOptimized measures the startup time of the optimized binary
func BenchmarkStartupTimeOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("./builds/chatmate-optimized", "--version")
		start := time.Now()
		err := cmd.Run()
		elapsed := time.Since(start)

		if err != nil {
			b.Fatalf("Command failed: %v", err)
		}

		if i == 0 {
			b.Logf("First optimized startup time: %v", elapsed)
		}
	}
}

// BenchmarkStartupTimeMinimal measures the startup time of the minimal binary
func BenchmarkStartupTimeMinimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("./builds/chatmate-minimal", "--version")
		start := time.Now()
		err := cmd.Run()
		elapsed := time.Since(start)

		if err != nil {
			b.Fatalf("Command failed: %v", err)
		}

		if i == 0 {
			b.Logf("First minimal startup time: %v", elapsed)
		}
	}
}

// BenchmarkListCommand benchmarks the list command performance
func BenchmarkListCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("./chatmate", "list", "--available")
		err := cmd.Run()
		if err != nil {
			b.Fatalf("List command failed: %v", err)
		}
	}
}

// BenchmarkStatusCommand benchmarks the status command performance
func BenchmarkStatusCommand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("./chatmate", "status")
		err := cmd.Run()
		if err != nil {
			b.Fatalf("Status command failed: %v", err)
		}
	}
}
