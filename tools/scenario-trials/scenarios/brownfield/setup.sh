set -e
mkdir -p /home/node/work && cd /home/node/work
git init -q
mkdir -p openspec/specs/export
printf 'schema: spec-driven\n' > openspec/config.yaml
cat > openspec/specs/export/spec.md <<'SPEC'
# Export Specification

## Purpose
Export a report as CSV so it can be opened in a spreadsheet.

## Requirements

### Requirement: Export writes a header row
The system SHALL write a header row before any data row.

#### Scenario: Header comes first [EX-01]
Test-type: Unit
- **WHEN** an export runs
- **THEN** the first line holds the column names

### Requirement: Export writes one row per record
The system SHALL write exactly one line for each record in the report.

#### Scenario: One line per record [EX-02]
Test-type: Unit
- **WHEN** a report with two records is exported
- **THEN** the file has a header line and two data lines

#### Scenario: Empty report [EX-03]
Test-type: Unit
- **WHEN** a report with no records is exported
- **THEN** the file holds only the header line
SPEC
printf '# Sales report\n\nA small report tool. Its specifications are in openspec/.\n' > README.md
clue init . >/dev/null
git add -A
git commit -qm init
