package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/estebgonza/go-richelieu/constants"
	"github.com/urfave/cli/v2"
)

// Execute Entrypoint of the generation plan
func Generate() error {
	// Read input json file plan
	p, err := GetPlanFromFile(constants.DefaultPlanFile)
	if err != nil {
		return err
	}

	// Generate rows
	if err := p.generate(); err != nil {
		return err
	}

	return nil
}

func (p *Plan) generate() error {
	if p == nil {
		return errors.New("Missing plan.json file")
	}
	_ = os.Mkdir("output", os.ModePerm) // Create output directory
	for _, schema := range p.Schemas {  // For each schema
		for _, table := range schema.Tables { // For each table
			err := table.generate(schema.Name) // Generate table dataset
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (table *Table) generate(schemaName string) error {

	// Prepare each column (pre-calculation for performance)
	for c := range table.Columns {
		table.Columns[c].init(table.Mode)
	}

	// Create output folder
	folderName := schemaName + "." + table.Name
	_ = os.Mkdir("output/"+folderName, os.ModePerm)

	// Multithread generation if several files are requested
	wg := sync.WaitGroup{}
	for i := 0; i < table.Files; i++ { // For each file
		wg.Add(1) // Start a dedicated thread

		go func(i int) {
			// Open output file
			fileName := "output/" + folderName + "/export_" + strconv.Itoa(i) + ".csv"
			csvFile, err := os.Create(fileName)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			// Insure last file has exactly appropriate number of rows
			rowsCurrentFile := (table.Rows / table.Files)
			firstRow := (i * rowsCurrentFile)
			if i == table.Files-1 {
				rowsCurrentFile = table.Rows - (table.Files-1)*(table.Rows/table.Files)
			}

			var row []string
			var rowsBuffer []string
			for j := 0; j < rowsCurrentFile; j++ {

				if j == 0 { // Build the first row
					for _, column := range table.Columns {
						row = append(row, column.getValue(firstRow+j, table.Rows))
					}
				} else { // For performane, only update columns with distinct > 1
					for c, column := range table.Columns {
						if column.Distinct > 1 {
							row[c] = column.getValue(firstRow+j, table.Rows)
						}
					}
				}

				rowsBuffer = append(rowsBuffer, strings.Join(row, ","))
				if j%10000 == 0 && j != 0 {
					// Note: for performance, use WriteString rather than a csvWriter
					_, err := csvFile.WriteString(strings.Join(rowsBuffer, "\n") + "\n")
					if err != nil {
						log.Println(err)
						os.Exit(1)
					}
					rowsBuffer = nil

					// Display a progress status
					if table.Rows >= 1000000 && j%100000 == 0 {
						fmt.Printf(".")
					}
				}
			}
			_, err = csvFile.WriteString(strings.Join(rowsBuffer, "\n") + "\n")
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			rowsBuffer = nil
			csvFile.Close()
			wg.Done()
		}(i)
	}
	wg.Wait()
	if table.Rows >= 1000000 {
		fmt.Printf("\n")
	}
	log.Println("Done generating " + schemaName + "." + table.Name)
	return nil
}

// Validate Plan (rows and cardinalities)
func (p *Plan) validate() error {
	// For each column of each table of each schema
	for _, s := range p.Schemas {
		for _, t := range s.Tables {
			if t.Rows < 0 {
				m := fmt.Sprintf("Error. Table %s: Expected rows can't be negative.", t.Name)
				return errors.New(m)
			}
			for _, c := range t.Columns {
				cardinality := c.Distinct
				if cardinality < 1 {
					m := fmt.Sprintf("Error. Column %s.%s.%s: cardinality can't be lower than 1.", s.Name, t.Name, c.Name)
					return errors.New(m)
				}

			}
		}
	}
	return nil
}

func GetPlanFromFile(planfile string) (*Plan, error) {
	var p Plan

	if _, err := os.Stat(planfile); os.IsNotExist(err) {
		return nil, nil // planfile doesnt exist
	}

	planFile, err := os.Open(planfile)
	if err != nil {
		return nil, errors.New("No plan.json found")
	}

	byteValue, _ := ioutil.ReadAll(planFile)
	json.Unmarshal(byteValue, &p)
	planFile.Close()

	// Control input plan
	if err := p.validate(); err != nil {
		return nil, err
	}

	return &p, nil
}

// Read a list of types from arguments (eg: INT, STRING, INT) and initiate a plan.json with it
func CreateFromColumn(args cli.Args) error {

	if args.Len() == 0 {
		return errors.New("Please specify columns type to init a generation plan")
	}

	typeList := strings.ReplaceAll(args.Get(0), " ", "")
	cols := strings.Split(typeList, ",")
	var columns []Column
	for index, t := range cols {
		if !ChecksSupportedType(t) {
			return errors.New("Unsupported column type " + t)
		}
		var pc Column
		pc.Name = strings.ToLower(strconv.Itoa(index) + "_" + t)
		pc.Distinct = 1
		pc.Type = t
		columns = append(columns, pc)
	}

	var table = Table{Name: "table1", Rows: 10000, Files: 1, Columns: columns}
	var schema = Schema{Name: "schema1", Tables: []Table{table}}
	var plan = Plan{Schemas: []Schema{schema}}

	return WriteToFile(plan, constants.DefaultPlanFile)
}

func WriteToFile(plan Plan, planfile string) error {

	// If planfile exist, we merge new plan with existing plan
	existingPlan, err := GetPlanFromFile(planfile)
	if err != nil {
		return err
	}
	if existingPlan != nil {
		MergePlanParts(&plan, existingPlan)
	}

	json, err := json.MarshalIndent(plan, "", "    ")
	if err != nil {
		return err
	}

	// Cosmetic corrections to have columns on one line (usefull for large schema)
	json = []byte(strings.ReplaceAll(string(json), "\n                            \"type\"", "\"type\""))
	json = []byte(strings.ReplaceAll(string(json), "\n                            \"name\"", "\"name\""))
	json = []byte(strings.ReplaceAll(string(json), "\n                            \"distinct\"", "\"distinct\""))
	json = []byte(strings.ReplaceAll(string(json), "\n                            \"start\"", "\"start\""))
	json = []byte(strings.ReplaceAll(string(json), "\n                            \"end\"", "\"end\""))
	json = []byte(strings.ReplaceAll(string(json), "\n                            \"values\"", "\"values\""))
	json = []byte(strings.ReplaceAll(string(json), "\n                        }", "}"))

	ioutil.WriteFile(planfile, json, 0644)

	return nil
}
