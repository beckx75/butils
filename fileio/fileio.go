package fileio

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"github.com/tealeg/xlsx/v3"
)

type FileSettings struct {
	Sheetname string
	HeaderRow int

	IdCol      string
	NameCol    string
	ParentCol  string
	ChildCol   string
	JsonParent string

	FieldSep string
}

func ReadFile(fp string, fs *FileSettings, srcIsUtf8 bool, parentAttr string) ([][]string, error) {
	ext := strings.ToLower(filepath.Ext(fp))
	if (ext == ".xlsx") || (ext == ".xlsm") {
		return ReadExlFile(fp, fs.Sheetname)
	} else if ext == ".csv" {
		return ReadCsvFile(fp, fs.FieldSep, srcIsUtf8)
	} else if ext == ".json" {
		return ReadJsonFile(fp, srcIsUtf8, parentAttr)
	} else {
		return nil, fmt.Errorf("unsupported filetype %s", ext)
	}
}

func ReadExlFile(fp string, sheetname string) ([][]string, error) {
	wb, err := xlsx.OpenFile(fp)
	if err != nil {
		return nil, err
	}
	sh, ok := wb.Sheet[sheetname]
	if !ok {
		sh = wb.Sheets[0]
		if sheetname != "" {
			return nil, fmt.Errorf("sheet '%s' does not exist...\n", sheetname)
		}
	}
	defer sh.Close()

	maxrow := sh.MaxRow
	maxcol := sh.MaxCol
	rec := [][]string{}
	for i := 0; i < maxrow; i++ {
		row := []string{}
		for k := 0; k < maxcol; k++ {
			thecell, err := sh.Cell(i, k)
			if err != nil {
				panic(err)
			}
			row = append(row, thecell.String())
		}
		rec = append(rec, row)
	}
	return rec, nil
}

func WriteExlFile(fp string, sheetname string, rec [][]string) error {
	wb := xlsx.NewFile()
	if sheetname == "" {
		sheetname = "Tabelle1"
	}
	sh, err := wb.AddSheet(sheetname)
	if err != nil {
		return err
	}
	for i := 0; i < len(rec); i++ {
		row := sh.AddRow()
		for k := 0; k < len(rec[i]); k++ {
			cell := row.AddCell()
			cell.SetFormat(`@`)
			cell.SetString(rec[i][k])
		}
	}
	return wb.Save(fp)
}

func WriteExlSheets(fp string, sheetnames []string, recs [][][]string) error {
	wb := xlsx.NewFile()
	if len(sheetnames) == 0 {
		for i := 1; i < len(recs); i++ {
			sheetnames = append(sheetnames, fmt.Sprintf("Tabelle%d", i))
		}
	}
	for i := 0; i < len(recs); i++ {
		err := writeSheetsToWorkbook(wb, sheetnames[i], recs[i])
		if err != nil {
			return err
		}
	}
	return wb.Save(fp)
}

func writeSheetsToWorkbook(wb *xlsx.File, sheetname string, rec [][]string) error {
	sh, err := wb.AddSheet(sheetname)
	if err != nil {
		return err
	}
	for i := 0; i < len(rec); i++ {
		row := sh.AddRow()
		for k := 0; k < len(rec[i]); k++ {
			cell := row.AddCell()
			cell.SetFormat(`@`)
			cell.SetString(rec[i][k])
		}
	}
	return nil
}

func ReadCsvFile(fp string, seperator string, fileEncIsUtf8 bool) ([][]string, error) {
	rec := [][]string{}
	file, err := os.Open(fp)
	if err != nil {
		return rec, err
	}
	defer file.Close()
	var sep rune
	if seperator == "" {
		sep = ';'
	} else {
		sep = []rune(seperator)[0]
	}
	r := csv.NewReader(file)
	r.Comma = sep
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	rec, err = r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	if fileEncIsUtf8 {
		return rec, nil
	} else {
		nr := [][]string{}
		cp1252Decoder := charmap.Windows1252.NewDecoder()
		for _, row := range rec {
			r := []string{}
			for _, c := range row {
				nc, err := cp1252Decoder.String(c)
				if err != nil {
					return nil, err
				}
				r = append(r, nc)
			}
			nr = append(nr, r)
		}
		return nr, nil
	}
}

func WriteCsvFile(rec [][]string, seperator string, fp string, encValues bool) error {
	var sep rune
	if seperator == "" {
		sep = ';'
	} else {
		sep = []rune(seperator)[0]
	}
	file, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	w.Comma = sep
	if !encValues {
		return w.WriteAll(rec)
	}

	cp1252Encoder := charmap.Windows1252.NewEncoder()
	for _, row := range rec {
		cprow := []string{}
		for _, c := range row {
			nc, err := cp1252Encoder.String(c)
			if err != nil {
				return err
			}
			cprow = append(cprow, nc)
		}
		err = w.Write(cprow)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadJsonFile(fp string, fileEncIsUtf8 bool, parentAttr string) ([][]string, error) {
	return JsonToSpreadsheet(fp, parentAttr)
}
