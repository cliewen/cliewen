package corpus

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	jvmAnnotationRe     = regexp.MustCompile(`@(?:[A-Za-z_$][A-Za-z0-9_$]*\.)*([A-Za-z_$][A-Za-z0-9_$]*)`)
	jvmAnnotationCallRe = regexp.MustCompile(`@(?:[A-Za-z_$][A-Za-z0-9_$]*\.)*[A-Za-z_$][A-Za-z0-9_$]*\s*\([^)]*\)`)
	jvmLiteralTagRe     = regexp.MustCompile(`@(?:[A-Za-z_$][A-Za-z0-9_$]*\.)*Tag\(\s*"([^"]+)"\s*\)`)
	jvmTagCallRe        = regexp.MustCompile(`@(?:[A-Za-z_$][A-Za-z0-9_$]*\.)*Tag\s*\(`)
	jvmClassRe          = regexp.MustCompile(`\b(?:class|interface|enum|object|record)\s+[A-Za-z_$][A-Za-z0-9_$]*`)
	jvmKotlinFunRe      = regexp.MustCompile(`\bfun\s+(?:` + "`([^`]+)`" + `|([A-Za-z_$][A-Za-z0-9_$]*))\s*\(`)
	jvmJavaMethodRe     = regexp.MustCompile(`(?:^|\s)(?:(?:public|protected|private|static|final|synchronized|abstract|native|strictfp|default)\s+)*(?:<[^>{};]+>\s+)?(?:void|[A-Za-z_$][A-Za-z0-9_$]*(?:\.[A-Za-z_$][A-Za-z0-9_$]*)*)(?:\s*<[^;{}()]+>)?(?:\s*\[\])?\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*\([^;{}]*\)\s*(?:throws\s+[^{=]+)?\s*(?:\{|=)`)
	jvmNamedRe          = regexp.MustCompile(`^test([A-Z][A-Z0-9]*?)(\d+)_(Unit|Integration|E2E|Performance)(Positive|Negative)(?:_[A-Za-z0-9_$]+)+$`)
)

type jvmEvidence struct {
	subject   string
	ac        string
	testType  string
	direction string
}

type jvmAnnotationBlock struct {
	tags           []string
	hasExecutable  bool
	unsupportedTag bool
}

// harvestJVMEvidence recognizes only source shapes ADR-036 can attribute to
// one executable. It is deliberately a conservative source scanner, not a
// Java/Kotlin compiler or a test-runner discovery engine.
func harvestJVMEvidence(text, path string, prefixes map[string]bool) ([]jvmEvidence, []Issue) {
	lines := strings.Split(stripJvmCommentsAndTextBlocks(text), "\n")
	var evidence []jvmEvidence
	var issues []Issue
	var pending jvmAnnotationBlock
	inTripleString := false

	for lineNumber, line := range lines {
		if strings.Count(line, `"""`)%2 == 1 {
			inTripleString = !inTripleString
			continue
		}
		if inTripleString {
			continue
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		lineBlock := parseJvmAnnotations(line)
		if strings.HasPrefix(trimmed, "@") && len(jvmAnnotationRe.FindAllStringSubmatch(line, -1)) > 0 {
			pending.tags = append(pending.tags, lineBlock.tags...)
			pending.hasExecutable = pending.hasExecutable || lineBlock.hasExecutable
			pending.unsupportedTag = pending.unsupportedTag || lineBlock.unsupportedTag
		}

		if jvmClassRe.MatchString(line) {
			classACs, _, _ := classifyJvmTags(pending.tags, prefixes)
			for _, ac := range classACs {
				issues = append(issues, Issue{path, "class-level JVM tag \"" + strings.ReplaceAll(ac, "-", "_") + "\" cannot be attributed to one executable (ADR-036)"})
			}
			if pending.unsupportedTag {
				issues = append(issues, Issue{path, "unsupported JVM @Tag syntax near line " + strconv.Itoa(lineNumber+1) + " receives no evidence credit (ADR-036: use a one-line literal @Tag on the executable)"})
			}
			pending = jvmAnnotationBlock{}
			continue
		}

		methodName := jvmMethodName(line, strings.HasSuffix(path, ".kt"))
		if methodName != "" {
			blockEvidence, blockIssues := classifyJvmExecutable(methodName, pending, prefixes, path, lineNumber+1)
			evidence = append(evidence, blockEvidence...)
			issues = append(issues, blockIssues...)
			pending = jvmAnnotationBlock{}
			continue
		}

		// Annotation-only lines remain pending. Any other source construct
		// breaks attachment rather than letting metadata drift to a later
		// method.
		if strings.HasPrefix(trimmed, "@") {
			continue
		}
		if pending.hasExecutable || len(pending.tags) > 0 || pending.unsupportedTag {
			acs, _, _ := classifyJvmTags(pending.tags, prefixes)
			if pending.unsupportedTag {
				issues = append(issues, Issue{path, "unsupported JVM @Tag syntax near line " + strconv.Itoa(lineNumber+1) + " receives no evidence credit (ADR-036: use a one-line literal @Tag on the executable)"})
			} else if len(acs) > 0 {
				issues = append(issues, Issue{path, "JVM evidence annotations near line " + strconv.Itoa(lineNumber+1) + " cannot be attributed to a supported executable and receive no credit (ADR-036)"})
			}
			pending = jvmAnnotationBlock{}
		}
	}

	if len(pending.tags) > 0 || pending.unsupportedTag {
		acs, _, _ := classifyJvmTags(pending.tags, prefixes)
		if len(acs) > 0 || pending.unsupportedTag {
			issues = append(issues, Issue{path, "JVM evidence annotations at end of file cannot be attributed to a supported executable and receive no credit (ADR-036)"})
		}
	}
	return evidence, issues
}

func parseJvmAnnotations(line string) jvmAnnotationBlock {
	var block jvmAnnotationBlock
	for _, match := range jvmAnnotationRe.FindAllStringSubmatch(line, -1) {
		switch match[1] {
		case "Test", "ParameterizedTest", "RepeatedTest", "TestFactory", "TestTemplate":
			block.hasExecutable = true
		case "Tags":
			block.unsupportedTag = true
		}
	}
	for _, match := range jvmLiteralTagRe.FindAllStringSubmatch(line, -1) {
		block.tags = append(block.tags, match[1])
	}
	if len(jvmTagCallRe.FindAllStringIndex(line, -1)) != len(jvmLiteralTagRe.FindAllStringIndex(line, -1)) {
		block.unsupportedTag = true
	}
	return block
}

func classifyJvmExecutable(name string, block jvmAnnotationBlock, prefixes map[string]bool, path string, lineNumber int) ([]jvmEvidence, []Issue) {
	subject := "JVM executable " + name
	tagACs, tagTypes, tagDirections := classifyJvmTags(block.tags, prefixes)
	nameEvidence, named := namedJvmEvidence(name, subject, prefixes)
	var issues []Issue

	if block.unsupportedTag {
		issues = append(issues, Issue{path, subject + " uses unsupported JVM @Tag syntax near line " + strconv.Itoa(lineNumber) + " and receives no classified credit (ADR-036)"})
		return unclassifiedJvmEvidence(subject, tagACs), issues
	}

	hasTaggedEvidence := len(tagACs) > 0 || len(tagTypes) > 0 || len(tagDirections) > 0
	if named && hasTaggedEvidence {
		if !block.hasExecutable || len(tagACs) != 1 || len(tagTypes) != 1 || len(tagDirections) != 1 ||
			tagACs[0] != nameEvidence.ac || tagTypes[0] != nameEvidence.testType || tagDirections[0] != nameEvidence.direction {
			issues = append(issues, Issue{path, subject + " has disagreeing or incomplete named and JUnit evidence carriers and receives no classified credit (ADR-036)"})
			return unclassifiedJvmEvidence(subject, append(tagACs, nameEvidence.ac)), issues
		}
		return []jvmEvidence{nameEvidence}, nil
	}

	if named {
		return []jvmEvidence{nameEvidence}, nil
	}

	if len(tagACs) == 0 {
		return nil, nil
	}
	if !block.hasExecutable {
		issues = append(issues, Issue{path, subject + " carries JVM AC tags without a supported executable annotation and receives no credit (ADR-036)"})
		return nil, issues
	}
	if len(tagACs) > 1 || len(tagTypes) > 1 || len(tagDirections) > 1 {
		issues = append(issues, Issue{path, subject + " must carry at most one JVM AC identity, proof type, and direction; ambiguous metadata receives no classified credit (ADR-036)"})
		return unclassifiedJvmEvidence(subject, tagACs), issues
	}

	result := unclassifiedJvmEvidence(subject, tagACs)
	if len(tagTypes) == 1 && len(tagDirections) == 1 {
		result[0].testType = tagTypes[0]
		result[0].direction = tagDirections[0]
	}
	return result, nil
}

func classifyJvmTags(tags []string, prefixes map[string]bool) (acs, testTypes, directions []string) {
	acSet, typeSet, directionSet := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, tag := range tags {
		switch strings.ToLower(tag) {
		case "unit":
			typeSet["Unit"] = true
		case "integration":
			typeSet["Integration"] = true
		case "e2e":
			typeSet["E2E"] = true
		case "performance":
			typeSet["Performance"] = true
		case "positive", "negative":
			directionSet[strings.ToLower(tag)] = true
		default:
			norm := strings.ReplaceAll(tag, "_", "-")
			if match := jvmACRe.FindStringSubmatch(norm); match != nil && prefixes[match[1]] {
				acSet[norm] = true
			}
		}
	}
	for ac := range acSet {
		acs = append(acs, ac)
	}
	for testType := range typeSet {
		testTypes = append(testTypes, testType)
	}
	for direction := range directionSet {
		directions = append(directions, direction)
	}
	return acs, testTypes, directions
}

func namedJvmEvidence(name, subject string, prefixes map[string]bool) (jvmEvidence, bool) {
	match := jvmNamedRe.FindStringSubmatch(name)
	if match == nil || !prefixes[match[1]] {
		return jvmEvidence{}, false
	}
	return jvmEvidence{
		subject:   subject,
		ac:        match[1] + "-" + match[2],
		testType:  match[3],
		direction: strings.ToLower(match[4]),
	}, true
}

func unclassifiedJvmEvidence(subject string, acs []string) []jvmEvidence {
	seen := map[string]bool{}
	var result []jvmEvidence
	for _, ac := range acs {
		if seen[ac] {
			continue
		}
		seen[ac] = true
		result = append(result, jvmEvidence{subject: subject, ac: ac})
	}
	return result
}

func jvmMethodName(line string, kotlin bool) string {
	withoutAnnotations := jvmAnnotationCallRe.ReplaceAllString(line, "")
	withoutAnnotations = jvmAnnotationRe.ReplaceAllString(withoutAnnotations, "")
	withoutAnnotations = maskJvmQuotedLiterals(withoutAnnotations)
	if kotlin {
		if match := jvmKotlinFunRe.FindStringSubmatch(withoutAnnotations); match != nil {
			if match[1] != "" {
				return match[1]
			}
			return match[2]
		}
		return ""
	}
	for _, match := range jvmJavaMethodRe.FindAllStringSubmatch(withoutAnnotations, -1) {
		switch match[1] {
		case "if", "for", "while", "switch", "catch", "synchronized":
			continue
		default:
			return match[1]
		}
	}
	return ""
}

func maskJvmQuotedLiterals(line string) string {
	masked := []byte(line)
	inString, inChar, escaped := false, false, false
	for i, ch := range masked {
		if inString || inChar {
			masked[i] = ' '
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if inString && ch == '"' {
				inString = false
			} else if inChar && ch == '\'' {
				inChar = false
			}
			continue
		}
		if ch == '"' {
			masked[i] = ' '
			inString = true
		} else if ch == '\'' {
			masked[i] = ' '
			inChar = true
		}
	}
	return string(masked)
}

// stripJvmCommentsAndTextBlocks removes comments while preserving quoted
// annotation arguments. Triple-quoted Kotlin strings and Java text blocks are
// blanked because source-looking lines inside them are not annotations.
func stripJvmCommentsAndTextBlocks(source string) string {
	var out strings.Builder
	inBlockComment, inString, inChar, inTriple, escaped := false, false, false, false, false
	for i := 0; i < len(source); {
		if inBlockComment {
			if i+1 < len(source) && source[i] == '*' && source[i+1] == '/' {
				out.WriteString("  ")
				i += 2
				inBlockComment = false
				continue
			}
			if source[i] == '\n' {
				out.WriteByte('\n')
			} else {
				out.WriteByte(' ')
			}
			i++
			continue
		}
		if inTriple {
			if i+2 < len(source) && source[i:i+3] == `"""` {
				out.WriteString("   ")
				i += 3
				inTriple = false
				continue
			}
			if source[i] == '\n' {
				out.WriteByte('\n')
			} else {
				out.WriteByte(' ')
			}
			i++
			continue
		}
		if !inString && !inChar && i+1 < len(source) && source[i] == '/' && source[i+1] == '/' {
			for i < len(source) && source[i] != '\n' {
				out.WriteByte(' ')
				i++
			}
			continue
		}
		if !inString && !inChar && i+1 < len(source) && source[i] == '/' && source[i+1] == '*' {
			out.WriteString("  ")
			i += 2
			inBlockComment = true
			continue
		}
		if !inString && !inChar && i+2 < len(source) && source[i:i+3] == `"""` {
			out.WriteString("   ")
			i += 3
			inTriple = true
			continue
		}
		ch := source[i]
		out.WriteByte(ch)
		if inString || inChar {
			if escaped {
				escaped = false
			} else if ch == '\\' {
				escaped = true
			} else if inString && ch == '"' {
				inString = false
			} else if inChar && ch == '\'' {
				inChar = false
			}
		} else if ch == '"' {
			inString = true
		} else if ch == '\'' {
			inChar = true
		}
		i++
	}
	return out.String()
}
