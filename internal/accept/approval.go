package accept

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Approval records a human decision and binds the exact approved inputs.
// Its provenance is declarative, like the existing prompt and Git identity;
// the command validates bindings, not the approver's authenticity.
type Approval struct {
	Type        string `yaml:"type"`
	Decision    string `yaml:"decision"`
	Candidate   string `yaml:"candidate"`
	Base        string `yaml:"base"`
	BriefSHA256 string `yaml:"brief-sha256"`
	ApprovedBy  string `yaml:"approved-by"`
	Source      string `yaml:"approval-source"`
	Statement   string `yaml:"statement"`
	RecordedAt  string `yaml:"recorded-at"`
}

func (p *Candidate) approvalData(name string) ([]byte, error) {
	data, e := os.ReadFile(name)
	if e != nil {
		return nil, fmt.Errorf("read human approval: %w", e)
	}
	var a Approval
	if e = decode(data, &a); e != nil {
		return nil, fmt.Errorf("human approval: %w", e)
	}
	digest := sha256.Sum256([]byte(p.brief))
	if a.Type != "acceptance-approval" || a.Decision != "accept" || a.Candidate != p.request.Candidate || a.Base != p.request.Base || a.BriefSHA256 != hex.EncodeToString(digest[:]) {
		return nil, fmt.Errorf("human approval must accept this exact candidate, base and complete brief hash")
	}
	if incompleteApprovalField(a.ApprovedBy) || incompleteApprovalField(a.Source) || incompleteApprovalField(a.Statement) {
		return nil, fmt.Errorf("human approval requires approver, approval source and actual statement")
	}
	if _, e = time.Parse(time.RFC3339, a.RecordedAt); e != nil {
		return nil, fmt.Errorf("human approval requires a recorded-at RFC3339 timestamp")
	}
	return data, nil
}

// VerifyApproval is a read-only approval-record preflight.
func (p *Candidate) VerifyApproval(name string) error { _, e := p.approvalData(name); return e }

// ConfirmRecorded executes a recorded human decision without a terminal prompt.
// It never decides acceptance for the human, authenticates a statement or pushes.
func (p *Candidate) ConfirmRecorded(name string, output io.Writer) (string, error) {
	original, e := p.approvalData(name)
	if e != nil {
		return "", e
	}
	fmt.Fprintf(output, "Executing recorded human approval for candidate %s at base %s. Provenance and verification are recorded claims, not authenticated human presence. No push will be performed.\n", p.request.Candidate, p.request.Base)
	fresh, e := p.approvalData(name)
	if e != nil {
		return "", e
	}
	if !bytes.Equal(original, fresh) {
		return "", fmt.Errorf("human approval changed before execution; obtain and record a fresh decision")
	}
	return p.complete(string(original))
}

func incompleteApprovalField(value string) bool {
	value = strings.TrimSpace(value)
	return value == "" || strings.HasPrefix(value, "REPLACE_WITH_") || strings.Contains(value, "<REQUIRED>")
}
