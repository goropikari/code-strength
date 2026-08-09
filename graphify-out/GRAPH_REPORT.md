# Graph Report - . (2026-08-10)

## Corpus Check

- Corpus is ~6,133 words - fits in a single context window. You may not need a graph.

## Summary

- 111 nodes · 163 edges · 13 communities (7 shown, 6 thin omitted)
- Extraction: 85% EXTRACTED · 15% INFERRED · 0% AMBIGUOUS · INFERRED: 25 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)

- CLI Commands
- Interactive Selection
- Testing and CI
- AI Requirements
- Selector Tests
- Directory Scanner
- Generate Flags
- Application Entry
- Quality Governance
- Selector Output Tests
- Agent Configuration
- Excluded Paths
- Project Identity

## God Nodes (most connected - your core abstractions)

1. `Select()` - 11 edges
2. `selectLines()` - 9 edges
3. `selectInteractive()` - 8 edges
4. `render()` - 8 edges
5. `newGenerateCommand()` - 7 edges
6. `LevelForPath()` - 6 edges
7. `Testing Guidelines` - 6 edges
8. `newLevelCommand()` - 5 edges
9. `List()` - 5 edges
10. `handleKey()` - 5 edges

## Surprising Connections (you probably didn't know these)

- `development requirement level` --semantically_similar_to--> `Production implementation practices` [INFERRED] [semantically similar]
  docs/prd.md → plugins/code-strength-ai/skills/code-strength-ai/SKILL.md
- `newGenerateCommand()` --calls--> `List()` [INFERRED]
  cmd/generate.go → internal/scanner/scanner.go
- `newGenerateCommand()` --calls--> `Select()` [INFERRED]
  cmd/generate.go → internal/selector/selector.go
- `go test ./...` --conceptually_related_to--> `Testing Guidelines` [INFERRED]
  .github/workflows/ci.yml → TESTING.md
- `newGenerateCommand()` --calls--> `BuildEntries()` [INFERRED]
  cmd/generate.go → internal/generator/generator.go

## Import Cycles

- None detected.

## Hyperedges (group relationships)

- **Repository quality checks** — github_workflows_ci_go_test, github_workflows_ci_make_lint, github_workflows_ci_make_fmt, quality_readme_coding_quality_gate [INFERRED 0.85]
- **AI requirement level workflow** — ai_requirements_directory_levels, readme_code_strength_cli, plugins_code_strength_ai_skills_code_strength_ai_skill_code_strength_ai_skill, docs_prd_ai_requirement_level_guide [INFERRED 0.95]
- **Testing quality principles** — agents_testing_conventions, testing_testing_guidelines, testing_testify, testing_aaa_pattern, testing_race_detection [EXTRACTED 1.00]

## Communities (13 total, 6 thin omitted)

### Community 0 - "CLI Commands"

Cohesion: 0.15
Nodes (18): Command, newGenerateCommand(), Command, newLevelCommand(), init(), Definition, Directory, BuildEntries() (+10 more)

### Community 1 - "Interactive Selection"

Cohesion: 0.23
Nodes (20): File, handleKey(), isPathOrChild(), mark(), matchingDirectories(), printMatches(), readEscapeKey(), render() (+12 more)

### Community 2 - "Testing and CI"

Cohesion: 0.14
Nodes (14): go test ./..., make fmt, make lint, Repository testing conventions, CI workflow, go test ./..., make fmt, make lint (+6 more)

### Community 3 - "AI Requirements"

Cohesion: 0.16
Nodes (14): Directory-specific AI requirement levels, Directory-specific AI requirement level guide, Go/Cobra CLI, development requirement level, production requirement level, Recursive dependency traversal from AI starting points, code-strength-ai skill, code-strength level query (+6 more)

### Community 4 - "Selector Tests"

Cohesion: 0.29
Nodes (9): T, TestShouldSelectInteractive(), Select(), shouldSelectInteractive(), T, TestSelectIgnoresInvalidIndexes(), TestSelectParentSelectsDescendants(), TestSelectSearchesAndTogglesMultipleDirectories() (+1 more)

### Community 5 - "Directory Scanner"

Cohesion: 0.36
Nodes (7): List(), matchesExclude(), countString(), T, TestListFollowsDirectorySymlinks(), TestListIncludesRootAndSkipsDefaultsAndExtras(), walkDirectories()

### Community 8 - "Quality Governance"

Cohesion: 0.50
Nodes (4): Versioned quality contract, coding-quality-gate, Repository quality state, Latest quality gate evidence

## Knowledge Gaps

- **21 isolated node(s):** `github.com/goropikari/code-strength`, `Directory-specific AI requirement levels`, `Excluded paths`, `make fmt`, `golangci-lint configuration` (+16 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions

_Questions this graph is uniquely positioned to answer:_

- **Why does `newGenerateCommand()` connect `CLI Commands` to `Selector Tests`, `Directory Scanner`, `Generate Flags`?**
  _High betweenness centrality (0.305) - this node is a cross-community bridge._
- **Why does `Select()` connect `Selector Tests` to `CLI Commands`, `Interactive Selection`?**
  _High betweenness centrality (0.260) - this node is a cross-community bridge._
- **Why does `List()` connect `Directory Scanner` to `CLI Commands`?**
  _High betweenness centrality (0.091) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `Select()` (e.g. with `newGenerateCommand()` and `TestSelectIgnoresInvalidIndexes()`) actually correct?**
  _`Select()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `newGenerateCommand()` (e.g. with `BuildEntries()` and `Write()`) actually correct?**
  _`newGenerateCommand()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/goropikari/code-strength`, `Directory-specific AI requirement levels`, `Excluded paths` to the rest of the system?**
  _21 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `CLI Commands` be split into smaller, more focused modules?**
  _Cohesion score 0.1471861471861472 - nodes in this community are weakly interconnected._
