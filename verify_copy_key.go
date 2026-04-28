// +build ignore

// Simple verification script to check that Copy key is properly defined
package main

import (
	"fmt"
	"github.com/leg100/pug/internal/tui/keys"
)

func main() {
	// Verify Copy key is defined
	copy := keys.Global.Copy
	help := copy.Help()

	fmt.Println("Copy Key Binding Verification")
	fmt.Println("===============================")
	fmt.Printf("Key(s): %v\n", copy.Keys())
	fmt.Printf("Help Key: %s\n", help.Key)
	fmt.Printf("Help Desc: %s\n", help.Desc)
	fmt.Println("\nThe Copy command will appear in the help menu when you press '?'")

	// Verify it's in the KeyMapToSlice output
	allKeys := keys.KeyMapToSlice(keys.Global)
	found := false
	for _, k := range allKeys {
		if k.Keys()[0] == "C" {
			found = true
			break
		}
	}

	if found {
		fmt.Println("✓ Copy key found in Global keys slice")
	} else {
		fmt.Println("✗ Copy key NOT found in Global keys slice")
	}
}

