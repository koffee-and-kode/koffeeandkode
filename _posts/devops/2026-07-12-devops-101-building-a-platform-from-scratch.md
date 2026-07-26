---
title: "DevOps 101: Building a Production-Ready Platform from Scratch"
date: 2026-07-12
tags: [devops, ci-cd, docker, kubernetes, gitops, terraform, observability, platform-engineering]
excerpt: "DevOps maturity isn't measured by how many tools you run. Follow a five-engineer startup from a single deploy script to a GitOps-driven Kubernetes platform, stage by stage — the pressure that forces each jump, the tool that answers it, the cost it adds, and what to deliberately not build yet."
read_time: 22
---

Most DevOps content is a glossary. Here is Kubernetes, here is Terraform, here is Argo CD — useful as a reference, useless as a sequence. Nobody hands a five-person startup a Kubernetes cluster on day one, and yet plenty of engineers read "Kubernetes basics" before they've ever felt the problem it solves. The result is a familiar failure mode: teams adopt the tool before the pressure that justifies it, and end up operating something more complex than their actual problem.

This post is the sequence instead of the glossary. It follows one startup — five engineers, a backend service and a frontend, no dedicated ops hire — from a single deploy script to a GitOps-driven Kubernetes platform with real observability. At every stage: what's breaking, why the current setup can't absorb it, what capability has to come next, which tool provides it, what that tool costs you in complexity, and — just as important — what you should deliberately not build yet.

The team in this post is building a B2B analytics product: a Go backend, a frontend, a Postgres database, and a handful of background jobs that sync customer data overnight. Nothing about the shape of the product matters much; swap in a marketplace, a SaaS tool, whatever you're building. What matters is the sequence of pressure.

One thesis carries the whole post, and it's worth stating up front so the rest reads as evidence for it rather than a tool tour:

> DevOps maturity is not measured by how many tools a company uses. It is measured by how safely, quickly, and predictably the company can deliver software.

Kubernetes, Terraform, and Argo CD all show up in this post. None of them show up on day one, and for a lot of startups, none of them show up at all in year one. That's not a criticism of any of these tools — it's the actual order operations happen in, and the actual bar each one has to clear before it's worth its cost.

## Stage 0 — five engineers, one deploy script

The team ships a Go API, a Postgres database, and a frontend. Deployment is `git push` to a small VM, or `git push heroku main` if they're using a PaaS. There's no staging environment. Code review is "hey, can you glance at this" in Slack before merging to `main`.

This is correct. Not a compromise, not something to feel behind on — correct. Every hour spent on infrastructure at this stage is an hour not spent finding out if anyone wants the product. A staging environment, a CI pipeline, and a Terraform module are all solving problems this team doesn't have yet: nobody's shipped a bad deploy that mattered, nobody's had two people's changes collide, and the "infrastructure" is small enough to hold in one person's head.

**What's still fine here:** direct pushes to main, manual deploys, a single environment, a managed database, no IaC, no dedicated monitoring beyond default hosting-platform logs.

**What starts the next stage:** a second engineer joins, or the first one ships a bug that breaks something a design partner was using. The instant two people are touching the same codebase without a review step, or a broken `main` costs more than the five minutes it takes to fix it, the pressure has arrived.

## Stage 1 — repo governance: the platform starts at GitHub

Two things happen close together as the team grows past one or two people: someone's change breaks someone else's work-in-progress, and someone ships a bug that a thirty-second review would have caught. Both are symptoms of the same gap — there's no gate between "I wrote this" and "this is in main."

GitHub is the first piece of the delivery platform, and most of what it offers here is nearly free:

- **Pull requests and code review** as the default path to `main` — not because five engineers can't trust each other, but because a second pair of eyes catches the class of bug that "I already read my own diff" never will.
- **Branch protection** — require at least one approval, require the branch to be up to date before merge. Ten minutes of setup, and it closes the "I force-pushed over your work" failure mode for good.
- **Required status checks** — once CI exists (next section), branch protection can require it to pass before merge is even offered. This is what turns CI from "a thing that runs" into "a thing that matters."
- **CODEOWNERS** — at five engineers, overkill. Once the team splits into "backend people" and "frontend people," it routes review to whoever actually knows the code, instead of whoever happened to be online.
- **Release tags and semantic versioning** — worth adopting the moment anything outside the repo depends on a specific version: a mobile client pinned to an API version, a customer running a self-hosted agent. Skip it entirely if nothing external depends on your version number yet.
- **Dependabot (or equivalent)** — cheap, automatic, and it converts "we're three major versions behind and it's terrifying to upgrade" into a steady trickle of small, reviewable bumps.
- **Secret scanning** — free on GitHub, catches the `AWS_SECRET_ACCESS_KEY` someone pastes into a commit at 11pm. Turn it on and forget about it.

**What to skip for now:** a formal branching model beyond "short-lived branches off `main`." Trunk-based development is the right default for a SaaS product shipping continuously — GitFlow's release branches solve a problem (multiple supported versions in the wild) that a hosted product doesn't have.

## Stage 2 — CI: fast feedback before it's a rule, not after

Code review catches logic mistakes; it doesn't catch "this doesn't compile" or "this breaks the test suite," and human reviewers are slow to notice either. Continuous integration is the automated version of "does this actually work," running on every push instead of relying on a human to remember to run tests locally.

A CI pipeline is a graph of steps that runs on every commit: lint, test, build, and — once there's a deployable artifact — package and publish. The one rule that matters more than any tool choice: **fast feedback wins.** A ten-minute pipeline that people wait for is worth more than a ninety-minute one people learn to ignore.

Here's a realistic GitHub Actions workflow for the Go backend — lint, test, build a Docker image, and push it to a registry, tagged by both a mutable branch tag and an immutable commit digest:

{% raw %}
```yaml
name: ci

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.27"
          cache: true
      - name: lint
        run: go vet ./...
      - name: test
        run: go test ./... -race -cover

  build-and-push:
    needs: test
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            ghcr.io/acme/api:latest
            ghcr.io/acme/api:${{ github.sha }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```
{% endraw %}

A few things worth calling out that aren't obvious from the YAML itself:

- **`test` and `build-and-push` are separate jobs.** The build only runs after tests pass, and only on `main` — a PR from a fork gets the fast feedback without ever getting registry credentials.
- **`cache: true` on setup-go and `cache-from`/`cache-to` on the Docker build** turn a two-minute dependency-download tax into a ten-second cache hit. This is usually the single biggest lever on pipeline speed.
- **The image is tagged twice** — `latest` for humans skimming the registry, `github.sha` for anything that needs to reference an exact, immutable build. More on why that distinction matters in the next section.
- **Secrets are scoped to the job that needs them.** `GITHUB_TOKEN` here is enough for GHCR; a deploy step needing cloud credentials would get its own environment-scoped secret, not a blanket one available to every job.

**Where CI ends:** the moment there's a pushed image, CI's job is done. Whether that image reaches production is a separate concern — continuous *delivery* — and conflating the two is how you end up with a "CI" pipeline that's secretly also your deploy mechanism, with no gate in between. Keep them separate even before you need the gate; splitting them later, after deploy logic has grown into the test job, is a bigger refactor than it sounds.

## Stage 2, continued — Docker: a consistent artifact

Before CI can publish something, there has to be something worth publishing. A Docker image solves a specific problem: "works on my machine" becomes "works in this container," full stop, no dependency-version archaeology on the host.

Here's a small production-oriented Dockerfile for the Go backend — multi-stage, so the final image ships a binary and nothing else:

```dockerfile
# --- build stage ---
FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

# --- runtime stage ---
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/api /api

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/api"]
```

The build stage has the Go toolchain and every dependency; the runtime stage has one static binary and nothing else — no shell, no package manager, no attack surface beyond the binary itself. `USER nonroot` means a container breakout doesn't hand the attacker root. Copying `go.mod`/`go.sum` before the rest of the source means the dependency-download layer only invalidates when dependencies actually change, not on every code edit — the same caching principle as the CI YAML above, applied to the image build itself.

**Tags versus digests:** a tag like `api:latest` is a mutable pointer — it can point at a different image tomorrow. A digest (`api@sha256:abc123...`) is the image, immutably. The failure mode this prevents: staging and production both deploy `api:latest`, but staging pulled it Tuesday and production pulled it Thursday, and now two environments running "the same" tag are running different code. Deploy manifests should reference a digest (or at minimum a commit-SHA tag, as above), not a floating tag — the whole point of a build artifact is that it doesn't change out from under you.

**What to skip for now:** image signing (Cosign), SBOM generation, a private registry beyond GHCR/ECR/GAR's default. All real, all worth having eventually — none of them are what's breaking today.

## Stage 3 — cloud infrastructure: what to provision, and what to skip

Somewhere around the first handful of real customers, "one VM" stops being enough — not because of load, but because of blast radius. A bad deploy or a database hiccup now affects people who are paying, not just the team.

This is the point to be deliberate about infrastructure, and it's worth being explicit that **none of it requires Kubernetes yet.** The building blocks, kept vendor-neutral (AWS naming first, with GCP/Azure equivalents):

- **Compute** — a managed container runtime (AWS App Runner or ECS Fargate; GCP Cloud Run; Azure Container Apps) instead of a hand-managed VM. You get restarts, basic scaling, and a health check for free, without managing an OS.
- **Managed database** — RDS/Aurora, Cloud SQL, or Azure Database for Postgres. Database choices are sticky; pick boring and move on. Don't self-host Postgres on a VM you also deploy code to.
- **Object storage** — S3, GCS, or Azure Blob Storage for anything that isn't a database row: uploads, generated files, logs.
- **DNS and TLS** — Route 53 / Cloud DNS / Azure DNS, plus a managed certificate (ACM, Google-managed certs, or Azure's equivalent). Both are close to free and there's no reason to hand-roll either.
- **Load balancer** — sits in front of compute the moment there's more than one instance, or the moment you want zero-downtime deploys (new instances up, health-checked, then old ones drained).
- **Secrets manager** — Secrets Manager / Secret Manager / Key Vault. The alternative — API keys in `.env` files copied between machines — is the single most common way secrets end up in a Git history.
- **IAM** — least-privilege roles per service, not one shared key with admin access. This is cheap to get right early and expensive to unwind later, because by the time it hurts, a dozen things depend on the overprivileged key.
- **Backups** — automated snapshots on the managed database, and — critically — a restore that's actually been tested once. An untested backup is a hypothesis, not a backup.

**The Kubernetes question, answered early and directly:** Kubernetes is not the default starting point. It solves problems — packing many services efficiently onto shared compute, a consistent deployment API across dozens of services, self-healing at a scale where manual intervention doesn't scale — that a five-to-fifteen-engineer team with one or two services usually doesn't have yet. A managed container runtime gives you most of the operational benefit (restarts, rolling deploys, health checks) with a fraction of the operational surface. Section 7 covers exactly when that trade flips.

If you want the AWS-specific version of this progression — which services, in what order, as an AI startup's user count climbs from ten to a million — that's its own post: [0 to 1 Million with AWS]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}). This section is the vendor-neutral shape; that one is the vendor-specific depth.

## Stage 3, continued — Terraform: making infrastructure reviewable

Once there's a VPC, a database, a load balancer, and IAM roles, "I clicked through the console to set this up" stops being sustainable — nobody remembers exactly what's provisioned, and there's no diff to review before a change goes live. Infrastructure as Code fixes the actual problem: infrastructure changes become pull requests, with a plan you can read before it's applied.

Terraform's core ideas, kept brief because the depth belongs in its own post:

- **State** is the source of truth Terraform diffs against. Losing it, or two people applying against stale state at once, is the classic Terraform incident — hence a remote backend (S3 + a DynamoDB lock, or Terraform Cloud) from day one, never local state.
- **Modules** are reusable building blocks — a `database` module, a `service` module — instead of one 2,000-line file. Structure around logical groupings, not one mega-module with thirty inputs.
- **Plan, then apply** — `terraform plan` shows the diff before anything changes, and that diff is what gets reviewed in the pull request, the same way application code does.

A simplified layout that separates environments cleanly:

```text
infrastructure/
  modules/
    database/
    service/
    networking/
  environments/
    development/
    staging/
    production/
```

Each environment directory has its own state file and its own variable values, pointing at the same modules. This is what prevents "I tested this in staging" from meaning something different than "this is what's running in production" — they're built from identical modules, just different inputs.

**What never goes in Terraform state or Git:** actual secret values. Reference a secret by ARN/name and let the secrets manager hold the value — Terraform state is often stored in plaintext, and a leaked state file with embedded secrets is a full breach, not a near-miss.

## Stage 3, continued — from manual deploys to CD

Up to this point, "deploy" has meant a person running a command. The next failure mode is human error at the worst possible time — someone fat-fingers a deploy during a customer's business hours, or forgets a step that only exists in their memory.

Worth being precise about three different things people call "CI/CD":

- **Continuous integration** — every commit is tested and merged safely (Stage 2).
- **Continuous delivery** — every commit that passes CI is *deployable* — a human still decides when to release it.
- **Continuous deployment** — every commit that passes CI is *actually deployed*, automatically.

Most startups want continuous delivery with a fast, low-friction release step, not full continuous deployment — the difference between "deploy" and "release" earns its keep here. You can deploy new code to production *dark* (running, but not receiving traffic or not gated behind a flag) and release it separately by flipping a flag. That decouples "is this code safely running" from "should users see this feature," which is a much better question to answer under pressure than "should we deploy right now."

A GitHub Actions job that deploys directly (pushing a new image and update to the managed compute service) is the right level of ceremony for one or two services. Compare that to the GitOps model in Section 8 — direct-deploy-from-CI is simpler and has fewer moving parts; it starts to strain once there's more than a couple of services, because now N pipelines all independently know how to deploy, and there's no single place to see "what's actually running in prod right now" other than checking each one.

**Deployment strategy, chosen by blast radius, not by what's fashionable:** a **rolling deployment** (new instances up and healthy before old ones are drained) is the correct default for a stateless web service and is what your managed compute runtime already does for you. **Blue-green** (a full duplicate environment, instant cutover, instant rollback) earns its cost when a deploy is high-stakes and you want zero risk during the swap — expensive to run two full environments, so save it for when downtime is truly not an option. **Canary** (a small slice of traffic first, expand if healthy) is the right tool once you have enough traffic and enough monitoring to actually tell if the canary is healthy before it matters — canary without observability is just a slower, more confusing rolling deploy. At this stage, rolling is enough; canary becomes worth the setup around Stage 4, once there's real observability behind it.

**What "production-ready" concretely means**, beyond "it deploys automatically" — this is the checklist worth having *before* the first real customer, not after the first incident:

- **Health checks** the load balancer actually uses to route traffic, not just a `200 OK` on `/`.
- **Graceful shutdown** — handle `SIGTERM`, finish in-flight requests, then exit. Without this, every rolling deploy drops a percentage of requests.
- **Timeouts** on every outbound call — a downstream service hanging forever becomes your service hanging forever.
- **Retries with backoff and jitter** — a bare retry loop against a struggling downstream service is how you turn a blip into an outage; jitter prevents every client from retrying in lockstep and creating a thundering herd against the thing that's already struggling.
- **Idempotency** for anything that might get retried — a payment or a sync job that runs twice by accident is a data-integrity bug, not an edge case.
- **Database migrations** that ship as a controlled step, not "run this SQL by hand and hope."

Circuit breakers, formal disaster-recovery drills, and rate limiting matter — but they earn their place a bit later, once there's enough traffic and enough downstream dependencies for a struggling dependency to actually take you down with it.

## Stage 4 — first real growth: why Kubernetes, specifically, now

The team's grown past the backend and frontend into three or four services — maybe a sync worker split out, maybe a separate notifications service. Each one currently means another App Runner/Fargate service, another set of alarms, another deploy pipeline to maintain by hand. The operational overhead is starting to scale with the number of services, not with the amount of traffic.

This is specifically the pressure Kubernetes answers: a single, consistent API for running many services, instead of N slightly different hand-rolled setups. It is not "the industry standard, so you should use it" — plenty of companies stay on managed container runtimes well past this point, and that's a legitimate choice, not a lesser one. The trade only makes sense once the *number of services* — not the traffic to any one of them — is the thing generating operational cost.

Kept deliberately at the level of "what problem does each primitive solve," not a command reference:

- **Cluster and nodes** — the pool of machines Kubernetes schedules work onto. You mostly stop thinking about individual nodes once this is set up correctly.
- **Pod** — one or more co-scheduled containers; almost always one app container, sometimes a sidecar.
- **Deployment** — the declarative wrapper that manages rolling out Pods and keeping the right number running.
- **Service** — a stable network identity for a set of Pods, so nothing needs to track individual Pod IPs that come and go.
- **Ingress** — the L7 entry point that routes external HTTP traffic to the right Service.
- **ConfigMap / Secret** — configuration and sensitive values, injected as env vars or mounted files, kept out of the container image itself.
- **Namespace** — a logical partition inside one cluster, most often one per environment or one per team.
- **Horizontal Pod Autoscaler** — adds/removes Pods based on load, so capacity roughly tracks demand instead of sitting fixed.
- **Resource requests and limits** — what a Pod is guaranteed and what it's capped at; get this wrong and one noisy service starves its neighbors on the same node.
- **Probes** — readiness (ready for traffic?), liveness (still alive?), startup (give slow-starting apps more runway before liveness kicks in). A readiness probe that's too generous is a classic way to let a broken Pod take real traffic.
- **Jobs and CronJobs** — for the nightly sync job that doesn't fit the "long-running service" shape.

**What Kubernetes costs you that a managed runtime didn't:** you're now operating a scheduler, a networking layer, and an API server, even if a cloud provider manages the control plane for you. Nodes need patching (or a node auto-upgrade policy), YAML sprawl needs a templating or overlay strategy (Helm or Kustomize — its own post, 112), and the failure modes get more interesting: a misconfigured resource limit or an over-eager readiness probe can now take down a service in ways a simpler platform never would have let happen. This is the point where the operational cost is real, not hypothetical — worth crossing only once the "many services, one deploy API" problem is actually the one you have.

**Managing environments and clusters as this grows** is its own trade-off ladder, not a single correct architecture:

- **One cluster, namespaces per environment** — the simplest option, and the right default when starting out with Kubernetes. Lower cost, lower operational overhead, but a bad actor or a runaway workload in one namespace can, if RBAC and resource limits aren't tight, spill into another.
- **Separate production and non-production clusters** — the natural next step. A staging mistake can no longer take down production infrastructure. Cost roughly doubles for the cluster itself, but the blast-radius isolation is usually worth it well before you have a compliance reason to require it.
- **Separate clusters per region** — driven by latency or by data-residency requirements (GDPR-style rules), not by scale alone. Real operational overhead: now you're running N clusters that all need to look the same.
- **Separate clusters per business unit or compliance boundary** — usually a security or regulatory requirement (a HIPAA-scoped workload, an enterprise customer that requires dedicated infrastructure), not a default architecture choice.

The sensible progression is exactly the order above: namespaces first, split prod out next when the blast radius stops being acceptable, region/BU splits only when a real forcing function — latency, compliance, an actual customer requirement — shows up. Skipping straight to "one cluster per region" for a fifteen-person team is the same mistake as adopting Kubernetes too early, one level up.

## Stage 4, continued — GitOps: Git as the source of truth for what's running

With Kubernetes in the picture and more than one service deploying into it, a new question becomes hard to answer quickly: what's actually running in the cluster right now, and who changed it? If deploys happen via `kubectl apply` from N different CI pipelines, there's no single place that shows the answer, and manual `kubectl edit` sessions during an incident quietly drift the cluster away from anything version-controlled.

GitOps is a specific answer to that: **Git holds the desired state, and a reconciler inside the cluster (Argo CD or Flux) continuously makes the cluster match it.** Nobody runs `kubectl apply` by hand; they open a pull request against the deployment repo instead.

```text
Developer pushes code
        ↓
GitHub Actions runs tests
        ↓
A Docker image is built and published
        ↓
The deployment repository is updated
        ↓
Argo CD detects the change
        ↓
Kubernetes applies the desired state
```

<div class="flow-diagram" role="img" aria-label="End-to-end delivery pipeline: developer pushes to GitHub, GitHub Actions builds and pushes an image to the container registry, the deployment repository is updated, Argo CD detects the change, and Kubernetes applies the desired state.">
  <div class="flow-step"><strong>Developer</strong><small>git push</small></div>
  <div class="flow-arrow">→</div>
  <div class="flow-step"><strong>GitHub</strong><small>PR + review</small></div>
  <div class="flow-arrow">→</div>
  <div class="flow-step"><strong>GitHub Actions</strong><small>test, build, scan</small></div>
  <div class="flow-arrow">→</div>
  <div class="flow-step"><strong>Container Registry</strong><small>image by digest</small></div>
  <div class="flow-arrow">→</div>
  <div class="flow-step"><strong>Deployment Repo</strong><small>manifest updated</small></div>
  <div class="flow-arrow">→</div>
  <div class="flow-step"><strong>Argo CD</strong><small>detects + syncs</small></div>
  <div class="flow-arrow">→</div>
  <div class="flow-step"><strong>Kubernetes</strong><small>desired state applied</small></div>
</div>

In practice that's two repositories and one watcher:

- **The application repository** — the Go/frontend source, plus the CI workflow from Stage 2 that builds and publishes an image tagged by commit SHA.
- **The deployment repository** — Kubernetes manifests (or a Helm chart, or Kustomize overlays) per environment, where a merged PR bumps the image tag to the new SHA.
- **Argo CD**, pointed at the deployment repository, continuously diffing the live cluster state against it and syncing on drift.

The benefits are less about the sync mechanism and more about what it gives you for free: every change to production is a Git commit, so `git log` on the deployment repo *is* your deploy history. Drift is visible — if someone hand-edits a resource, Argo CD flags or reverts it depending on sync policy. Rollback is `git revert`, not a separate rollback procedure to remember under pressure.

**What GitOps isn't:** it's not continuous delivery by itself. Something still has to decide *when* the deployment repo gets updated with a new image tag — that's still a CI step, a bot, or a person opening a PR. GitOps is the reconciliation half of the story, not the promotion half.

**When it's unnecessary overhead:** one or two services on a single environment, deployed by one or two people who already know exactly what's running. The auditability and drift-detection Argo CD buys you matter most once enough people and enough services exist that "what's running in prod" stops being something any one person can answer from memory.

## Stage 4, continued — configuration, secrets, and observability

**Configuration and secrets** get one more layer of indirection once Kubernetes is in the picture: application config as env vars or mounted files sourced from ConfigMaps, sensitive values from Kubernetes Secrets — which are base64-encoded, not encrypted, by default, worth knowing before treating them as a real security boundary. The moment secrets live in a cloud secret manager rather than hand-typed into `kubectl create secret`, an External Secrets Operator (or equivalent) becomes the bridge, syncing from Vault/Secrets Manager into Kubernetes Secrets so there's one real source of truth and rotation actually propagates.

**Observability** is where "the app is running" and "the app is working" stop being the same question. Four categories, each answering something different:

- **Metrics** — pre-aggregated numbers, cheap to store, good for "is this trending wrong." Prometheus is the default collector.
- **Logs** — high-cardinality, expensive per byte, best when structured (JSON, with a correlation ID) instead of free text. Loki is a common low-cost aggregator.
- **Traces** — the causal path of a single request across services, essential the moment "which of these five services is actually slow" stops being answerable by staring at logs. OpenTelemetry for instrumentation, Tempo or Jaeger as the backend.
- **Alerts** — the layer that turns the first three into "someone gets paged," via Alertmanager or equivalent routing.

<div class="flow-diagram flow-diagram--vertical" role="img" aria-label="Observability architecture: applications and Kubernetes emit data through OpenTelemetry, Prometheus, and log shippers, producing metrics, logs, and traces, which feed Grafana dashboards and alerts.">
  <div class="flow-step"><strong>Applications & Kubernetes</strong><small>your services, pods, nodes</small></div>
  <div class="flow-arrow">↓</div>
  <div class="flow-step"><strong>OpenTelemetry / Prometheus / Log shippers</strong><small>collection layer</small></div>
  <div class="flow-arrow">↓</div>
  <div class="flow-step"><strong>Metrics, Logs, and Traces</strong><small>stored and indexed</small></div>
  <div class="flow-arrow">↓</div>
  <div class="flow-step"><strong>Grafana Dashboards and Alerts</strong><small>what a human actually looks at</small></div>
</div>

The four golden signals — **latency, traffic, errors, saturation** — are the right starting set of dashboards: API latency (p50/p95/p99, not just average), request volume, error rate, and resource saturation (CPU/memory, database connection pool usage, queue backlog). Add Kubernetes pod restarts and deployment health once GitOps is in the picture, so a bad rollout shows up as a graph, not a support ticket.

The mistake worth naming explicitly: **dashboards and alerts are not free just because they're easy to add.** A dashboard nobody looks at is decoration. An alert that pages someone for something that isn't actionable trains that person to ignore pages — and the next real incident gets the same shrug. Start with the golden signals and the handful of alerts tied directly to user-facing symptoms (error rate, latency, hard saturation), and add more only when a specific incident proves you were missing a signal — not preemptively.

## Security, developer experience, and build-versus-buy

**Security** belongs woven through everything above, not bolted on at the end — but a few foundations deserve naming directly: least-privilege IAM and Kubernetes RBAC (a service account that can do everything is a liability the moment it's compromised), short-lived credentials and workload identity over long-lived static keys, network policies restricting which Pods can talk to which, audit logs on anything touching production, and a hard line between who can deploy application code and who can change infrastructure. None of this needs to be a dedicated security team's job at fifteen engineers — it needs to be nobody's afterthought.

**Developer experience** is the multiplier on everything else: a good platform makes the correct path the easiest path, not the most-documented one. In practice that's a repo template so a new service starts with CI, a Dockerfile, and health checks already wired up; a Docker Compose file so local dev doesn't require touching the real cluster; and — as team and service count grow — preview environments per PR and self-service deploys, so shipping a change doesn't require asking someone who "knows how the deploy works." This is the seed of a "paved road," and it's worth resisting the urge to formalize it into a dedicated platform team before there's more than one team actually walking the road.

**Build versus buy** is the question underneath half the decisions in this post. Every tool named above has a managed or hosted equivalent: a managed Kubernetes control plane instead of running your own, a hosted Prometheus/Grafana instead of self-hosting the stack, Argo CD as a managed service instead of another thing your team operates. The trade-offs are the same shape every time — engineering time, infrastructure cost, operational burden, vendor lock-in, reliability, and how fast you need to move — and the honest answer for most startups, most of the time, is to buy the undifferentiated parts and spend engineering time on the product. A platform is an investment that pays off at sufficient scale; below that scale, the "investment" is just a distraction with good intentions.

## The staged roadmap

<div class="flow-diagram" role="img" aria-label="Startup DevOps maturity roadmap across four stages: first production deployment, growing application, multiple services and teams, and platform maturity.">
  <div class="flow-step"><strong>Stage 1</strong><small>First deploy</small></div>
  <div class="flow-arrow">→</div>
  <div class="flow-step"><strong>Stage 2</strong><small>Growing app</small></div>
  <div class="flow-arrow">→</div>
  <div class="flow-step"><strong>Stage 3</strong><small>Multiple services & teams</small></div>
  <div class="flow-arrow">→</div>
  <div class="flow-step"><strong>Stage 4</strong><small>Platform maturity</small></div>
</div>

**Stage 1 — first production deployment.** Essential: GitHub with pull requests, GitHub Actions for lint/test/build, Docker, managed application hosting, a managed database, cloud logging, and basic uptime monitoring. Optional: CODEOWNERS, semantic versioning. Premature: Kubernetes, Terraform, a dedicated observability stack, multiple environments.

**Stage 2 — growing application.** Essential: Infrastructure as Code, separated staging and production, a container registry with real tags, an automated database-migration step, structured logging with basic metrics, real secrets management, and a rollback procedure someone has actually rehearsed. Optional: canary-style gradual rollout, a full APM tool. Premature: Kubernetes (usually), a multi-region footprint, a platform team.

**Stage 3 — multiple services and teams.** Essential: Kubernetes (only once the "many services" pressure is real), Helm or Kustomize, Argo CD, Prometheus and Grafana, OpenTelemetry, centralized logs, clear service ownership, and reusable CI workflows across services. Optional: formal SLOs, a dedicated on-call rotation tool. Premature: multi-cluster-per-region, a full internal developer platform, policy-as-code enforcement.

**Stage 4 — platform maturity.** Consider: a deliberate multi-cluster strategy, progressive delivery (canary as the default, not the exception), preview environments per PR, policy enforcement (OPA or similar), internal developer platform capabilities, automated security and compliance controls, real cost visibility, and rehearsed disaster recovery. Even here, most of this is "consider," not "require" — plenty of successful companies run Stage 3's toolset well past a hundred engineers.

The distinction that matters at every stage isn't which column a tool is in today — it's asking, before adopting anything, *which specific failure, cost, or velocity problem is this fixing right now.* "We might need it later" is not a yes.

## Closing the loop

Line up the whole sequence and it isn't a list of unrelated tools — it's one story about deploy safety and blast radius, told from five different angles. GitHub and branch protection reduce the blast radius of a single bad commit. CI reduces the blast radius of a broken build reaching anyone. GitOps and Argo CD reduce the blast radius of an untracked change to the cluster. Observability reduces the blast radius of *not noticing* something's wrong. Kubernetes reduces the operational blast radius of running many services by hand. Every tool in this post is answering some version of the same question — how do we ship this safely — at a different point in the company's growth.

That's the whole thesis again, and it's worth ending on it exactly because it doesn't get less true as the toolchain grows: DevOps maturity is not measured by how many tools a company uses. It's measured by how safely, quickly, and predictably the company can deliver software — and the fastest way to lose that is to adopt Stage 4's answer to a Stage 1 problem.

This post is the front door to a longer series working through each of these topics at real depth — Linux fundamentals, Git workflows, CI internals, Docker, the four-part Kubernetes arc, Terraform, GitOps, the observability and SRE practices, and platform engineering as the capstone. Each one picks up exactly where this post left it at survey depth and goes deep.
