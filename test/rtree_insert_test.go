package test

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"rtree/src"
)

// generateUUID generates a random UUID v4
func generateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		log.Fatal("Failed to generate UUID:", err)
	}

	// Set version (4) and variant bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16])
}

// TestRTreeWithUUIDs performs comprehensive testing of RTree with nodeToCreate UUIDs
func TestRTreeWithUUIDs(t *testing.T) {
	nodeToCreate := 100000

	fmt.Println("🚀 Starting RTree UUID Performance Test")
	fmt.Println(strings.Repeat("=", 50))

	// Create new tree
	tree := src.NewRTree()
	uuids := make([]string, nodeToCreate)

	// Generate nodeToCreate UUIDs
	fmt.Println("📝 Generating nodeToCreate UUIDs...")
	for i := 0; i < nodeToCreate; i++ {
		uuids[i] = generateUUID()
	}
	fmt.Printf("✅ Generated %d UUIDs\n\n", len(uuids))

	// Test 1: Insert Performance
	fmt.Println("🔧 Test 1: INSERT Performance")
	fmt.Println(strings.Repeat("-", 30))

	startTime := time.Now()
	successCount := 0
	failureCount := 0

	for i, uuid := range uuids {
		if src.Add(uuid, uuid, tree) {
			successCount++
			// fmt.Printf("UUID %d: %s\n", i, uuid)
		} else {
			failureCount++
			fmt.Printf("❌ Failed to insert UUID %d: %s\n", i, uuid)
		}

		// Progress indicator every 100 insertions
		if (i+1)%int(nodeToCreate/10) == 0 {
			fmt.Printf("📊 Progress: %d/%d inserted\n", i+1, len(uuids))
		}
	}
	// src.PrintNode(tree.Root, true)
	// return

	insertDuration := time.Since(startTime)
	fmt.Printf("\n📈 INSERT RESULTS:\n")
	fmt.Printf("   ✅ Successful: %d\n", successCount)
	fmt.Printf("   ❌ Failed: %d\n", failureCount)
	fmt.Printf("   ⏱️  Duration: %v\n", insertDuration)
	fmt.Printf("   📊 Rate: %.2f insertions/second\n", float64(successCount)/insertDuration.Seconds())
	fmt.Printf("   📊 Avg per insert: %v\n\n", insertDuration/time.Duration(successCount))

	// Test 2: Search Performance
	fmt.Println("🔍 Test 2: SEARCH Performance")
	fmt.Println(strings.Repeat("-", 30))

	startTime = time.Now()
	foundCount := 0
	notFoundCount := 0

	for i, uuid := range uuids {
		node := src.Search(uuid, tree)
		if node != nil && node.IsEnd && node.Value == uuid {
			foundCount++
		} else {
			notFoundCount++
			fmt.Printf("❌ Search failed for UUID %d: %s\n", i, uuid)
		}

		// Progress indicator every 100 searches
		if (i+1)%int(nodeToCreate/10) == 0 {
			fmt.Printf("📊 Progress: %d/%d searched\n", i+1, len(uuids))
		}
	}

	searchDuration := time.Since(startTime)
	fmt.Printf("\n📈 SEARCH RESULTS:\n")
	fmt.Printf("   ✅ Found: %d\n", foundCount)
	fmt.Printf("   ❌ Not Found: %d\n", notFoundCount)
	fmt.Printf("   ⏱️  Duration: %v\n", searchDuration)
	fmt.Printf("   📊 Rate: %.2f searches/second\n", float64(foundCount)/searchDuration.Seconds())
	fmt.Printf("   📊 Avg per search: %v\n\n", searchDuration/time.Duration(len(uuids)))

	// Test 3: Tree Structure Analysis
	fmt.Println("🌳 Test 3: TREE STRUCTURE Analysis")
	fmt.Println(strings.Repeat("-", 35))

	treeStats := analyzeTreeStructure(tree.Root, 0)
	fmt.Printf("📊 Tree Statistics:\n")
	fmt.Printf("   🔢 Total nodes: %d\n", treeStats.TotalNodes)
	fmt.Printf("   🍃 Leaf nodes: %d\n", treeStats.LeafNodes)
	fmt.Printf("   🌿 Internal nodes: %d\n", treeStats.InternalNodes)
	fmt.Printf("   📏 Max depth: %d\n", treeStats.MaxDepth)
	fmt.Printf("   📊 Avg depth: %.2f\n", treeStats.AvgDepth)
	fmt.Printf("   🎯 Terminal nodes: %d\n", treeStats.TerminalNodes)

	// Test 4: Compact Operation
	fmt.Println("\n🗜️  Test 4: COMPACT Operation")
	fmt.Println(strings.Repeat("-", 30))

	startTime = time.Now()
	src.Compact(tree)
	compactDuration := time.Since(startTime)

	treeStatsAfterCompact := analyzeTreeStructure(tree.Root, 0)
	fmt.Printf("📊 Compact Results:\n")
	fmt.Printf("   ⏱️  Duration: %v\n", compactDuration)
	fmt.Printf("   🔢 Nodes before: %d\n", treeStats.TotalNodes)
	fmt.Printf("   🔢 Nodes after: %d\n", treeStatsAfterCompact.TotalNodes)
	fmt.Printf("   📉 Nodes reduced: %d\n", treeStats.TotalNodes-treeStatsAfterCompact.TotalNodes)

	// Test 5: Random Search Test (search some random UUIDs that don't exist)
	fmt.Println("\n🎲 Test 5: RANDOM SEARCH Test (Non-existent UUIDs)")
	fmt.Println(strings.Repeat("-", 50))

	randomUUIDs := make([]string, 100)
	for i := 0; i < 100; i++ {
		randomUUIDs[i] = generateUUID()
	}

	startTime = time.Now()
	randomFoundCount := 0
	for _, uuid := range randomUUIDs {
		node := src.Search(uuid, tree)
		if node != nil && node.IsEnd {
			randomFoundCount++
		}
	}
	randomSearchDuration := time.Since(startTime)

	fmt.Printf("📊 Random Search Results:\n")
	fmt.Printf("   🔍 Searched: 100 random UUIDs\n")
	fmt.Printf("   ✅ Found (unexpected): %d\n", randomFoundCount)
	fmt.Printf("   ❌ Not found (expected): %d\n", 100-randomFoundCount)
	fmt.Printf("   ⏱️  Duration: %v\n", randomSearchDuration)

	// Test 6: Delete Performance (delete 10% of UUIDs)
	fmt.Println("\n🗑️  Test 6: DELETE Performance (10% of UUIDs)")
	fmt.Println(strings.Repeat("-", 40))

	deleteCount := len(uuids) / 10
	startTime = time.Now()
	deleteSuccessCount := 0

	for i := 0; i < deleteCount; i++ {
		if src.Delete(uuids[i], tree) {
			deleteSuccessCount++
		}
	}

	deleteDuration := time.Since(startTime)
	fmt.Printf("📊 Delete Results:\n")
	fmt.Printf("   🗑️  Attempted: %d\n", deleteCount)
	fmt.Printf("   ✅ Successful: %d\n", deleteSuccessCount)
	fmt.Printf("   ❌ Failed: %d\n", deleteCount-deleteSuccessCount)
	fmt.Printf("   ⏱️  Duration: %v\n", deleteDuration)

	// Verify deletions
	verifyStartTime := time.Now()
	deletedVerifyCount := 0
	for i := 0; i < deleteCount; i++ {
		node := src.Search(uuids[i], tree)
		if node == nil || !node.IsEnd {
			deletedVerifyCount++
		}
	}
	verifyDuration := time.Since(verifyStartTime)

	fmt.Printf("   ✅ Verified deleted: %d/%d (in %v)\n", deletedVerifyCount, deleteCount, verifyDuration)

	// Final Summary
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 TEST SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("📊 Total UUIDs processed: %d\n", len(uuids))
	fmt.Printf("⏱️  Total test duration: %v\n", time.Since(time.Now().Add(-insertDuration-searchDuration-compactDuration-randomSearchDuration-deleteDuration-verifyDuration)))
	fmt.Printf("✅ Insert success rate: %.2f%%\n", float64(successCount)/float64(len(uuids))*100)
	fmt.Printf("✅ Search success rate: %.2f%%\n", float64(foundCount)/float64(len(uuids))*100)
	fmt.Printf("✅ Delete success rate: %.2f%%\n", float64(deleteSuccessCount)/float64(deleteCount)*100)
	fmt.Printf("🌳 Final tree nodes: %d\n", treeStatsAfterCompact.TotalNodes)
	fmt.Printf("🏆 Overall Status: %s\n", getOverallStatus(successCount, foundCount, len(uuids)))
}

// TreeStats holds statistics about the tree structure
type TreeStats struct {
	TotalNodes    int
	LeafNodes     int
	InternalNodes int
	TerminalNodes int
	MaxDepth      int
	AvgDepth      float64
	TotalDepth    int
	DepthCount    int
}

// analyzeTreeStructure recursively analyzes the tree structure
func analyzeTreeStructure(node *src.Node, depth int) TreeStats {
	if node == nil {
		return TreeStats{}
	}

	stats := TreeStats{
		TotalNodes: 1,
		MaxDepth:   depth,
		TotalDepth: depth,
		DepthCount: 1,
	}

	if node.IsEnd {
		stats.TerminalNodes = 1
	}

	if len(node.Children) == 0 {
		stats.LeafNodes = 1
	} else {
		stats.InternalNodes = 1
		for _, child := range node.Children {
			childStats := analyzeTreeStructure(child, depth+1)
			stats.TotalNodes += childStats.TotalNodes
			stats.LeafNodes += childStats.LeafNodes
			stats.InternalNodes += childStats.InternalNodes
			stats.TerminalNodes += childStats.TerminalNodes
			stats.TotalDepth += childStats.TotalDepth
			stats.DepthCount += childStats.DepthCount
			if childStats.MaxDepth > stats.MaxDepth {
				stats.MaxDepth = childStats.MaxDepth
			}
		}
	}

	if stats.DepthCount > 0 {
		stats.AvgDepth = float64(stats.TotalDepth) / float64(stats.DepthCount)
	}

	return stats
}

// getOverallStatus returns overall test status
func getOverallStatus(successCount, foundCount, totalCount int) string {
	insertRate := float64(successCount) / float64(totalCount)
	searchRate := float64(foundCount) / float64(totalCount)

	if insertRate >= 1.0 && searchRate >= 1.0 {
		return "🎉 EXCELLENT - All operations successful!"
	} else if insertRate >= 0.95 && searchRate >= 0.95 {
		return "✅ GOOD - Most operations successful"
	} else if insertRate >= 0.90 && searchRate >= 0.90 {
		return "⚠️  ACCEPTABLE - Some issues detected"
	} else {
		return "❌ POOR - Significant issues found"
	}
}
