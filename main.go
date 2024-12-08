package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rivo/tview"
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
	// create a new TUI app
	app := tview.NewApplication()
	loadInventory()
	inventoryList := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)

	inventoryList.SetBorder(true).SetTitle("Inventory Items")

	refreshInventory := func() {
		inventoryList.Clear()
		if len(inventory) == 0 {
			fmt.Fprint(inventoryList, "No item in inventory")
		} else {
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (stock: %d)\n", i+1, item.Name, item.Stock)
			}
		}
	}
	itemNameInput := tview.NewInputField().SetLabel("Item name: ")
	itemStockInput := tview.NewInputField().SetLabel("Stock: ")
	itemIDInput := tview.NewInputField().SetLabel("Item ID to delete: ")

	form := tview.NewForm().
		AddFormItem(itemNameInput).AddFormItem(itemStockInput).AddFormItem(itemIDInput).
		AddButton("Add Item", func() {
			name := itemNameInput.GetText()
			stock := itemStockInput.GetText()
			if name != "" && stock != "" {
				quantity, err := strconv.Atoi(stock)
				if err != nil {
					fmt.Fprintln(inventoryList, "Invalid Stock value")
				}
				inventory = append(inventory, Item{Name: name, Stock: quantity})
				saveInventory()
				refreshInventory()
				itemNameInput.SetText("")
				itemStockInput.SetText("")
			}
		}).
		AddButton("Delete Item ", func() {
			idStr := itemIDInput.GetText()
			if idStr == "" {
				fmt.Fprintln(inventoryList, "Please enter an item ID to delete.")
				return
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Fprintln(inventoryList, "Invalid item ID.")
			}
			deleteItem(id - 1)
			fmt.Fprintln(inventoryList, "Item [%d] delete. \n", id)
			refreshInventory()
			itemIDInput.SetText("")
		}).
		AddButton("Exit ", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Manage Inventory").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().AddItem(inventoryList, 0, 1, false).AddItem(form, 0, 1, true)

	refreshInventory()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
