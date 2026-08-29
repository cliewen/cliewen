package skills

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestAC137_UnitPositive_GenerationProducesRoutedStandaloneSkillDirectories(t *testing.T) {
	root := t.TempDir()
	if err := Write(root); err != nil {
		t.Fatal(err)
	}
	drifts, err := Check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(drifts) != 0 {
		t.Fatalf("freshly generated skills drifted: %v", drifts)
	}

	for _, file := range mustRender(t) {
		agentPath := filepath.Join(root, ".agents", "skills", filepath.FromSlash(file.relativePath))
		templatePath := filepath.Join(root, "internal", "scaffold", "templates", "skills", filepath.FromSlash(file.relativePath))
		agent, agentErr := os.ReadFile(agentPath)
		if agentErr != nil {
			t.Fatal(agentErr)
		}
		embedded, embeddedErr := os.ReadFile(templatePath)
		if embeddedErr != nil {
			t.Fatal(embeddedErr)
		}
		if string(agent) != string(embedded) {
			t.Fatalf("%s differs between generated output trees", file.relativePath)
		}
		if filepath.Base(file.relativePath) == "skill.md" && !strings.Contains(string(agent), "\ncliewen-skill: true\n") {
			t.Fatalf("%s carries no Cliewen ownership marker", file.relativePath)
		}
	}

	for name, definition := range skillDefinitions {
		entrypoint := mustRenderFile(t, path.Join(name, "skill.md"))
		for _, route := range definition.routes {
			link := "[" + route.heading + "](references/" + route.file + ")"
			if !strings.Contains(entrypoint, route.condition) || !strings.Contains(entrypoint, link) {
				t.Errorf("%s/skill.md does not route %q to %s under its condition", name, route.heading, route.file)
			}
			if strings.Contains(entrypoint, "## "+route.heading+"\n") {
				t.Errorf("%s/skill.md eagerly carries deferred section %q", name, route.heading)
			}
			if content := mustRenderFile(t, path.Join(name, "references", route.file)); !strings.HasPrefix(content, "## "+route.heading+"\n") {
				t.Errorf("%s reference %s does not carry routed section %q", name, route.file, route.heading)
			}
		}
	}
}

func TestAC137_UnitNegative_EntrypointAndReferenceDriftIsRejected(t *testing.T) {
	tests := map[string]func(*testing.T, string){
		"changed entrypoint": func(t *testing.T, root string) {
			target := filepath.Join(root, ".agents", "skills", "clue-delta", "skill.md")
			if err := os.WriteFile(target, []byte("edited generated output\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
		"missing reference": func(t *testing.T, root string) {
			target := filepath.Join(root, ".agents", "skills", "clue-delta", "references", "review-boundary.md")
			if err := os.Remove(target); err != nil {
				t.Fatal(err)
			}
		},
		"unexpected reference": func(t *testing.T, root string) {
			target := filepath.Join(root, ".agents", "skills", "clue-delta", "references", "manual.md")
			if err := os.WriteFile(target, []byte("not generated\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
		"changed reference in template tree": func(t *testing.T, root string) {
			target := filepath.Join(root, "internal", "scaffold", "templates", "skills", "clue-delta", "references", "change-loop.md")
			if err := os.WriteFile(target, []byte("edited generated output\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		},
	}

	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			if err := Write(root); err != nil {
				t.Fatal(err)
			}
			mutate(t, root)
			drifts, err := Check(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(drifts) == 0 {
				t.Fatal("expected generated skill drift to be rejected")
			}
			if !strings.Contains(drifts[0].Path, "clue-delta") {
				t.Fatalf("drift did not name the affected skill: %v", drifts)
			}
		})
	}
}

func TestAC147_UnitPositive_EveryLifecycleSkillCarriesSubjectTypedDecisionGuidance(t *testing.T) {
	for name := range skillDefinitions {
		content := mustRenderFile(t, path.Join(name, "references", "decision-records.md"))
		for _, want := range []string{"future-shaping choice", "ADR for software or corpus architecture", "PDR for how the project", "IDR for implementation", "timeless and compact"} {
			if !strings.Contains(content, want) {
				t.Errorf("%s decision guidance is missing %q", name, want)
			}
		}
	}
}

func TestAC147_UnitNegative_GeneratedSkillsOmitLegacyDecisionRouting(t *testing.T) {
	for name := range skillDefinitions {
		content := mustRenderFile(t, path.Join(name, "references", "decision-records.md"))
		for _, legacy := range []string{"Route every decision by reversal cost", "docs/decisions/log.md"} {
			if strings.Contains(content, legacy) {
				t.Errorf("%s decision guidance retains legacy contract %q", name, legacy)
			}
		}
	}
}

func TestAC146_UnitPositive_ExtractionClassifiesDecisionsAndInventoriesLegacyRows(t *testing.T) {
	target := mustRenderFile(t, "clue-extract/references/target-contract.md")
	mapping := mustRenderFile(t, "clue-extract/mappings/madr.md")
	for _, want := range []string{"decisions route by subject", "ADR for architecture", "PDR for project/process", "IDR for implementation", "legacy decision log", "never guesses or silently drops"} {
		if !strings.Contains(target, want) {
			t.Errorf("extraction target contract is missing %q", want)
		}
	}
	for _, want := range []string{"decision's subject", "ADR-xxx", "PDR-xxx", "IDR-xxx", "legacy decision log table", "never infer one mechanically"} {
		if !strings.Contains(mapping, want) {
			t.Errorf("MADR mapping is missing %q", want)
		}
	}
}

func TestAC146_UnitNegative_ExtractionDoesNotGuessOrRouteEveryRecordToADR(t *testing.T) {
	target := mustRenderFile(t, "clue-extract/references/target-contract.md")
	mapping := mustRenderFile(t, "clue-extract/mappings/madr.md")
	for _, legacy := range []string{"their ADR/PDR record type is already the high-cost route", "the number becomes the `ADR-xxx` ID"} {
		if strings.Contains(target, legacy) || strings.Contains(mapping, legacy) {
			t.Errorf("extraction retains legacy routing %q", legacy)
		}
	}
}

func TestSanity_EveryMappingSourceHasAGeneratedCounterpart(t *testing.T) {
	mappingsDir := filepath.Join("source", "resources", "clue-extract", "mappings")

	rendered := map[string]bool{}
	for _, file := range mustRender(t) {
		rendered[file.relativePath] = true
	}

	found := 0
	err := filepath.WalkDir(mappingsDir, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		found++
		rel, relErr := filepath.Rel(mappingsDir, filePath)
		if relErr != nil {
			return relErr
		}
		want := path.Join("clue-extract", "mappings", filepath.ToSlash(rel))
		if !rendered[want] {
			t.Errorf("mapping source %s has no generated counterpart %s", filepath.ToSlash(rel), want)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if found == 0 {
		t.Fatal("no mapping sources found under source/resources/clue-extract/mappings")
	}
}

func TestSanity_CommittedSkillsMatchCanonicalSources(t *testing.T) {
	root := filepath.Join("..", "..")
	drifts, err := Check(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, drift := range drifts {
		t.Error(drift)
	}
}

func TestSanity_EveryManagedSkillAppearsInBothRoutingHubs(t *testing.T) {
	for _, hub := range []string{
		filepath.Join("..", "..", "AGENTS.md"),
		filepath.Join("..", "scaffold", "templates", "AGENTS.md"),
	} {
		content, err := os.ReadFile(hub)
		if err != nil {
			t.Fatal(err)
		}
		for _, name := range skillNames {
			row := "| [`" + name + "`](.agents/skills/" + name + "/skill.md) |"
			if !strings.Contains(string(content), row) {
				t.Errorf("%s is missing managed skill row %q", filepath.ToSlash(hub), name)
			}
		}
	}
}

func TestSanity_RoutedSkillSectionsUseLocalReferencePointers(t *testing.T) {
	verify := mustRenderFile(t, "clue-verify/references/verification-checklist.md")
	for _, want := range []string{
		"under [Change routing](change-scope-and-tiers.md)",
		"satisfy the [Review boundary](review-boundary.md)",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify verification reference does not carry accurate local pointer %q", want)
		}
	}

	extract := mustRenderFile(t, "clue-extract/references/source-mappings.md")
	for _, want := range []string{"[openspec.md](../mappings/openspec.md)", "[madr.md](../mappings/madr.md)"} {
		if !strings.Contains(extract, want) {
			t.Errorf("clue-extract source mapping reference does not resolve local mapping %q", want)
		}
	}
}

func TestSanity_SpecFirstPauseReportsProposalAndImplementationStatus(t *testing.T) {
	delta := mustRenderSkill(t, "clue-delta/skill.md")
	for _, want := range []string{
		"report briefly what the proposal says and what implementation involves",
		"ask whether implementation should begin — the proposal is already committed, pushed, and visible on the draft PR",
	} {
		if !strings.Contains(delta, want) {
			t.Errorf("clue-delta/skill.md does not carry complete spec-first pause instruction %q", want)
		}
	}
}

// AC-054 is an extraction-guidance criterion: what it promises is that the
// canonical clue-extract skill instructs criterion-level phasing, so the
// rendered skill text is its evidence. The mechanism that guidance relies on
// is proven separately against the validator — AC-045 for the Human class and
// AC-046 for per-criterion @draft, both in internal/corpus.
func TestAC054_UnitPositive_ExtractionSupportsCriterionLevelPhasing(t *testing.T) {
	extract := mustRenderSkill(t, "clue-extract/skill.md")
	for _, want := range []string{
		"Whole-file draft phasing remains available",
		"tag each genuinely not-yet-proven criterion `@draft`",
		"`Test-type: Human` criterion is already proven by naming it in the pull request acceptance brief",
		"supported Go, per-executable JVM, or Cucumber evidence",
		"A capability is therefore not the smallest activation unit",
	} {
		if !strings.Contains(extract, want) {
			t.Errorf("clue-extract/skill.md does not carry criterion-level phasing rule %q", want)
		}
	}
}

func TestAC142_UnitPositive_GeneratedUpgradeSkillRoutesAHumanAuthorizedCoordinatedUpgradeAsSimpleWork(t *testing.T) {
	upgrade := mustRenderSkill(t, "clue-upgrade/skill.md")
	for _, want := range []string{
		"Run `clue latest`",
		"`### Migration` section",
		"Ask the human whether to upgrade now or later",
		"Do nothing to the repository until they explicitly choose now",
		"**An upgrade is simple work.**",
		"Recommend `Recommended route: simple`",
		"needs no CH identity, workspace, plan declaration, digest, acceptance brief, or mandatory agentic review",
		"How many files the migration rewrites, and whether corpus or skill paths are among them, never makes the route full",
		"Escalate only on a semantic discovery",
		"recommend the full loop for that decision on its own terms",
		"make the repository green and create a branch",
		"resolve every finding and notice — including those no command may repair — except the explicitly non-blocking `MIG-009` competing-wall notice",
		"reconcile that job by hand after applying",
		"Never merge it",
	} {
		if !strings.Contains(upgrade, want) {
			t.Errorf("clue-upgrade/skill.md does not carry upgrade boundary %q", want)
		}
	}
}

func TestAC142_UnitNegative_GeneratedUpgradeSkillDoesNotInventAPlatformRouteSelfAuthorizeOrBindSimpleWorkToTheFullLoopBoundary(t *testing.T) {
	upgrade := mustRenderSkill(t, "clue-upgrade/skill.md")
	for _, forbidden := range []string{
		"curl -fsSL https://cliewen.dev/install.sh | sh",
		"irm https://cliewen.dev/install.ps1 | iex",
		"automatically upgrade",
		"merge the pull request",
		"mark the upgrade's pull request ready under the review boundary",
	} {
		if strings.Contains(upgrade, forbidden) {
			t.Errorf("clue-upgrade/skill.md exceeds its route and authority boundary with %q", forbidden)
		}
	}
}

// The upgrade skill's router must reach its workflow — where the route is
// stated — before the generic tier text, or the agent recommends a route from
// text that cannot know the contract change already happened upstream.
func TestAC142_UnitNegative_UpgradeRouterDoesNotSendTheAgentToGenericTiersForItsRoute(t *testing.T) {
	upgrade := mustRenderSkill(t, "clue-upgrade/skill.md")
	workflow := strings.Index(upgrade, "references/upgrade-workflow.md")
	tiers := strings.Index(upgrade, "references/change-scope-and-tiers.md")
	if workflow < 0 || tiers < 0 {
		t.Fatalf("clue-upgrade/skill.md does not route both the upgrade workflow and change routing")
	}
	if workflow > tiers {
		t.Errorf("clue-upgrade/skill.md routes generic change routing before its own workflow, so the route is decided without it")
	}
	if strings.Contains(upgrade, "before recommending its route") {
		t.Errorf("clue-upgrade/skill.md still sends the agent to the generic tier text to recommend the upgrade's route")
	}
}

func TestAC148_UnitPositive_GeneratedUpgradeSkillPreviewsCarrierDriftBeforeAuthorization(t *testing.T) {
	upgrade := mustRenderSkill(t, "clue-upgrade/skill.md")
	latest := strings.Index(upgrade, "Run `clue latest`")
	preview := strings.Index(upgrade, "Immediately after `clue latest`, run `clue migrate` without `--apply` as a preview")
	approval := strings.Index(upgrade, "Ask the human whether to upgrade now or later")
	if latest < 0 || preview < 0 || approval < 0 {
		t.Fatalf("clue-upgrade/skill.md does not carry the release, preview, and authorization sequence")
	}
	if !(latest < preview && preview < approval) {
		t.Fatalf("clue-upgrade/skill.md does not preview repository state after availability and before authorization")
	}
	for _, want := range []string{
		"even when `clue latest` says the installed release is the newest",
		"checks whether this repository's managed carriers match the installed binary",
		"call the repository current only when the preview reports no changes and no findings",
		"use the already reviewed `clue migrate` preview",
	} {
		if !strings.Contains(upgrade, want) {
			t.Errorf("clue-upgrade/skill.md does not carry carrier-drift preview rule %q", want)
		}
	}
}

func TestAC148_UnitNegative_GeneratedUpgradeSkillDoesNotConflateReleaseAvailabilityWithCarrierCurrency(t *testing.T) {
	upgrade := mustRenderSkill(t, "clue-upgrade/skill.md")
	for _, forbidden := range []string{
		"If no newer release is available, stop",
		"call the repository current because `clue latest` says the installed release is the newest",
		"Ask the human whether to upgrade now or later before running `clue migrate`",
	} {
		if strings.Contains(upgrade, forbidden) {
			t.Errorf("clue-upgrade/skill.md still conflates release availability with carrier currency through %q", forbidden)
		}
	}
}

func TestAC054_UnitNegative_ExtractionRejectsCapabilityOnlyPhasing(t *testing.T) {
	extract := mustRenderSkill(t, "clue-extract/skill.md")
	for _, stale := range []string{
		"Activation is per criteria file, not per criterion",
		"a capability is the smallest unit a phasing change can take",
		"Every active acceptance criterion has positive and negative tests, or its capability honestly stays `draft`",
	} {
		if strings.Contains(extract, stale) {
			t.Errorf("clue-extract/skill.md still carries capability-only phasing rule %q", stale)
		}
	}
}

func TestAC055_UnitPositive_AnalysisQualifiesEnvironmentAndPopulationEvidence(t *testing.T) {
	analysis := mustRenderSkill(t, "clue-analysis/skill.md")
	for _, want := range []string{
		"either a clean disposable environment or a prepared environment",
		"A clean result supports onboarding reproducibility only when it has no local prerequisites",
		"any local prerequisite, documented or not, makes the result prepared",
		"A prepared result names its prerequisites and establishes only what that prepared environment demonstrated",
		"versioned corpus and population, eligibility rules, exclusions and their reasons, sampling or repetition method, uncertainty",
		"deterministic-versus-quality boundary",
		"Do not turn an environment-sensitive quality claim into a deterministic acceptance criterion",
		"name the governance or process changes it introduces",
		"do not describe scaffolding as neutral",
	} {
		if !strings.Contains(analysis, want) {
			t.Errorf("clue-analysis/skill.md does not qualify analysis evidence with %q", want)
		}
	}
}

// The negative direction rejects the weaker forms this obligation can decay
// into, not invented opposites: an unqualified skill never claims that a
// percentage needs no population, it simply stops asking for one. Each string
// below is text the skill would carry again if the rule were relaxed.
func TestAC055_UnitNegative_AnalysisRejectsUnqualifiedEvidenceClaims(t *testing.T) {
	analysis := mustRenderSkill(t, "clue-analysis/skill.md")
	for _, stale := range []string{
		// The pre-CH-080 order, in which investigation begins straight after
		// the evidence boundary with no population or adoption qualification.
		"says otherwise.\n3. Run a **spike**",
		// A prerequisite that only counts when it is undocumented.
		"no unstated local prerequisites",
		"any undocumented local prerequisite",
		// A population claim untied to the corpus version it was drawn from.
		"versioned population, eligibility rules, exclusions and their reasons",
	} {
		if strings.Contains(analysis, stale) {
			t.Errorf("clue-analysis/skill.md still permits unqualified analysis evidence claim %q", stale)
		}
	}
}

func TestAC056_UnitPositive_ExtractionRehearsesBeforeMutation(t *testing.T) {
	extract := mustRenderSkill(t, "clue-extract/skill.md")
	for _, want := range []string{
		"mandatory report-only pass",
		"under `/changes/<CH-xxx-slug>/`",
		"do not change the target source corpus, Cliewen `/docs` corpus, tests, routing, or hosted state",
		"source formats and entry points, proposed artifact mappings, preserved and minted IDs, confidence and reversal cost, test-purpose work, instruction conflicts, planned deletions, and named plan doors",
		"An unresolved conflict becomes an `open-questions.md` entry and stops before mutation",
		"Only explicit human direction begins the existing full extraction change's mutate phase",
		"digests the rehearsal into the durable extraction report under `/docs/analysis`",
	} {
		if !strings.Contains(extract, want) {
			t.Errorf("clue-extract/skill.md does not require report-only rehearsal detail %q", want)
		}
	}
}

func TestAC061_UnitPositive_ExtractionPreservesAndMintsExtendedCriterionIDs(t *testing.T) {
	extract := mustRenderSkill(t, "clue-extract/skill.md")
	for _, want := range []string{
		"keep source IDs verbatim",
		"one or more uppercase alphanumeric segments joined by single hyphens",
		"a lowercase letter suffix after its numeric portion",
		"the next numeric slot after the maximum numeric component already declared",
		"ignoring letter suffixes for the maximum",
		"an empty namespace starts at one",
		"the same source state always mints the same IDs",
		"preserved and minted mapping",
	} {
		if !strings.Contains(extract, want) {
			t.Errorf("clue-extract/skill.md does not carry extended-ID rule %q", want)
		}
	}
}

func TestAC061_UnitNegative_ExtractionRejectsIDRenumberingAndAdHocMinting(t *testing.T) {
	extract := mustRenderSkill(t, "clue-extract/skill.md")
	for _, stale := range []string{
		"renumber source IDs to fit the default AC grammar",
		"mint IDs in arbitrary order",
		"choose a fresh ID independently for each extraction run",
	} {
		if strings.Contains(extract, stale) {
			t.Errorf("clue-extract/skill.md carries an unstable extended-ID rule %q", stale)
		}
	}
}

// The negative direction rejects the pre-CH-081 shape this obligation decays
// back into, not invented opposites: a skill that permits mutation first never
// says so, it simply stops naming the checkpoint and describes the extraction
// report as a record written about work already done. Each string below is
// text the skill carried before the checkpoint existed.
func TestAC056_UnitNegative_ExtractionRejectsMutationBeforeRehearsal(t *testing.T) {
	extract := mustRenderSkill(t, "clue-extract/skill.md")
	for _, stale := range []string{
		// The apply list without the rehearsal section.
		"Apply the **Decision records**, **Repository-local conventions**, and **Review boundary** below.",
		// The durable report as a free-standing record rather than the
		// mutate phase's digest of the rehearsal.
		"`/docs/analysis`:** Record what was found",
	} {
		if strings.Contains(extract, stale) {
			t.Errorf("clue-extract/skill.md permits mutation before rehearsal through %q", stale)
		}
	}
}

func TestSanity_MethodologyContractChangesMoveEveryLiveCarrierTogether(t *testing.T) {
	for _, name := range skillNames {
		content := mustRenderSkill(t, name+"/skill.md")
		for _, want := range []string{
			"A decision that changes a methodology contract inventories every live carrier",
			"updates that complete inventory in the same change",
			"Historical analyses, completed plans, and changelog entries remain pinned history",
			"that general obligation remains agent-enforced",
		} {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not carry same-change methodology-carrier rule %q", name, want)
			}
		}
	}
}

func TestSanity_VerifyRecognizesTheCompleteEvidenceContract(t *testing.T) {
	verify := mustRenderSkill(t, "clue-verify/skill.md")
	for _, want := range []string{
		"supported Go, JVM, or Cucumber evidence",
		"positive/negative direction",
		"`(single-direction)`",
		"a genuine `Human` criterion is named in the acceptance brief as its proof",
		"an individual not-yet-proven criterion carries `@draft`",
		"an unannotated legacy criterion has its one supported reference",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not recognize evidence-model case %q", want)
		}
	}
	for _, stale := range []string{
		"Every active acceptance criterion has positive and negative tests",
		"its capability honestly stays `draft`",
	} {
		if strings.Contains(verify, stale) {
			t.Errorf("clue-verify/skill.md still carries stale evidence-model rule %q", stale)
		}
	}
}

func TestAC058_UnitPositive_GeneratedSkillsStatePerExecutableJVMContract(t *testing.T) {
	rendered := map[string]string{}
	for _, name := range skillNames {
		rendered[name+"/skill.md"] = mustRenderSkill(t, name+"/skill.md")
	}
	required := map[string][]string{
		"clue-delta/skill.md": {
			"all three evidence parts attach to the same Java or Kotlin executable",
			"`test<PREFIX><digits>[lowercase-suffix]_<Type><Direction>_<description>`",
			"class tags, comments, and unrelated methods cannot supply missing parts",
		},
		"clue-extract/skill.md": {
			"normalize each supported Java or Kotlin executable",
			"dynamic or multi-line tag expressions",
			"instead of installing an external rule or letting `clue` guess",
		},
		"clue-verify/skill.md": {
			"JVM evidence carries its AC identity, type, and direction on the same Java or Kotlin executable",
			"literal JUnit method tags or the stable named-executable form",
		},
	}
	for name, fragments := range required {
		for _, fragment := range fragments {
			if !strings.Contains(rendered[name], fragment) {
				t.Errorf("%s does not carry per-executable JVM rule %q", name, fragment)
			}
		}
	}
}

func TestAC058_UnitNegative_GeneratedSkillsRejectTheObsoleteJVMCarrier(t *testing.T) {
	for _, name := range skillNames {
		content := mustRenderSkill(t, name+"/skill.md")
		for _, stale := range []string{
			"install an ArchUnit or equivalent rule enforcing one purpose tag per test",
			"`clue` only harvests at file level",
			"JVM test files use JUnit tags harvested at file level",
		} {
			if strings.Contains(content, stale) {
				t.Errorf("%s still carries obsolete JVM evidence rule %q", name, stale)
			}
		}
	}
}

func TestUnit_ReviewBoundaryRequiresExactHostedHandoff(t *testing.T) {
	rendered := map[string]string{}
	for _, name := range skillNames {
		rendered[name+"/skill.md"] = mustRenderSkill(t, name+"/skill.md")
	}

	for _, name := range []string{"clue-delta/skill.md", "clue-extract/skill.md", "clue-verify/skill.md"} {
		content := rendered[name]
		for _, want := range []string{
			"authorization and protected-integration boundary",
			"not a demand for duplicate human code review",
			"only a human-controlled PR merge accepts it",
			"A PR alone displays hosted CI but does not enforce it",
			"branch protection makes its required status check a merge precondition",
			"commit every intended edit",
			"`git status --porcelain` to be empty",
			"head branch and SHA equal the current local branch and `HEAD`",
			"If either side differs",
			"Push is durability, never a signal",
			"Every working turn that changed anything ends by committing and pushing the change branch",
			"The PR exists from first publication and starts as a draft",
			"Marking the PR ready for review is the explicit act that claims a candidate",
			"A substantive edit to a ready PR returns it to draft",
			"Stopping anywhere else is ordinary, not an exception",
			"say where the work stands rather than that a ready PR exists",
			"Review fixes stay on the same branch and PR and repeat the complete updater handoff",
		} {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not contain review-handoff rule %q", name, want)
			}
		}
	}

	verify := rendered["clue-verify/skill.md"]
	for _, want := range []string{
		"Every intended edit, including each review fix, is committed and `git status --porcelain` is empty",
		"Every working turn on this change that changed anything ended by pushing the change branch",
		"the PR existed as a draft from first publication rather than appearing only at readiness",
		"When the PR is marked ready, the current branch is its head branch",
		"reported verification ran against that commit",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not contain hosted verification item %q", want)
		}
	}
}

func TestAC040_ReviewResultsAreDurableAndCommitBound(t *testing.T) {
	for _, name := range []string{"clue-delta", "clue-extract", "clue-verify"} {
		content := mustRenderSkill(t, name+"/skill.md")
		for _, want := range []string{
			"Every review of an existing hosted PR is bound to its observed head SHA",
			"A clean result applies only to that commit; every substantive edit invalidates it",
			"publish the finding there and leave it unresolved until a hosted commit contains the reviewed repair",
		} {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not carry durable review state %q", name, want)
			}
		}
	}
}

func TestAC040_ReviewWithoutResolvableHostStateFailsOpenly(t *testing.T) {
	verify := mustRenderSkill(t, "clue-verify/skill.md")
	for _, want := range []string{
		"If the reviewer cannot publish a resolvable finding, report the PR as not merge-ready",
		"never claim a chat-only finding has equivalent protection",
		"the isolated reviewer itself remains read-only",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not expose the unenforced-review fallback %q", want)
		}
	}
}

func TestAC132_UnitPositive_AnyEditorOwnsTheExactFastForwardHandoff(t *testing.T) {
	for _, name := range []string{"clue-delta", "clue-extract", "clue-verify"} {
		content := mustRenderSkill(t, name+"/skill.md")
		for _, want := range []string{
			"Any agent that edits an existing PR becomes the updater for that turn",
			"record its hosted head",
			"push only normal fast-forward updates, never force",
			"When accepted `main` advances while the change is open, incorporate it by merging `main` into the change branch with a normal push, never by rewriting hosted history",
			"Resolve satisfied review conversations only after the hosted head contains their reviewed repair",
			"the agent may review or help update an existing PR under the handoff above",
		} {
			if !strings.Contains(content, want) {
				t.Errorf("%s does not carry the exact updater handoff %q", name, want)
			}
		}
	}
}

// AC-041 was retired because its handoff put the push after the clean review
// and made incorporating `main` conditional on publication. Both readings are
// still natural to write, and nothing else would notice them returning: the
// generated skill directory is the whole carrier.
func TestAC132_UnitNegative_CarriersDoNotRestoreTheRetiredHandoffOrdering(t *testing.T) {
	for _, name := range skillNames {
		content := mustRenderSkill(t, name+"/skill.md")
		for _, retired := range []string{
			"obtains a clean review of the repaired commit, pushes without force",
			"before publishing, recheck that head",
			"After a PR is published, incorporate a newer accepted `main`",
			"Rebasing an unpublished local branch before its first publication remains allowed",
			"never as a draft",
			`local stopping point such as "commit only"`,
		} {
			if strings.Contains(content, retired) {
				t.Errorf("%s restored the retired handoff ordering %q", name, retired)
			}
		}
	}
}

func TestAC132_UnitPositive_ConcurrentOrClosedPRStateFailsSafely(t *testing.T) {
	verify := mustRenderSkill(t, "clue-verify/skill.md")
	for _, want := range []string{
		"If the head changed underneath the turn or a push is rejected as non-fast-forward",
		"fetch and reconcile without overwriting remote work",
		"A repair pushed to a ready PR returns it to draft until the repaired head has its own verification and clean review pass",
		"If the PR merged or closed, stop without pushing — the one case where a turn ends unpushed",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not fail safely on contested PR state %q", want)
		}
	}
}

func TestUnit_AgenticReviewLoopConvergesOnCurrentCommit(t *testing.T) {
	rendered := map[string]string{}
	for _, name := range skillNames {
		rendered[name+"/skill.md"] = mustRenderSkill(t, name+"/skill.md")
	}

	verify := rendered["clue-verify/skill.md"]
	for _, want := range []string{
		"never ask the human to clear context or initiate a separate review",
		"start a new read-only reviewer without the implementation conversation",
		"recover a full change's proposal from branch history",
		"label it `in-context fallback`",
		"only findings about correctness, intent mismatch, regressions, security, missing evidence, or unjustified complexity",
		"operative requirement or declared intent that is violated",
		"the concrete consequence",
		"Apply authoritative decisions and the repository's explicit lifecycle rules",
		"human-controlled merge does not require duplicate human code review",
		"lifecycle-successor evidence satisfies a requirement when the repository declares that transition",
		"lifecycle-correct state are not actionable defects by themselves",
		"a previous clean result applies only to the commit it reviewed",
		"Do not mark the PR ready with unresolved blocking findings or without such a pass",
		"Report the final review mode, reviewed commit, number of review passes run, and advisory findings left open",
		// The severity gate is what makes the loop terminate: an undifferentiated
		// list spends the same effort on a stale figure as on a corpus-breaking
		// defect, and repairing the figure used to restart the whole pass.
		"Every finding is classified **blocking** or **advisory**",
		"a blocking finding is actionable",
		"an advisory is a non-actionable observation for the readiness gate",
		"a reviewer brief cannot redefine the severity model",
		"A finding whose substance is a count, total, population figure, or arithmetic disagreement is **advisory** whatever the brief called it",
		"a wrong, missing, or reused identity remains **blocking**",
		"The reviewer spends no pass re-deriving figures",
		"An advisory repair may ride before a review pass already required by a blocking repair",
		"an advisory first reported by a pass with no blocking findings stays in the handoff for a later change",
		"advisories stay in the verification handoff rather than becoming repair-required conversations",
		"An advisory finding alone does not start another pass",
		"without changing its reviewed commit",
		"Scope that pass to the diff since the reviewed commit plus the carriers those files declare",
		// Reaching the maximum is a report and a question, never permission to
		// publish (PDR-036). The previous rule let a blocking finding on the last
		// pass earn the next one, which read as a bound and never was one. No
		// digit is pinned: the default is C-017's to state and an adopter's to
		// change, so raising it must not move a test.
		"passes** run for one change. That is a maximum, never a quota",
		"a further pass runs only when the immediately preceding pass returned at least one blocking finding",
		"When the maximum is reached and blocking findings remain, stop and report them to the human",
		"ask whether to run further passes; only that answer runs another",
		"a repository that wants a different maximum states it in its own `AGENTS.md` conventions",
		"whether or not the maximum was reached",
		"number of review passes run",
		"the bounded loop is over, and an edit would create a new candidate that this exact-commit rule requires reviewing",
		"Advisories do not become repair-required conversations",
	} {
		if !strings.Contains(verify, want) {
			t.Errorf("clue-verify/skill.md does not contain agentic-review rule %q", want)
		}
	}
	commitCandidate := strings.Index(verify, "commit the complete candidate")
	verifyCandidate := strings.Index(verify, "run the applicable local checks against that commit")
	reviewCandidate := strings.Index(verify, "start a new read-only reviewer")
	if commitCandidate < 0 || verifyCandidate <= commitCandidate || reviewCandidate <= verifyCandidate {
		t.Error("clue-verify must commit the candidate, verify that commit, then start agentic review")
	}

	for _, name := range []string{"clue-delta/skill.md", "clue-extract/skill.md"} {
		if !strings.Contains(rendered[name], "automatic agentic review loop") {
			t.Errorf("%s does not invoke the automatic agentic review loop", name)
		}
	}
}

func TestSanity_GeneratedSkillsCarryMergeHistoryBoundary(t *testing.T) {
	for _, name := range []string{"clue-delta", "clue-extract", "clue-verify"} {
		content := mustRenderSkill(t, name+"/skill.md")
		if !strings.Contains(content, "human accepts the ready pull request with a merge commit") {
			t.Errorf("%s does not state the supported merge-commit acceptance mode", name)
		}
		if !strings.Contains(content, "disable squash and rebase-and-merge") {
			t.Errorf("%s does not reject provenance-losing merge modes", name)
		}
	}
}

func mustRenderSkill(t *testing.T, relativePath string) string {
	t.Helper()
	directory := strings.TrimSuffix(relativePath, "/skill.md") + "/"
	var content strings.Builder
	for _, file := range mustRender(t) {
		if strings.HasPrefix(file.relativePath, directory) {
			content.Write(file.content)
			content.WriteByte('\n')
		}
	}
	if content.Len() == 0 {
		t.Fatalf("%s was not rendered", relativePath)
	}
	return content.String()

}

func mustRenderFile(t *testing.T, relativePath string) string {
	t.Helper()
	for _, file := range mustRender(t) {
		if file.relativePath == relativePath {
			return string(file.content)
		}
	}
	t.Fatalf("%s was not rendered", relativePath)
	return ""
}

func mustRender(t *testing.T) []renderedFile {
	t.Helper()
	files, err := render()
	if err != nil {
		t.Fatal(err)
	}
	return files
}
