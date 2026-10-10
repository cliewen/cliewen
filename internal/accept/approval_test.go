package accept

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func approvalFixture(t *testing.T, r Request, change func(*Approval)) string {
	t.Helper()
	data, e := os.ReadFile(r.BriefPath)
	if e != nil {
		t.Fatal(e)
	}
	digest := sha256.Sum256(data)
	a := Approval{Type: "acceptance-approval", Decision: "accept", Candidate: r.Candidate, Base: r.Base, BriefSHA256: hex.EncodeToString(digest[:]), ApprovedBy: "Fixture observer", Source: "Synthetic test conversation", Statement: "accept this fixture candidate", RecordedAt: time.Now().UTC().Format(time.RFC3339)}
	if change != nil {
		change(&a)
	}
	body, e := yaml.Marshal(a)
	if e != nil {
		t.Fatal(e)
	}
	name := filepath.Join(r.Root, ".git/clue/approval.yaml")
	write(t, filepath.Dir(name), filepath.Base(name), string(body))
	return name
}

func TestAC237_IntegrationPositive_RecordedDecisionRetainsExactTreeParentsAndApproval(t *testing.T) {
	r := fixture(t, nil)
	name := approvalFixture(t, r, nil)
	original, _ := os.ReadFile(name)
	p, e := Check(r)
	if e != nil {
		t.Fatal(e)
	}
	if e = p.VerifyApproval(name); e != nil {
		t.Fatal(e)
	}
	if got := mustGit(t, r.Root, "rev-parse", "HEAD"); got != r.Base {
		t.Fatal("approval preflight integrated")
	}
	var out bytes.Buffer
	commit, e := p.ConfirmRecorded(name, &out)
	if e != nil {
		t.Fatal(e)
	}
	if got := mustGit(t, r.Root, "show", "-s", "--format=%P", commit); got != r.Base+" "+r.Candidate {
		t.Fatal(got)
	}
	if got := mustGit(t, r.Root, "rev-parse", commit+"^{tree}"); got != mustGit(t, r.Root, "rev-parse", r.Candidate+"^{tree}") {
		t.Fatal("tree differs")
	}
	message := mustGit(t, r.Root, "show", "-s", "--format=%B", commit)
	if !strings.Contains(message, strings.TrimSpace(string(original))) || !strings.Contains(message, p.brief) {
		t.Fatal("approval or brief lost")
	}
	if !strings.Contains(out.String(), "recorded claims") && !strings.Contains(out.String(), "not authenticated human presence") {
		t.Fatal("trust boundary omitted")
	}
}

func TestAC237_IntegrationNegative_MissingStaleMalformedOrRejectedRecordsNeverIntegrate(t *testing.T) {
	cases := map[string]func(*Approval){"candidate": func(a *Approval) { a.Candidate = strings.Repeat("1", 40) }, "base": func(a *Approval) { a.Base = strings.Repeat("1", 40) }, "brief": func(a *Approval) { a.BriefSHA256 = strings.Repeat("0", 64) }, "decision": func(a *Approval) { a.Decision = "reject" }, "type": func(a *Approval) { a.Type = "other" }, "observer": func(a *Approval) { a.ApprovedBy = "" }, "source": func(a *Approval) { a.Source = "" }, "statement": func(a *Approval) { a.Statement = "" }, "time": func(a *Approval) { a.RecordedAt = "yesterday" }}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			r := fixture(t, nil)
			record := approvalFixture(t, r, change)
			p, e := Check(r)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = p.ConfirmRecorded(record, io.Discard); e == nil {
				t.Fatal("invalid record accepted")
			}
			if got := mustGit(t, r.Root, "rev-parse", "HEAD"); got != r.Base {
				t.Fatal("ref moved")
			}
		})
	}
	r := fixture(t, nil)
	record := approvalFixture(t, r, nil)
	p, e := Check(r)
	if e != nil {
		t.Fatal(e)
	}
	body, _ := os.ReadFile(record)
	for _, bad := range []string{string(body) + "unknown: ignored\n", string(body) + "---\ndecision: accept\n", ""} {
		write(t, filepath.Dir(record), filepath.Base(record), bad)
		if e = p.VerifyApproval(record); e == nil {
			t.Fatal("malformed YAML accepted")
		}
	}
	if _, e = p.ConfirmRecorded(record+".missing", io.Discard); e == nil {
		t.Fatal("missing record accepted")
	}
}

func TestAC238_IntegrationPositive_ApprovalPreflightAndCancellationRemainNonIntegrating(t *testing.T) {
	r := fixture(t, nil)
	record := approvalFixture(t, r, nil)
	p, e := Check(r)
	if e != nil {
		t.Fatal(e)
	}
	if e = p.VerifyApproval(record); e != nil {
		t.Fatal(e)
	}
	if _, e = p.Confirm(strings.NewReader("cancel\n"), io.Discard, true); e == nil {
		t.Fatal("cancellation accepted")
	}
	if got := mustGit(t, r.Root, "rev-parse", "HEAD"); got != r.Base {
		t.Fatal("preflight or cancellation integrated")
	}
}

func TestAC238_IntegrationNegative_ChangedBriefOrDirtyBaseCannotReplayRecordedApproval(t *testing.T) {
	for _, name := range []string{"brief", "dirty", "base"} {
		t.Run(name, func(t *testing.T) {
			r := fixture(t, nil)
			record := approvalFixture(t, r, nil)
			p, e := Check(r)
			if e != nil {
				t.Fatal(e)
			}
			switch name {
			case "brief":
				write(t, filepath.Dir(r.BriefPath), filepath.Base(r.BriefPath), stringMustRead(t, r.BriefPath)+"\nChanged verification claim.\n")
			case "dirty":
				write(t, r.Root, "app.txt", "new local edit")
			case "base":
				mustGit(t, r.Root, "commit", "--allow-empty", "-m", "Advance accepted base")
			}
			before := mustGit(t, r.Root, "rev-parse", "HEAD")
			if _, e = p.ConfirmRecorded(record, io.Discard); e == nil {
				t.Fatal("stale inputs accepted")
			}
			if got := mustGit(t, r.Root, "rev-parse", "HEAD"); got != before {
				t.Fatal("stale record moved ref")
			}
		})
	}
}

func stringMustRead(t *testing.T, name string) string {
	t.Helper()
	data, e := os.ReadFile(name)
	if e != nil {
		t.Fatal(e)
	}
	return string(data)
}

type approvalOutputMutation struct {
	action func()
	done   bool
}

func (w *approvalOutputMutation) Write(p []byte) (int, error) {
	if !w.done {
		w.done = true
		w.action()
	}
	return len(p), nil
}
func TestAC238_IntegrationNegative_ApprovalChangedAtExecutionIsRejected(t *testing.T) {
	r := fixture(t, nil)
	record := approvalFixture(t, r, nil)
	p, e := Check(r)
	if e != nil {
		t.Fatal(e)
	}
	out := &approvalOutputMutation{action: func() {
		body := stringMustRead(t, record)
		write(t, filepath.Dir(record), filepath.Base(record), strings.Replace(body, "Fixture observer", "Other observer", 1))
	}}
	if _, e = p.ConfirmRecorded(record, out); e == nil || !strings.Contains(e.Error(), "changed") {
		t.Fatalf("changed record accepted: %v", e)
	}
	if got := mustGit(t, r.Root, "rev-parse", "HEAD"); got != r.Base {
		t.Fatal("changed approval moved ref")
	}
}

func TestAC237_IntegrationNegative_ApprovalTemplatePlaceholdersAreNotProvenance(t *testing.T) {
	r := fixture(t, nil)
	record := approvalFixture(t, r, func(a *Approval) { a.Statement = "REPLACE_WITH_ACTUAL_HUMAN_APPROVAL" })
	p, e := Check(r)
	if e != nil {
		t.Fatal(e)
	}
	if e = p.VerifyApproval(record); e == nil {
		t.Fatal("placeholder statement accepted")
	}
	if !incompleteApprovalField(" <REQUIRED> ") {
		t.Fatal("incomplete marker accepted")
	}
}
