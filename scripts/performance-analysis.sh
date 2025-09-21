#!/bin/bash

# Performance profiling script for chatmate CLI
set -e

echo "🔍 Running performance analysis for chatmate CLI..."

# Create profiles directory
mkdir -p profiles

# Build optimized binary if it doesn't exist
if [ ! -f "builds/chatmate-optimized" ]; then
    echo "Building optimized binary..."
    ./scripts/build-optimized.sh > /dev/null 2>&1
fi

echo ""
echo "📊 Running benchmarks..."

# Run benchmarks
go test -bench=. -benchmem -cpuprofile=profiles/cpu.prof -memprofile=profiles/mem.prof -run=^$ .

echo ""
echo "🚀 Analyzing startup performance..."

# Test startup times
echo "Measuring startup times (10 runs each):"

echo -n "Standard build: "
{ time for i in {1..10}; do ./builds/chatmate-standard --version >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{print $2}'

echo -n "Optimized build: "
{ time for i in {1..10}; do ./builds/chatmate-optimized --version >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{print $2}'

echo -n "Minimal build: "
{ time for i in {1..10}; do ./builds/chatmate-minimal --version >/dev/null 2>&1; done; } 2>&1 | grep real | awk '{print $2}'

echo ""
echo "💾 Memory usage analysis..."

# Memory usage for different builds
echo "Memory usage (RSS in KB):"
echo -n "Standard: "
/usr/bin/time -l ./builds/chatmate-standard --version 2>&1 | grep "maximum resident set size" | awk '{print $1}' || echo "N/A"

echo -n "Optimized: "
/usr/bin/time -l ./builds/chatmate-optimized --version 2>&1 | grep "maximum resident set size" | awk '{print $1}' || echo "N/A"

echo -n "Minimal: "
/usr/bin/time -l ./builds/chatmate-minimal --version 2>&1 | grep "maximum resident set size" | awk '{print $1}' || echo "N/A"

echo ""
echo "📈 Binary analysis..."

echo "Binary sizes:"
ls -lh builds/chatmate-* | awk '{print $9 ": " $5}'

echo ""
echo "File type analysis:"
file builds/chatmate-optimized

echo ""
echo "📁 Profile files generated:"
ls -lh profiles/

echo ""
echo "🔧 Analysis suggestions:"
echo "1. Run 'go tool pprof profiles/cpu.prof' for CPU profiling"
echo "2. Run 'go tool pprof profiles/mem.prof' for memory profiling" 
echo "3. Use 'go tool pprof -http=:8080 profiles/cpu.prof' for web interface"

echo ""
echo "✅ Performance analysis complete!"
