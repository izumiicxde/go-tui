package main

import (
	"encoding/json"
	"log"
	"os"
)

// Structure of the inventory Item.
type Item struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

// Global variable
var (
	inventory     = []Item{}
	inventoryFile = "inventory.json"
)

// 4 required functions
// 1- Load inventory
func loadInventory() {
	if _, err := os.Stat(inventoryFile); err == nil {
		data, err := os.ReadFile(inventoryFile)
		if err != nil {
			log.Fatal("error reading inventory file: ", err)
		}
		json.Unmarshal(data, &inventory)
	}
}

func saveInventory() {
	data, err := json.MarshalIndent(inventory, "", " ")
	if err != nil {
		log.Fatal("error saving inventory: ", err)
	}
	os.WriteFile(inventoryFile, data, 0644) // only we can write to the file
}

func deleteItem(index int) {
	if index < 0 || index >= len(inventory) {
		log.Fatal("Invalid item index")
		return
	}
	// this is to delete an item from the slice
	// first is till the index and next is after the index ignoring the index
	inventory = append(inventory[:index], inventory[index+1:]...)
	saveInventory()
}

func main() {
}
