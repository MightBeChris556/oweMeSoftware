package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/dixonwille/wmenu/v5"
	"github.com/google/uuid"
	c "github.com/ostafen/clover"
	"log"
	"os"
)

type Debt struct {
	Amount   int64  `json:"amount"`
	Name     string `json:"name"`
	id       string
	DebtName string `json:"debtname"`
}

func main() {
	databaseSetup()

	menu := wmenu.NewMenu("Make a selection")
	menu.Action(func(opts []wmenu.Opt) error {
		fmt.Println(opts[0].Value)
		if opts[0].Value == "Add" {
			addDebt()

		} else if opts[0].Value == "Show" {
			showDatabase()

		} else if opts[0].Value == "Increase" {
			increaseDebt()

		} else if opts[0].Value == "Increase" {
			decreaseDebt()

		} else if opts[0].Value == "Delete" {
			deleteDebt()

		} else if opts[0].Value == "Export" {
			exportToCSV("data.csv")
		}
		menuError := menu.Run()
		if menuError != nil {
			fmt.Println("menu error")
			log.Fatal(menuError)
		}

		return nil
	})
	menu.Option("Add Debt", "Add", false, nil)
	menu.Option("Delete Debt", "Delete", false, nil)
	menu.Option("Increase Debt", "Increase", false, nil)
	menu.Option("Decrease Debt", "Decrease", false, nil)

	menu.Option("Show Debts", "Show", false, nil)
	menu.Option("Export Debts", "Export", false, nil)

	menuError := menu.Run()
	if menuError != nil {
		fmt.Println("menu error")
		log.Fatal(menuError)
	}

}

func databaseSetup() {
	db, _ := c.Open("clover-db")
	defer db.Close()
	collectionExists, _ := db.HasCollection("debts")

	if !collectionExists {
		// Create a collection named 'todos'
		db.CreateCollection("debts")
	} else {
		db.Close()
	}

}
func showDatabase() {
	db, _ := c.Open("clover-db")

	debts, _ := db.Query("debts").FindAll()

	t := tabby.New()
	t.AddHeader("NAME", "DEBT", "Amount")

	for _, debt := range debts {
		t.AddLine(debt.Get("name"), debt.Get("debtName"), debt.Get("amount"))
	}

	t.Print()
	db.Close()
}

func addDebt() {
	db, _ := c.Open("clover-db")

	var name string
	var amount int
	var debtName string

	var id = uuid.New().String()

	fmt.Println("Enter the debtors name: ")
	fmt.Scanln(&name)
	fmt.Println("Enter the debt amount: ")
	fmt.Scanln(&amount)
	fmt.Println("What is the debt for?: ")
	fmt.Scanln(&debtName)
	//var _ = Debt{name: name, id: id, amount: amount, debtName: debtName}
	doc := c.NewDocument()
	doc.Set("name", name)
	doc.Set("amount", amount)
	doc.Set("debtName", debtName)
	doc.Set("id", id)

	docId, _ := db.InsertOne("debts", doc)
	fmt.Println(docId)
	db.Close()

}

func deleteDebt() {
	db, _ := c.Open("clover-db")
	defer db.Close()

	var name string
	var debtName string

	fmt.Println("Enter the debtors name: ")
	fmt.Scanln(&name)

	fmt.Println("What is the debt for?: ")
	fmt.Scanln(&debtName)

	db.Query("debts").Where(c.Field("name").Eq(name).And(c.Field("debtName").Eq(debtName))).Delete()

}

func increaseDebt() {
	db, _ := c.Open("clover-db")

	var name string
	var debtName string
	var debtIncreaseAmount int64

	fmt.Println("Enter the debtors name: ")
	fmt.Scanln(&name)

	fmt.Println("What is the debt for?: ")
	fmt.Scanln(&debtName)

	fmt.Println("increase by how much?: ")
	fmt.Scanln(&debtIncreaseAmount)

	var debt, _ = db.Query("debts").Where(c.Field("name").Eq(name).And(c.Field("debtName").Eq(debtName))).FindFirst()
	var debtID = debt.Get("_id").(string)
	var currentDebt = debt.Get("amount").(int64)
	var newDebt = currentDebt + debtIncreaseAmount

	db.Query("debts").UpdateById(debtID, map[string]interface{}{"amount": newDebt})

	db.Close()

}

func decreaseDebt() {
	db, _ := c.Open("clover-db")

	var name string
	var debtName string
	var debtDecreaseAmount int64

	fmt.Println("Enter the debtors name: ")
	fmt.Scanln(&name)

	fmt.Println("What is the debt for?: ")
	fmt.Scanln(&debtName)

	fmt.Println("increase by how much?: ")
	fmt.Scanln(&debtDecreaseAmount)

	var debt, _ = db.Query("debts").Where(c.Field("name").Eq(name).And(c.Field("debtName").Eq(debtName))).FindFirst()
	var debtID = debt.Get("_id").(string)
	var currentDebt = debt.Get("amount").(int64)
	var newDebt = currentDebt - debtDecreaseAmount

	db.Query("debts").UpdateById(debtID, map[string]interface{}{"amount": newDebt})

	db.Close()

}

func exportToCSV(destination string) error {
	db, _ := c.Open("clover-db")

	db.ExportCollection("debts", "debts.json")
	db.Close()
	sourceFile, err := os.Open("debts.json")
	if err != nil {
		return err
	}
	// close file at end
	defer sourceFile.Close()

	var debtList []Debt
	if err := json.NewDecoder(sourceFile).Decode(&debtList); err != nil {
		return err
	}

	// 3. Create a new file
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 4. Write the header of the CSV file
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"name", "debtName", "amount"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range debtList {
		var csvRow []string
		csvRow = append(csvRow, r.Name, r.DebtName, fmt.Sprint(r.Amount))
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
