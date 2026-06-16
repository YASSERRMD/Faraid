package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-pdf/fpdf"

	"github.com/YASSERRMD/Faraid/internal/core/solver"
)

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")

	var req solveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	c, m, err := toCase(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := solver.Solve(c, m)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	dto := toSolveResult(result)

	switch format {
	case "json":
		w.Header().Set("Content-Disposition", `attachment; filename="faraid.json"`)
		writeJSON(w, http.StatusOK, dto)
	case "pdf":
		pdf, err := resultPDF(dto)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not render pdf")
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="faraid.pdf"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(pdf)
	default:
		writeError(w, http.StatusBadRequest, "format must be json or pdf")
	}
}

// fixedPDFDate makes the PDF output deterministic by pinning the creation date.
var fixedPDFDate = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

// resultPDF renders a printable result document. Compression is disabled and
// the creation date is fixed so the output is deterministic and its text is
// visible in the bytes for structural testing. The document is laid out left to
// right with English labels; rendering Arabic labels right to left requires
// embedding an Arabic font, which is left to a later enhancement.
func resultPDF(res solveResultDTO) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetCompression(false)
	pdf.SetCreationDate(fixedPDFDate)
	pdf.SetModificationDate(fixedPDFDate)
	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(0, 10, "Faraid Inheritance Result")
	pdf.Ln(12)

	pdf.SetFont("Helvetica", "", 11)
	pdf.Cell(0, 7, "School: "+res.Madhhab)
	pdf.Ln(7)
	pdf.Cell(0, 7, "Base of the problem: "+strconv.FormatInt(res.Base, 10))
	pdf.Ln(7)
	pdf.Cell(0, 7, "Distributable estate: "+res.Distributable)
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(0, 8, "Shares")
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(70, 7, "Heir", "1", 0, "L", false, 0, "")
	pdf.CellFormat(20, 7, "Count", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, "Fraction", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 7, "Amount", "1", 1, "C", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	for _, sh := range res.Shares {
		pdf.CellFormat(70, 7, sh.Relation, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 7, strconv.Itoa(sh.Count), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 7, sh.Fraction, "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 7, sh.Amount, "1", 1, "C", false, 0, "")
	}
	pdf.Ln(6)

	if res.SpecialCase != "" {
		pdf.SetFont("Helvetica", "I", 10)
		pdf.Cell(0, 7, "Special case: "+res.SpecialCase)
		pdf.Ln(8)
	}

	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(0, 8, "Derivation")
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 9)
	for _, step := range res.Derivation {
		line := "[" + step.Stage + "] "
		if step.Relation != "" {
			line += step.Relation + ": "
		}
		line += step.Detail
		if step.Fraction != "" {
			line += " = " + step.Fraction
		}
		if step.Reference != "" {
			line += " (" + step.Reference + ")"
		}
		pdf.MultiCell(0, 5, line, "", "L", false)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
