package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func getExcelCol(x int) string {
	chars := [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	if x%26 == 0 { // Z
		modulus := chars[25]
		remainder := (x / 26) - 1
		if remainder > 0 {
			return chars[remainder-1] + modulus
		}
		return modulus
	} else {
		modulus := chars[x%26-1]
		remainder := (x / 26)
		if remainder > 0 {
			return chars[remainder-1] + modulus
		}
		return modulus
	}

}
func main() {

	properties := map[string]map[string]map[string]map[string]string{}
	// cartridges: from the less important to the most
	cartridges := [...]string{}
	// define locales: you must include all the locales
	locales := [...]string{"br", "es", "it", "ja_JP", "kr", "pt", "ru", "zh", "it_IT", "de", "fr", "fr_FR", "it_IT", "zh_CN", "en_GB", "jp", "en", "ja", "en_SE"}
	for _, cartridge := range cartridges {
		var processedFiles []string
		path := ""
		dir, err := os.Open(path)
		if err != nil {
			fmt.Println("Fail to open folder: " + path)
			continue
		}
		defer dir.Close()
		fileInfos, err := dir.Readdir(-1)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, exists := properties[cartridge]
		if !exists {
			properties[cartridge] = map[string]map[string]map[string]string{}
		}
		for _, locale := range locales {
			for _, fi := range fileInfos {
				if strings.HasSuffix(fi.Name(), "_"+locale+".properties") {
					defaultFile := strings.TrimSuffix(fi.Name(), "_"+locale+".properties")
					_, exists := properties[cartridge][defaultFile]
					if !exists {
						properties[cartridge][defaultFile] = map[string]map[string]string{}
						//fmt.Println(fi.Name())
						// process defaultFile
						filename := filepath.Join(path, defaultFile+".properties")
						file, _ := os.Open(filename)
						if err != nil {
							log.Fatal(err)
						}
						if file != nil {
							properties[cartridge][defaultFile]["default"] = map[string]string{}
							defer file.Close()
							//fmt.Println(file.Name())
							scanner := bufio.NewScanner(file)
							for scanner.Scan() {
								if !strings.HasPrefix(scanner.Text(), "#") && len(scanner.Text()) > 0 {
									property := strings.SplitN(scanner.Text(), "=", 2)
									if len(property) < 2 {
										fmt.Println("Error parsing for property: " + scanner.Text() + " // file:" + filename)
										continue
									}
									properties[cartridge][defaultFile]["default"][property[0]] = property[1]
								}
							}
							processedFiles = append(processedFiles, defaultFile+".properties")
						}
					}
					_, exists = properties[cartridge][defaultFile][locale]
					if !exists {
						properties[cartridge][defaultFile][locale] = map[string]string{}
					}
					// start processing file
					filename := filepath.Join(path, fi.Name())
					//fmt.Println(fi.Name())
					file, _ := os.Open(filename)
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						if !(strings.HasPrefix(scanner.Text(), "#")) && len(scanner.Text()) > 0 {
							text := scanner.Text()
							property := strings.SplitN(text, "=", 2)
							if len(property) < 2 {
								fmt.Println("Error parsing for property: " + scanner.Text() + " // file:" + filename)
								continue
							}
							properties[cartridge][defaultFile][locale][property[0]] = property[1]
						}
					}
					processedFiles = append(processedFiles, fi.Name())
				}
			}
		}
		// process file not in locales (files that have only default)
		for _, fi := range fileInfos {
			if strings.HasSuffix(fi.Name(), ".properties") {
				//fmt.Println(fi.Name())
				alreadyProcessed := false
				for _, b := range processedFiles {
					if b == fi.Name() {
						//fmt.Println("Already in file")
						alreadyProcessed = true
						break
					}
				}
				if !alreadyProcessed {
					defaultFile := strings.TrimSuffix(fi.Name(), ".properties")
					properties[cartridge][defaultFile] = map[string]map[string]string{}
					properties[cartridge][defaultFile]["default"] = map[string]string{}
					filename := filepath.Join(path, fi.Name())

					//fmt.Println(fi.Name())
					file, _ := os.Open(filename)
					if err != nil {
						log.Fatal(err)
					}
					defer file.Close()
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						if !strings.HasPrefix(scanner.Text(), "#") && len(scanner.Text()) > 0 {
							property := strings.SplitN(scanner.Text(), "=", 2)
							if len(property) < 2 {
								fmt.Println("Error parsing for property: " + scanner.Text() + " // file:" + filename)
								continue
							}
							properties[cartridge][defaultFile]["default"][property[0]] = property[1]
						}
					}
					processedFiles = append(processedFiles, fi.Name())
				}
			}
		}
	}

	keys := map[string]bool{}
	for _, cartridge := range cartridges {
		for k := range properties[cartridge] {
			keys[k] = true
		}
	}
	names := make([]string, 0, len(keys))
	for k := range keys {
		names = append(names, k)
	}
	sort.Strings(names)
	// start creating file.

	f := excelize.NewFile()
	boldStyle, err := f.NewStyle(`{"font":{"bold":true}}`)
	if err != nil {
		fmt.Println(err)
	}
	for sheetName := range keys {
		f.NewSheet(sheetName)
		f.SetCellValue(sheetName, "A2", "id")
		f.SetCellStyle(sheetName, "A2", "A2", boldStyle)
		f.SetCellValue(sheetName, "B1", "default")
		f.SetCellStyle(sheetName, "B1", "B1", boldStyle)
		f.MergeCell(sheetName, "B1", "C1")
		f.SetCellValue(sheetName, "B2", "cartridge")
		f.SetCellStyle(sheetName, "B2", "B2", boldStyle)
		f.SetCellValue(sheetName, "C2", "label")
		f.SetCellStyle(sheetName, "C2", "C2", boldStyle)
		i := 4
		for _, locale := range locales {
			f.SetCellValue(sheetName, getExcelCol(i)+"1", locale)
			f.SetCellStyle(sheetName, getExcelCol(i)+"1", getExcelCol(i)+"1", boldStyle)
			f.MergeCell(sheetName, getExcelCol(i)+"1", getExcelCol(i+1)+"1")
			f.SetCellValue(sheetName, getExcelCol(i)+"2", "cartridge")
			f.SetCellStyle(sheetName, getExcelCol(i)+"2", getExcelCol(i)+"2", boldStyle)
			f.SetCellValue(sheetName, getExcelCol(i+1)+"2", "label")
			f.SetCellStyle(sheetName, getExcelCol(i+1)+"2", getExcelCol(i+1)+"2", boldStyle)
			i += 2
		}

		type ExcelEntry struct {
			cartridge string
			label     string
		}
		entries := map[string]map[string]ExcelEntry{} // map id locale Excel Entry
		for _, cartridge := range cartridges {
			saveLocale := func(locale string) {
				_, exists := properties[cartridge][sheetName]
				if exists {
					propertiesList, exists := properties[cartridge][sheetName][locale]
					if exists {
						for property, label := range propertiesList {
							_, exists := entries[property]
							if !exists {
								entries[property] = map[string]ExcelEntry{}
							}
							entries[property][locale] = ExcelEntry{cartridge: cartridge, label: label}
						}

					}
				}
			}
			saveLocale("default")
			// TODO: default
			for _, locale := range locales {
				saveLocale(locale)
			}
		}

		// here our entries map for the current file is ready. We can add it to Excel
		j := 3 // row
		for entry, labels := range entries {
			f.SetCellValue(sheetName, "A"+strconv.Itoa(j), entry)
			// save default
			_, exists := labels["default"]
			if exists {
				f.SetCellValue(sheetName, "B"+strconv.Itoa(j), labels["default"].cartridge)
				f.SetCellValue(sheetName, "C"+strconv.Itoa(j), labels["default"].label)
			}
			i := 0
			for _, locale := range locales {
				_, exists := labels[locale]
				if exists {
					f.SetCellValue(sheetName, getExcelCol(i+4)+strconv.Itoa(j), labels[locale].cartridge)
					f.SetCellValue(sheetName, getExcelCol(i+5)+strconv.Itoa(j), labels[locale].label)
				}
				i += 2
			}
			j++
		}
		for k := range keys {
			names = append(names, k)
		}
	}
	f.DeleteSheet("Sheet1")
	if err := f.SaveAs("properties.xlsx"); err != nil {
		fmt.Println(err)
	}
}
