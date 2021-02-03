package generator

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/estebgonza/go-richelieu/constants"
)

// Execute Entrypoint of the generation plan
func Generate() error {
	// Read input json file plan
	p, err := ReadFromFile(constants.DefaultPlanFile)
	if err != nil {
		return err
	}

	// Initialize Column from PlanColumns
	if err := initializeColumns(p); err != nil {
		return err
	}

	// Generate rows
	if err := generate(p); err != nil {
		return err
	}

	// Export load commands
	if err := exportLoadDbCommands(p); err != nil {
		return err
	}

	return nil
}

func generate(p *Plan) error {
	_ = os.Mkdir("output", os.ModePerm) // Create output directory

	for _, schema := range p.Schemas { // For each schema
		for _, table := range schema.Tables { // For each table

			folderName := schema.Name + "." + table.Name
			_ = os.Mkdir("output/"+folderName, os.ModePerm)

			// Multithread generation if several files are requested
			wg := sync.WaitGroup{}
			for i := 0; i < table.Files; i++ { // For each file
				wg.Add(1) // Start a dedicated thread

				go func(i int) {
					fileName := "output/" + folderName + "/export_" + strconv.Itoa(i) + ".csv"
					csvFile, err := os.Create(fileName)
					if err != nil {
						log.Println(err)
						os.Exit(1)
					}
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
			log.Println("Done generating " + schema.Name + "." + table.Name)
		}
	}
	return nil
}

func initializeColumns(p *Plan) error {
	// For each column of each table of each schema:
	// - default prefix for string fields
	// - mode default to table mode
	// - split valuesList to slices
	// - pre-calculate float and date steps
	for s := range p.Schemas {
		for t := range p.Schemas[s].Tables {
			for c := range p.Schemas[s].Tables[t].Columns {
				cp := &p.Schemas[s].Tables[t].Columns[c]
				if cp.Type == "STRING" && cp.Prefix == "" {
					cp.Prefix = "txt_"
				}
				if cp.Mode == "" {
					cp.Mode = p.Schemas[s].Tables[t].Mode
				}
				if cp.Mode == "" {
					cp.Mode = "ALTERNATE"
				}
				if cp.ValuesList != "" {
					cp.ValuesSlice = strings.Split(cp.ValuesList, ";")
				}
				if cp.Type == "FLOAT" {
					v1, _ := strconv.ParseFloat(cp.Start, 32)
					v2, _ := strconv.ParseFloat(cp.End, 32)
					if v2 <= v1 {
						v2 = v1 + 1.0
					}
					cp.FloatStart = v1
					cp.FloatStep = (v2 - v1) / float64(cp.Distinct)
				}
				if cp.Type == "DATE" {
					if cp.Start == "" {
						cp.Start = "2020-01-01 00:00:00"
					}
					if cp.End == "" {
						cp.End = "2020-12-31 00:00:00"
					}
					v1, _ := time.Parse("2006-01-02 15:04:05", cp.Start)
					v2, _ := time.Parse("2006-01-02 15:04:05", cp.End)
					cp.DateStart = v1.Unix()
					cp.DateStep = (v2.Unix() - v1.Unix()) / int64(cp.Distinct)
				}
			}
		}
	}
	return nil
}

// Validate Plan (rows and cardinalities)
func validate(p *Plan) error {
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

func exportLoadDbCommands(p *Plan) error {
	var commands []string
	var cmd string

	cmd = "# Load data to S3 after adapting S3 repository"
	commands = append(commands, cmd)
	cmd = "aws s3 cp ./output " + constants.DefaultS3Repository + " --recursive"
	commands = append(commands, cmd)
	commands = append(commands, "\n")

	for _, s := range p.Schemas {
		for _, t := range s.Tables {
			fullName := s.Name + "." + t.Name
			cmd = "LOAD DATA INPATH '" + strings.ReplaceAll(constants.DefaultS3Repository, "s3://", "s3a://") + "/" + fullName + "' INTO TABLE " + fullName + " FORMAT CSV SEPARATOR ',';"
			commands = append(commands, cmd)
			cmd = "COMMIT " + fullName + ";"
			commands = append(commands, cmd)
		}
	}

	cmdFile, err := os.Create("./output/loadCommands.txt")
	if err != nil {
		return err
	}

	for _, l := range commands {
		cmdFile.WriteString(l + "\n")
	}
	cmdFile.Close()

	return nil
}
