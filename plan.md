# Delivery & DevOps Plan (plan.md)

Version: 0.1  
Date: 2025-10-05  
Aligns With: CRS v0.5 (Jane for Pool Professionals), ERS v0.5

> Objective: Deliver the *shortest path to stakeholder value* by completing a production-grade vertical slice (Health / Observability) first, then layering domain capabilities reusing the established DevOps/SRE patterns. Every later feature inherits: logging, metrics, tracing, SLO framing, CI/CD quality gates, security & audit patterns.

---
## 0. Guiding Principles
- Vertical Slice First: Each chapter ends in a deployable increment and a public blog post artifact.
- Sub-Chapter = Merge Unit: Completed when acceptance checklist passes (tests + quality gates) and merged to `develop`.
- Commit Granularity: Each bullet under a step represents a single logical commit (small, reviewable, revert-friendly).
- Branch Strategy:
  - `develop`: Integration branch (all sub-chapter merges land here first).
  - `staging`: Promotion branch for chapter-level release candidates (manual or automated promotion after Chapter QA). Mirrors production config except scale.
  - `production`: Only tagged, immutable release merges (`release-x.y.z`).
  - `feature/<chapter>-<subchapter>-<short-desc>`: Working branches per sub-chapter.
- Versioning: Semantic Versioning. First production deploy after Chapter 3 (core chemistry) = `v0.3.0` (Ch1 = 0.1.x, Ch2 = 0.2.x, etc.).
- Definition of Done (chapter): All mapped ERS IDs implemented & validated, blog post published, staging soak tests green.
- Naming Conventions: commits use prefixes (`feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`, `ops:`). Branch names kebab-case.

---
## 1. Chapter 1 – Observability & Operational Baseline (Health Endpoint Vertical Slice)
**Goal:** Implement ERS observability, security skeleton, build metadata, CI/CD gates using only the existing health endpoint as domain surface. Stakeholder value: demonstrates production-readiness culture early; enables future rapid diagnosis.

**Primary CRS/ERS Coverage:** CRS #15, #17 (partial), ERS IDs: E-ARCH-001..003, E-OBS-001..004, E-SEC-004, E-CICD-001..002, E-REL-001..002, partial E-AUD-001 (structure only).

### 1.1 Sub-Chapter: Repo & Metadata Foundations (Merge → develop)
Feature Branch: `feature/ch1-1-repo-metadata`

Commit Steps:
1. `chore: add CODEOWNERS and PR template`  
2. `chore: add commitlint / conventional commits config (optional)`  
3. `feat: embed build metadata (version, commit, build_date) in binary`  
4. `feat: extend /health to include version, commit, build_date, uptime`  
5. `test: add unit test for health handler metadata fields`  
6. `docs: update README with build metadata usage`  

### 1.2 Sub-Chapter: Structured Logging & Correlation
Branch: `feature/ch1-2-logging`

Commits:
1. `feat: enhance zap middleware with request_id generation`  
2. `feat: add trace_id placeholder field in logs (empty until tracing)`  
3. `chore: add log level via env (LOG_LEVEL)`  
4. `test: logging middleware unit tests (request id presence)`  
5. `docs: logging strategy section added to README`  

### 1.3 Sub-Chapter: Split Health Endpoints & Readiness Logic
Branch: `feature/ch1-3-health-split`

Commits:
1. `feat: add /health/live & /health/ready endpoints`  
2. `feat: readiness check DB placeholder (mocked / future real)`  
3. `feat: add dependency status array to readiness payload`  
4. `test: readiness endpoint returns 503 when dependency failure simulated`  
5. `docs: health endpoint contract in docs/`  

### 1.4 Sub-Chapter: Metrics (Prometheus)
Branch: `feature/ch1-4-metrics`

Commits:
1. `feat: add /metrics endpoint with Prometheus handler`  
2. `feat: add RED counters histograms for HTTP requests`  
3. `feat: add custom gauge uptime_seconds & health_checks_total`  
4. `test: metrics e2e scrape contains expected metric names`  
5. `docs: metrics catalog table`  

### 1.5 Sub-Chapter: Tracing (OpenTelemetry)
Branch: `feature/ch1-5-tracing`

Commits:
1. `feat: init OpenTelemetry tracer provider (OTLP optional)`  
2. `feat: gin tracing middleware integration`  
3. `feat: propagate trace_id into log context`  
4. `chore: local otel collector docker-compose snippet (optional)`  
5. `test: trace span creation integration test (assert headers)`  
6. `docs: tracing quickstart`  

### 1.6 Sub-Chapter: SLO & Alert Design Draft (No Infra Yet)
Branch: `feature/ch1-6-slo-drafts`

Commits:
1. `docs: add SLO.md (availability & latency objectives, error budget)`  
2. `docs: add ALERTS.md (fast/slow burn examples)`  
3. `feat: expose synthetic readiness probe script (scripts/synthetic-check.ps1)`  
4. `test: synthetic script exit codes`  

### 1.7 Sub-Chapter: CI/CD Pipeline Hardening
Branch: `feature/ch1-7-ci-cd`

Commits:
1. `feat: github action workflow for build test lint`  
2. `feat: add security scan step (gosec or trivy image scan)`  
3. `feat: add coverage report & status badge`  
4. `feat: add build artifact publishing (container image) with tags`  
5. `test: pipeline simulation doc`  
6. `docs: ci-cd architecture section`  

### 1.8 Sub-Chapter: Supply Chain & Image Integrity
Branch: `feature/ch1-8-supply-chain`

Commits:
1. `feat: generate SBOM (syft) in CI`  
2. `feat: trivy vulnerability scan gating`  
3. `feat: cosign keyless sign image (GitHub OIDC)`  
4. `docs: supply-chain.md (attestations & provenance)`  

### 1.9 Sub-Chapter: Staging Deploy Integration
Branch: `feature/ch1-9-staging-deploy`

Commits:
1. `feat: add k8s manifests (health-only slice) /k8s/base`  
2. `feat: add k8s overlay staging with resource requests & probes`  
3. `ops: add deployment job in CI for staging on tag v0.1.x`  
4. `test: helm template smoke test`  
5. `docs: deployment.md (promotion flow)`  

### Chapter 1 Promotion
- Merge all sub-chapter branches → `develop`.
- Run full CI; tag `v0.1.0` → auto deploy to `staging` (soak 24h synthetic checks).  
- Manual checklist (SLO baseline, metrics visible).  
- Optional production deploy of slice (not required yet).  

### Blog Post 1
Title: *“From Zero to Production-Ready Health Endpoint: Building an Observability-First Vertical Slice”*  
Outline: Problem (why start here) → Steps → Tooling choices → SLO draft → Lessons learned.

---
## 2. Chapter 2 – Security & RBAC Foundation
**Goal:** Implement AuthN/Z & audit skeleton so later domain entities inherit secure posture.

Coverage: CRS #1, #17; ERS E-SEC-001..005, E-AUD-001..002.

### Sub-Chapters & Commits
2.1 `feature/ch2-1-auth-model`  
- feat: user & role tables schema migration  
- feat: password hashing & user creation CLI  
- feat: login endpoint (JWT issuance)  
- feat: refresh token rotation  
- test: auth unit tests  
- docs: AUTH.md  

2.2 `feature/ch2-2-rbac-middleware`  
- feat: role-based authorization middleware  
- feat: add audit events for login success/fail  
- test: RBAC negative cases  
- docs: update SLO.md with auth failure metrics  

2.3 `feature/ch2-3-webhook-security`  
- feat: HMAC signature util  
- test: signature verification tests  
- docs: webhook security section  

2.4 `feature/ch2-4-audit-hardening`  
- feat: audit event persistence pattern (append-only txn)  
- chore: DB permission restrictions (no update/delete)  
- test: audit immutability tests  
- docs: AUDIT.md  

Promotion: Tag `v0.2.0` -> staging; smoke test user flows → optional production if needed.

### Blog Post 2
*“Securing the Backbone: Implementing Auth, RBAC & Immutable Audit Trails Early.”*

---
## 3. Chapter 3 – Core Chemistry & Dose Engine Slice
**Goal:** Deliver chemistry test input, dose recommendation, dose logging & visit report skeleton. First end-to-end business value.

Coverage: CRS #2, #5–#9, #6–#8 (chemistry & dosing), #16 (export partial), ERS E-DOM-001..005, E-PERF-001, E-OBS-004, E-INV-001.

### Sub-Chapters
3.1 Schema & Entities (`feature/ch3-1-schema`)
- feat: migrations for Pools, ServicePlans, Jobs, JobReadings, DoseEvents  
- test: schema migration tests  
- docs: ERD.svg  

3.2 Chemistry Input API (`feature/ch3-2-chem-input`)
- feat: POST /api/v1/jobs/{id}/readings  
- feat: validation errors & error envelope  
- test: readings validation  
- obs: add dosing_engine_latency metric placeholder  

3.3 Dose Engine v1 (`feature/ch3-3-dose-engine`)
- feat: engine formulas (baseline)  
- feat: recommended doses computation endpoint  
- test: dose calculation golden tests  
- perf: micro-bench for engine  
- docs: DOSE_ENGINE.md  

3.4 Dose Logging & Inventory (`feature/ch3-4-dose-logging`)
- feat: POST /api/v1/jobs/{id}/doses (actual)  
- feat: inventory decrement transactional  
- test: negative stock prevention  
- obs: increment dose_events_total metric  

3.5 Visit Report Assembly (`feature/ch3-5-visit-report`)
- feat: report aggregation service  
- feat: GET /api/v1/jobs/{id}/report  
- test: report snapshot test  
- docs: report contract  

3.6 Chemical Export v1 (`feature/ch3-6-export`)
- feat: CSV export endpoint (date range)  
- test: csv generation golden file  
- docs: compliance export format  

Promotion: Tag `v0.3.0` → staging → production (first true business value release).  

### Blog Post 3
*“Delivering Core Chemistry: From Readings to Actionable Dose Recommendations.”*

---
## 4. Chapter 4 – Scheduling & Technician Workflow
Goal: Add recurring service plans, manual routing, technician job UI endpoints.

Sub-Chapters: plan schema → job generation → route sequence → technician job lifecycle → performance tuning.
Tag: `v0.4.0` after completion.

### Blog Post 4
*“Optimizing the Day: Scheduling and Route Foundations.”*

---
## 5. Chapter 5 – Billing & Payments
Goal: Invoicing & Stripe integration.
Tag: `v0.5.0` (production). Blog 5.

---
## 6. Chapter 6 – Alerts & Threshold Monitoring
Goal: Out-of-band readings alert pipeline, acknowledgment workflow, metrics.
Tag: `v0.6.0`. Blog 6.

---
## 7. Chapter 7 – Inventory Refinement & Compliance Enhancements
Add override reasons, lot placeholder, extended export.
Tag: `v0.7.0`. Blog 7.

---
## 8. Chapter 8 – Integrations & Webhooks
Finalize events + delivery reliability dashboards.
Tag: `v0.8.0`. Blog 8.

---
## 9. Chapter 9 – Reliability Hardening & SLO Enforcement
Chaos tests, retry strategies, autoscaling policies.
Tag: `v0.9.0`. Blog 9.

---
## 10. Chapter 10 – Pre-MVP Launch Wrap & MVP Declaration
Audit gap review, documentation freeze, security scan sign-off. Tag `v1.0.0`.
Blog 10: *“Declaring MVP: Lessons from an Observability-First Journey.”*

---
## 11. Stakeholder Value Map
| Stakeholder | First Visible Value | Chapter |
|-------------|--------------------|---------|
| SRE/DevOps | Full observability slice | 1 |
| Security/Compliance | Auth & audit trail | 2 |
| Operations/Dispatch | Scheduling & routing basics | 4 |
| Technician | Mobile workflow & dosing | 3 & 4 |
| Finance | Invoices & payments | 5 |
| Compliance Officer | Export & audit | 3 & 7 |
| Customer | Visit report & portal | 3 (report), 4 (scheduling reliability) |
| Integration Partner | Webhooks | 8 |

---
## 12. Promotion Workflow
1. Feature Branch → PR → CI gates (lint, tests, coverage ≥ threshold, vuln scan) → merge to `develop`.
2. Chapter Completion: Tag pre-release (e.g., `v0.3.0-rc1`) → auto deploy to `staging`.
3. Staging Validation Checklist: SLO burn check, error rate < threshold, synthetic health success 24h, no P1 vulnerabilities.
4. Production Promotion: Create release tag (e.g., `v0.3.0`) → signed image pushed → `production` deployment.
5. Post-Deploy: Capture metrics & annotate dashboards.

---
## 13. Quality Gates & Automation Matrix
| Gate | Enforced At | Tools | Block? |
|------|-------------|-------|--------|
| Unit Tests | PR | `go test` | Yes |
| Coverage ≥ 70% core packages | PR | `go test -cover` | Yes |
| Lint | PR | `golangci-lint` | Yes |
| Vuln Scan (High+) | Merge & Release | `trivy`, `govulncheck` | Yes |
| SBOM Generation | Release | `syft` | No (artifact) |
| Image Sign | Release | `cosign` | Yes (if fail) |
| Synthetic Health | Staging soak | custom script | Yes (promotion) |
| SLO Regression (latency/error) | Pre-prod | Prometheus query | Manual gate |

---
## 14. ERS Mapping First Implementation Points
| ERS ID | Implemented Earliest In | Notes |
|--------|-------------------------|-------|
| E-ARCH-001 | Ch1.1 | Docker build metadata |
| E-OBS-001..004 | Ch1.4/1.5 | Logging, metrics, tracing |
| E-SEC-004 | Ch1.3 | Health only public |
| E-SEC-001..003 | Ch2 | Auth & audit |
| E-PERF-001 | Ch3.3 | Dose engine benchmarks |
| E-DOM-001..005 | Ch3.1 | Core schema |
| E-INV-001 | Ch3.4 | Decrement logic |
| E-AUD-001..002 | Ch2.4 | Audit immutability |
| E-API-001..002 | Ch8 | Webhooks |
| E-CICD-001..002 | Ch1.7 | CI/CD pipeline |

---
## 15. Risk Register (Active Monitoring)
| Risk | Phase Likely | Mitigation | Owner |
|------|--------------|-----------|-------|
| Scope creep in chemistry formulas | Ch3 | Lock v1 spec, version engine | Eng Lead |
| Latency > target for dose engine | Ch3 | Profiling early, caching constants | Eng |
| Security gaps (JWT misuse) | Ch2 | Pen test lite, automated tests | Sec Champion |
| Observability cost explosion | Ch1+ | Log sampling & metric cardinality review | SRE |
| Staging drift vs prod | All | GitOps manifests single source | DevOps |

---
## 16. Blog Series Index
| Blog # | Chapter | Working Title | Audience |
|--------|---------|---------------|----------|
| 1 | 1 | From Zero to Production-Ready Health Endpoint | DevOps/SRE/Founders |
| 2 | 2 | Securing the Backbone Early | Eng/Security |
| 3 | 3 | Chemistry Intelligence: Building a Dose Engine | Domain + Eng |
| 4 | 4 | Scheduling & Route Foundations | Ops/Eng |
| 5 | 5 | Monetizing the Flow: Invoices & Payments | Finance/Eng |
| 6 | 6 | Turning Readings into Actionable Alerts | Ops/SRE |
| 7 | 7 | Compliance & Inventory Integrity | Compliance/Exec |
| 8 | 8 | Event-Driven Expansion with Webhooks | Integrators |
| 9 | 9 | Hardening for Reliability & Scale | SRE/Eng |
| 10 | 10 | Declaring MVP & Next Horizons | Broad |

---
## 17. Execution Checklist Template (Per Sub-Chapter)
```
[ ] Feature branch created
[ ] ERS IDs referenced in PR description
[ ] Unit tests added/updated
[ ] Observability instrumentation added (log/metric/trace)
[ ] Security/RBAC verified (if applicable)
[ ] Docs updated (README or dedicated *.md)
[ ] CI pipeline green
[ ] Reviewer approvals (>=2)
```

---
## 18. Next Immediate Actions (You Can Start Now)
1. Create feature branch `feature/ch1-1-repo-metadata`.
2. Implement build metadata & enriched health response.
3. Add unit test for health handler.
4. Draft Blog #1 outline (skeleton file `blog/01-health-slice.md`).
5. Open tracking issue linking ERS IDs for Chapter 1.

---
## 19. Future Enhancements (Post MVP Fast Follows)
- Offline sync conflict resolution strategy (OQ-3) design doc.
- Route optimization heuristic (savings algorithm prototype) in MVP+.
- IoT ingestion pipeline (buffered writes, anomaly pre-processor) after stable chemistry logs.
- Feature flag service (engine versioning) for safe formula updates.

---
**END OF PLAN**
