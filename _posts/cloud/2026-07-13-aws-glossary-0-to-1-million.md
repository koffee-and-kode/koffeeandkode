---
title: "AWS Glossary: The Terms Behind 0 to 1 Million"
date: 2026-07-13
tags: [aws, cloud, glossary, reference]
excerpt: "Every AWS service, acronym, and pricing term from the 0-to-1-million-users post, defined — plus the popular and open-source alternatives to each, and a 1-to-5 gut check on whether the AWS-native pick is actually the right one."
read_time: 14
---

[0 to 1 Million with AWS]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}) moves fast through a lot of AWS vocabulary — by Tier 4 it's assuming you already know what a Firehose is. This is the glossary for that post: one definition per term, alphabetical, each linking back to the tier it's actually used in so you can jump straight to the context.

Two things this version adds beyond a plain definition list. First, **alternatives** — the equivalent service on GCP/Azure or from a common third-party vendor, plus an open-source option where a real one exists (some of these, like IAM or a VPC, don't have a meaningful "alternative" because every cloud has its own version of the same primitive — those are called out as such). Second, a **recommendation score, 1–5**, which rates the AWS-native pick specifically — not "is this a good idea," but "of the options listed, is AWS's own version the one I'd actually reach for." A 5 means don't overthink it; a 2 means it works, but there's usually a better-regarded option sitting right next to it.

This isn't a general AWS glossary — it only covers terms that show up in that specific post. For the broader "what is DevOps" vocabulary (Kubernetes, Terraform, GitOps, and friends), see [DevOps 101]({% post_url devops/2026-07-12-devops-101-building-a-platform-from-scratch %}) instead.

## The glossary

### ACM (AWS Certificate Manager)
Issues and auto-renews TLS certificates for free when they're attached to AWS-managed resources like a load balancer or CloudFront. No cert-expiry pages at 2am. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

- **Alternatives:** GCP-managed certs, Azure App Service managed certs
- **Open-source:** Let's Encrypt / certbot (self-managed)
- **Recommendation:** 5/5 — free, zero-maintenance, no reason to reach for anything else once you're on AWS

### ACV (Annual Contract Value)
The yearly dollar value of one customer's contract. It's the metric that makes "10,000 users" mean wildly different things: for a B2C app that's a small product; for a B2B tool at $50k ACV, that's a $500M business. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

*A business metric, not a product — no alternative to name or score.*

### ALB (Application Load Balancer)
An AWS load balancer that understands HTTP, so it can route by path or hostname and run health checks against individual instances, pulling unhealthy ones out of rotation automatically. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

- **Alternatives:** GCP Cloud Load Balancing, Azure Application Gateway
- **Open-source:** NGINX, HAProxy, Envoy (self-hosted)
- **Recommendation:** 5/5 — boring, reliable, exactly what you want from a load balancer

### API Gateway
A managed front door for HTTP APIs: handles routing, throttling, and auth in front of Lambda (or another backend) without you running a web server for it. At scale it also offers **usage plans** — per-API-key rate limits and quotas — useful the moment you expose an API to customers directly. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners) · [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** GCP API Gateway / Apigee, Azure API Management
- **Open-source:** Kong, KrakenD, Tyk
- **Recommendation:** 3/5 — capable, but the console/DX is clunkier than Kong or an Envoy-based gateway; fine as the default if you're already all-in on Lambda

### Athena
A query engine that runs SQL directly against files sitting in S3 (typically Parquet), with no cluster to provision or manage. You pay per query scanned, which makes it about the cheapest possible analytics warehouse for ad-hoc use. → [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** GCP BigQuery, Azure Synapse serverless SQL, Snowflake
- **Open-source:** Trino / Presto, DuckDB (smaller scale)
- **Recommendation:** 4/5 — great for ad-hoc queries against S3; BigQuery has nicer ergonomics if you're not AWS-committed

### Aurora Global Database
A version of Aurora that replicates across AWS regions with low replication lag, so a secondary region can take over quickly if the primary has an outage. → [Tier 5]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-5--10000-to-100000-users-scaled)

- **Alternatives:** GCP Cloud Spanner, Azure Cosmos DB
- **Open-source:** CockroachDB, YugabyteDB
- **Recommendation:** 3/5 — solid, but Spanner and CockroachDB have a stronger native multi-region consistency story if that's the actual requirement

### Aurora Serverless v2
A Postgres/MySQL-compatible database that scales its own capacity up and down automatically, billed by the fraction of a "capacity unit" actually used. A 0.5 ACU floor means it never scales below half a unit — cheap enough to leave running for an idle alpha product, without paying for a fixed-size instance around the clock. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

- **Alternatives:** GCP AlloyDB / Cloud SQL, Azure SQL Database serverless, Neon, PlanetScale
- **Open-source:** Postgres/MySQL self-hosted; Neon (built on an open-core model)
- **Recommendation:** 4/5 — good default; Neon's branching model is genuinely nicer for dev/preview workflows if that matters to you

### AWS Backup
A managed service that centralizes and schedules backups across RDS, DynamoDB, EFS, and other AWS storage, instead of each service running its own ad-hoc backup script. → [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** GCP Backup and DR, Azure Backup
- **Open-source:** Velero (Kubernetes-centric), Restic
- **Recommendation:** 4/5 — unglamorous but does the job; turn it on and forget about it

### AWS Config
Continuously records the configuration of your AWS resources and can flag when one drifts from a rule you've defined (e.g., "no S3 bucket should be public"). Compliance auditors like seeing this turned on. → [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** GCP Security Command Center, Azure Policy
- **Open-source:** Cloud Custodian, OPA / Conftest (policy-as-code)
- **Recommendation:** 3/5 — useful, but the rules DSL is awkward; many teams layer Cloud Custodian on top anyway

### AWS Organizations
Lets you run multiple AWS accounts (one per environment, one per business unit) under one billing and governance umbrella, with **Service Control Policies (SCPs)** as the guardrails that apply across every account underneath — e.g., "no resources outside these two regions," enforced at the account level, not just by convention. → [Tier 5]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-5--10000-to-100000-users-scaled)

- **Alternatives:** GCP Resource Manager + Org Policies, Azure Management Groups + Azure Policy
- **Open-source:** none meaningful — this is inherently a cloud-provider governance primitive
- **Recommendation:** 5/5 — set this up early; there's no real substitute once you have more than one account

### Bedrock
AWS's managed access point to foundation models (Claude, Llama, Mistral, and others) via a single API. You call `bedrock-runtime:InvokeModel` and get inference back — no GPUs, no model-serving infrastructure, no autoscaling group of inference workers to run yourself. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

- **Alternatives:** GCP Vertex AI, Azure AI Foundry, calling a provider's API directly (e.g. Anthropic's own API)
- **Open-source:** self-hosted serving via vLLM / Ollama, or LiteLLM as a provider-agnostic routing layer
- **Recommendation:** 4/5 — a genuinely good convenience layer; if you specifically need the newest frontier model on day one, going direct to the provider sometimes beats Bedrock's onboarding lag

### Bedrock cross-region inference profiles
Routes an inference request to whichever region currently has capacity or the lowest load, transparently, instead of your code hardcoding one region. Mainly a tail-latency win. → [Tier 5]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-5--10000-to-100000-users-scaled)

- **Alternatives:** Vertex AI multi-region endpoints, Azure OpenAI regional deployments behind Front Door
- **Open-source:** a hand-rolled routing layer (LiteLLM, or your own)
- **Recommendation:** 4/5 — does exactly what it says, worth turning on once tail latency actually matters

### Bedrock Guardrails
A filtering layer you attach to a Bedrock model that blocks categories of unwanted output (unsafe content, PII leakage, off-topic responses) before it reaches the user. Cheap insurance the moment model output is user-facing. → [Tier 2]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-2--10-to-100-users-private-beta)

- **Alternatives:** Azure AI Content Safety, Vertex AI safety filters
- **Open-source:** Guardrails AI, NeMo Guardrails, LLM-Guard
- **Recommendation:** 3/5 — a decent baseline; NeMo Guardrails or Guardrails AI give more control if the default categories aren't enough

### Bedrock Knowledge Bases
A managed RAG (retrieval-augmented generation) pipeline: point it at your documents, it handles chunking, embedding, and retrieval, backed by OpenSearch Serverless underneath. Covers most "search your own data" use cases without you standing up a vector database yourself. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

- **Alternatives:** Vertex AI Search, Azure AI Search
- **Open-source:** LangChain or LlamaIndex over your own vector store (pgvector, Qdrant, Weaviate)
- **Recommendation:** 3/5 — fast to stand up, but you trade away fine-grained control over chunking and retrieval quality that a hand-rolled pipeline gives you

### Bedrock Provisioned Throughput
Instead of paying per-token on demand, you commit to a fixed amount of model throughput for predictable cost and latency. Only worth it above a certain utilization — model your actual token volume before committing. → [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** Vertex AI provisioned throughput, Azure OpenAI PTUs (provisioned throughput units)
- **Open-source:** self-hosted serving on your own GPUs (vLLM, TGI) once volume justifies it
- **Recommendation:** 3/5 — the math has to actually pencil out; run the utilization numbers before committing, not after

### Chaos engineering / GameDay
Deliberately injecting failure (killing an instance, cutting a region's network) in a controlled way to verify your system actually survives what you claim it survives, instead of hoping. A GameDay is the scheduled, team-wide version of this exercise. → [Tier 6]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-6--100000-to-1000000-users-at-scale)

- **Alternatives:** Gremlin (commercial, multi-cloud), Azure Chaos Studio
- **Open-source:** Chaos Mesh, LitmusChaos (both Kubernetes-native)
- **Recommendation:** 3/5 for AWS's own tooling here — it's thin; most teams doing this seriously reach for Gremlin or Chaos Mesh regardless of cloud

### CloudFront
AWS's CDN: caches static assets and terminates TLS at edge locations close to the user, so a page load doesn't round-trip to your one origin region every time. → [Tier 2]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-2--10-to-100-users-private-beta)

- **Alternatives:** GCP Cloud CDN, Azure Front Door, Cloudflare, Fastly
- **Open-source:** none meaningful — a CDN is inherently a network of PoPs, not something you self-host
- **Recommendation:** 3/5 — works fine, but Cloudflare is frequently cheaper with a better free tier; worth comparing before defaulting to CloudFront just because you're on AWS

### CloudFront Functions / Lambda@Edge
Lightweight code that runs at CloudFront's edge locations, before a request reaches your origin — used for things like per-region routing or validating an auth token without a round-trip inland. → [Tier 5]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-5--10000-to-100000-users-scaled)

- **Alternatives:** Cloudflare Workers, Fastly Compute@Edge
- **Open-source:** none meaningful
- **Recommendation:** 2/5 — Cloudflare Workers is widely considered the nicer developer experience in this category by a real margin

### CloudWatch
AWS's default logging and metrics service. Every AWS resource ships logs and basic metrics here unless you point them somewhere else. Fine to read by hand at ten users; not a substitute for real observability once traffic grows. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

- **Alternatives:** GCP Cloud Logging / Monitoring, Azure Monitor, Datadog, Grafana Cloud
- **Open-source:** Prometheus + Grafana + Loki
- **Recommendation:** 2/5 — functional but clunky, and gets expensive fast; most teams outgrow it and move to Datadog or self-hosted Prometheus/Grafana within a year

### Cognito (and third-party equivalents: Clerk, WorkOS, Auth0)
Managed user-account systems: sign-up, login, sessions, and eventually **SSO/SAML** (single sign-on via the enterprise identity-federation standard SAML) for customers who need it. Worth choosing carefully — migrating auth providers later is genuinely painful. → [Tier 2]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-2--10-to-100-users-private-beta)

- **Alternatives:** Clerk, WorkOS, Auth0, GCP Identity Platform, Azure AD B2C
- **Open-source:** Keycloak, Ory (Kratos/Hydra), Authentik
- **Recommendation:** 2/5 for Cognito specifically — notoriously clunky DX and confusing pricing tiers; most engineers who've used both prefer Clerk/WorkOS, or self-hosted Keycloak if you want to own it

### Compliance frameworks (SOC 2 Type II, ISO 27001, HIPAA, FedRAMP)
Third-party audits/certifications that prove your controls actually work over time (SOC 2 Type II), meet an international security standard (ISO 27001), protect health data (HIPAA), or meet U.S. federal requirements (FedRAMP). Mostly process, but they lean on the cloud configuration underneath (KMS, CloudTrail, Config, GuardDuty, Security Hub, Macie) as their evidence. → [Tier 6]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-6--100000-to-1000000-users-at-scale)

*A regulatory/audit category, not a product — the AWS tools that support it (KMS, CloudTrail, etc.) are each scored on their own elsewhere in this glossary.*

### DynamoDB
A managed key-value/NoSQL database built for one access pattern done extremely fast: look something up by a known key, at very high request rates, with predictable low latency. The right tool for session state, rate-limit counters, and per-request agent state — the wrong tool the moment you need joins or ad-hoc queries. → [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** GCP Firestore / Bigtable, Azure Cosmos DB
- **Open-source:** ScyllaDB, Apache Cassandra, FoundationDB
- **Recommendation:** 4/5 — genuinely excellent for its access pattern; Cosmos DB is the closest multi-cloud analog if avoiding lock-in matters more to you

### ECS Fargate
Runs containers without you managing the underlying servers — you describe the container, Fargate finds and manages the compute. The "just run my Docker image" tier between a single VM and full Kubernetes. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

- **Alternatives:** GCP Cloud Run, Azure Container Apps / Container Instances
- **Open-source:** Nomad (self-hosted scheduler), or Kubernetes itself
- **Recommendation:** 4/5 — a very solid, boring default; Cloud Run is often cited as the nicer developer experience in this category, but Fargate won't let you down

### EKS (Elastic Kubernetes Service)
AWS's managed Kubernetes control plane. Only worth adopting once you have workloads that specifically benefit from Kubernetes' scheduling model — most companies should stay on ECS Fargate longer than they expect to; EKS is operational overhead you keep paying, not a one-time cost. → [Tier 5]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-5--10000-to-100000-users-scaled)

- **Alternatives:** GCP GKE, Azure AKS
- **Open-source:** self-managed via kubeadm or k3s
- **Recommendation:** 3/5 — GKE is widely regarded as the most polished managed Kubernetes offering; EKS works fine but has rougher edges (VPC CNI quirks, slower feature rollout)

### ElastiCache (Redis)
A managed Redis (or Memcached) instance, used for anything that needs to be fast and shared across app instances: sessions, rate-limit counters, a hot cache key that's getting hit constantly. → [Tier 2]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-2--10-to-100-users-private-beta)

- **Alternatives:** GCP Memorystore, Azure Cache for Redis, Upstash, Redis Cloud
- **Open-source:** Redis / Valkey self-hosted, KeyDB
- **Recommendation:** 4/5 — solid managed option; Upstash is worth a look for spiky/serverless workloads on cost

### EventBridge
A managed event bus: one service publishes an event (`user.created`), any number of other services subscribe to it, without the publisher knowing or caring who's listening. Decouples side effects (onboarding, billing, a Slack alert) from the endpoint that triggered them. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

- **Alternatives:** GCP Eventarc / Pub/Sub, Azure Event Grid
- **Open-source:** NATS, Kafka (via MSK or self-hosted)
- **Recommendation:** 4/5 — good for AWS-native fan-out; NATS is the OSS-first pick if avoiding lock-in is a priority

### GuardDuty
Continuously analyzes account activity and network traffic for signs of compromise (an unusual API call pattern, traffic to a known-bad IP) and flags it. Cheap to turn on, genuinely catches things a human wouldn't notice in time. → [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** GCP Security Command Center, Microsoft Defender for Cloud
- **Open-source:** Falco (runtime threat detection), Wazuh (SIEM-style)
- **Recommendation:** 3/5 — decent baseline threat detection for the price; Falco or a real SIEM goes further if you need depth

### IAM (Identity and Access Management)
The system for who — or what service — is allowed to do what in your AWS account. The core discipline is least privilege: a service gets a role scoped to exactly what it needs, never a shared key with admin access — cheap to get right early, expensive to unwind once a dozen things depend on the overprivileged key. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

- **Alternatives:** GCP IAM, Azure RBAC / Entra ID — every cloud has its own version of this same primitive
- **Open-source:** Open Policy Agent (as a policy layer on top of any of them)
- **Recommendation:** 5/5 — non-negotiable, and the same 5/5 applies to whichever cloud you're on; use it properly from day one

### Inferentia / Trainium
AWS's custom chips for running (Inferentia) or training (Trainium) ML models, an alternative to renting GPUs. Only enters the picture once your inference volume is large enough to justify a real build-vs-buy analysis against just calling Bedrock. → [Tier 6]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-6--100000-to-1000000-users-at-scale)

- **Alternatives:** GCP TPUs, Nvidia GPUs on any cloud
- **Open-source:** n/a (this is hardware, not software)
- **Recommendation:** 2/5 — niche; most teams are better served by Nvidia GPUs for the broader framework support, unless you're optimizing cost at genuinely huge scale

### Kinesis (Data Streams and Data Firehose)
**Data Streams** is a fire-and-forget event pipe for high-volume telemetry (every Bedrock call, every feature use) that something downstream processes later, decoupled from the request path that generated it. **Data Firehose** is the variant purpose-built to land that stream straight into S3 (as Parquet, ready for Athena) without you writing the plumbing yourself. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch) · [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** GCP Pub/Sub + Dataflow, Azure Event Hubs, Confluent Cloud
- **Open-source:** Apache Kafka (self-hosted or via MSK), Redpanda, NATS
- **Recommendation:** 3/5 — works, but the Kafka/Redpanda ecosystem is richer; a lot of teams that outgrow Kinesis land on Kafka anyway, so it's worth asking if you'll end up there regardless

### Lambda
A "function as a service" compute option: you deploy a function, AWS runs it on demand and scales it per-request, and you pay per invocation rather than for an always-on server. Paired with API Gateway for HTTP-triggered functions. Its limits (a hard execution-time cap, no true streaming response through API Gateway) are exactly why AI workloads with long or streaming responses tend to outgrow it. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

- **Alternatives:** GCP Cloud Functions / Cloud Run, Azure Functions
- **Open-source:** Knative, OpenFaaS (self-hosted FaaS)
- **Recommendation:** 4/5 — mature and huge ecosystem; Cloud Run edges it out specifically for anything needing longer execution time or true streaming

### Model artifacts
The actual files a trained model produces (weights, config, tokenizer files) that need to be stored somewhere and loaded at inference time. When you're calling Bedrock instead of hosting your own model, this is almost entirely someone else's problem. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

*A concept, not a product — storage is just S3/GCS/Blob Storage; the artifact format (safetensors, ONNX, etc.) matters more than the cloud it sits in.*

### MSK (Managed Streaming for Kafka)
AWS-managed Apache Kafka, for when streaming needs have outgrown Kinesis: higher fan-out, longer message retention, a more flexible consumer model. Most teams should stay on Kinesis longer than they think before reaching for this. → [Tier 6]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-6--100000-to-1000000-users-at-scale)

- **Alternatives:** Confluent Cloud (multi-cloud), Azure Event Hubs (Kafka-compatible)
- **Open-source:** self-hosted Kafka, Redpanda
- **Recommendation:** 3/5 — works, but Confluent Cloud or self-hosted Redpanda often have a better cost/operations trade-off for the same Kafka protocol

### NAT (gateway)
Lets resources in a private subnet (no direct inbound internet access) reach *out* to the internet — to call a third-party API, say — without being reachable from the internet themselves. Part of the standard public/private VPC split. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

- **Alternatives:** GCP Cloud NAT, Azure NAT Gateway
- **Open-source:** a self-managed NAT instance (an EC2 box doing the same job)
- **Recommendation:** 3/5 — notoriously expensive per-GB at real traffic volumes; a self-managed NAT instance is a common cost optimization once the bill actually stings

### OpenSearch (Serverless vs. a real cluster)
A search/vector engine. **Serverless** mode is the low-maintenance option Bedrock Knowledge Bases runs on for typical RAG workloads; a **full cluster** earns its keep once your vector or lexical search workload is too large or too custom for the serverless tier, at the cost of actually operating a cluster. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch) · [Tier 5]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-5--10000-to-100000-users-scaled)

- **Alternatives:** Vertex AI Vector Search, Azure AI Search, Elastic Cloud, Pinecone
- **Open-source:** self-hosted Elasticsearch/OpenSearch, Qdrant, Weaviate, pgvector
- **Recommendation:** 3/5 — capable, but Pinecone or Qdrant are generally considered nicer purpose-built vector databases if that's your primary use case, not an add-on to search

### PagerDuty / AWS Incident Manager
The tools that turn a triggered alarm into an actual page to a human, with escalation policies if the first person doesn't respond. The point at which "Slack whoever's online" stops being a viable on-call strategy. → [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** PagerDuty, Opsgenie
- **Open-source:** Grafana OnCall
- **Recommendation:** 2/5 for AWS Incident Manager specifically — PagerDuty remains the category leader by a wide margin; most teams pick PagerDuty or Opsgenie over AWS's native option

### PrivateLink / VPC endpoints / dedicated tenancy
Ways to let traffic reach an AWS service (or a specific enterprise customer's infrastructure) without ever crossing the public internet. Enterprise customers with strict network policies sometimes require this outright. → [Tier 6]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-6--100000-to-1000000-users-at-scale)

- **Alternatives:** GCP Private Service Connect, Azure Private Link
- **Open-source:** none meaningful
- **Recommendation:** 4/5 — does exactly what it says; worth using the moment a customer or compliance requirement actually needs it

### Redshift (and Snowflake / Databricks as alternatives)
A proper data warehouse: built for complex joins, scheduled transformations (dbt), and dashboards at a scale where Athena's ad-hoc-query model stops being enough. → [Tier 5]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-5--10000-to-100000-users-scaled)

- **Alternatives:** Snowflake, GCP BigQuery, Azure Synapse, Databricks
- **Open-source:** ClickHouse, DuckDB (smaller scale), Trino
- **Recommendation:** 2/5 — Redshift is one of the more frequently migrated-away-from AWS products; Snowflake and BigQuery are widely preferred for ergonomics and cleaner separation of storage and compute

### Reserved capacity / Savings Plans
Committing to a baseline amount of compute usage in advance in exchange for a meaningfully lower price than on-demand — a straightforward 30–50% cut on the part of your workload you're confident will keep running. → [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** GCP Committed Use Discounts, Azure Reserved Instances
- **Open-source:** n/a
- **Recommendation:** 4/5 — straightforward savings once usage is predictable; the only risk is over-committing before you're sure

### Route 53
AWS's DNS service: translates your domain name to the actual IP/resource behind it, and can also do health-check-based failover (route around a region that's down). → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

- **Alternatives:** GCP Cloud DNS, Azure DNS, Cloudflare DNS, NS1
- **Open-source:** BIND / PowerDNS (self-hosted, rarely worth it for DNS specifically)
- **Recommendation:** 4/5 — solid; Cloudflare DNS is a common preference for its free tier and built-in DDoS protection if you're not deep in the AWS ecosystem otherwise

### SageMaker
AWS's platform for training and fine-tuning your own ML models. The build-vs-buy line for most AI startups: keep using Bedrock's off-the-shelf models as long as possible, and only reach for SageMaker once a custom or fine-tuned model is a genuine product differentiator. → [Tier 5]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-5--10000-to-100000-users-scaled)

- **Alternatives:** Vertex AI Training, Azure Machine Learning
- **Open-source:** Kubeflow, Ray (for training/serving), MLflow (for tracking)
- **Recommendation:** 3/5 — capable but genuinely complex; many ML teams prefer the more composable Ray/Kubeflow stack, or Vertex AI's smoother notebook-to-deploy flow

### Secrets Manager
Stores API keys and credentials centrally, with access controlled by IAM, instead of sitting in `.env` files copied between machines (and occasionally into a Git history by accident). → [Tier 2]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-2--10-to-100-users-private-beta)

- **Alternatives:** GCP Secret Manager, Azure Key Vault
- **Open-source:** HashiCorp Vault, Infisical
- **Recommendation:** 4/5 — solid default; Vault remains the more powerful (and more complex) option if you need dynamic, short-lived secrets across multiple clouds

### Single compute target
Shorthand for "pick one place your app runs" (one App Runner service, one Lambda function) instead of splitting workload across multiple compute options before you have a reason to. The Tier 1 advice is to have exactly one. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

*A principle, not a product — no alternative to name or score.*

### SQS (Simple Queue Service)
A managed message queue: drop a job on the queue, a worker picks it up later. The default answer the moment you have a "this doesn't need to block the request, retry it if it fails" need — sending a notification, kicking off a long-running job. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

- **Alternatives:** GCP Pub/Sub, Azure Service Bus / Queue Storage
- **Open-source:** RabbitMQ, NATS, Redis-based queues (BullMQ)
- **Recommendation:** 4/5 — simple, reliable, boring in the best way; there's rarely a good reason to reach for something fancier at this stage

### Step Functions
A managed state machine for multi-step workflows — chaining `retrieve → think → tool-call → respond` for an AI agent, for example — that makes each step visible, individually retryable, and observable, instead of one long fragile function. → [Tier 4]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-4--1000-to-10000-users-growing)

- **Alternatives:** GCP Workflows, Azure Logic Apps / Durable Functions
- **Open-source:** Temporal, Apache Airflow (batch-oriented), Argo Workflows (Kubernetes-native)
- **Recommendation:** 3/5 — works well for AWS-native orchestration; Temporal has a notably better developer experience for complex, long-running workflows if you're willing to run it yourself

### Transit Gateway
A central hub that connects many VPCs together without wiring up a peering connection between every pair individually. Matters once your VPC count grows past a handful. → [Tier 5]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-5--10000-to-100000-users-scaled)

- **Alternatives:** GCP Network Connectivity Center, Azure Virtual WAN
- **Open-source:** n/a
- **Recommendation:** 4/5 — does exactly what it says, and the alternative (a full peering mesh) gets ugly fast without it

### Triton (NVIDIA Triton Inference Server)
Open-source software for serving your own ML models at scale, including autoscaling the fleet of GPUs it runs on. Mentioned mainly as the contrast case: calling Bedrock means you never have to stand this up yourself. → [Tier 1]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-1--0-to-10-users-alpha-design-partners)

- **Alternatives (other self-hosted serving options):** vLLM, Hugging Face TGI (Text Generation Inference), Ray Serve
- **Open-source:** Triton itself is open-source (it's the alternative to Bedrock, not the other way around)
- **Recommendation:** 3/5 — vLLM has largely overtaken Triton in popularity specifically for LLM serving, mostly on setup simplicity; Triton remains stronger for mixed model types beyond just LLMs

### Vector search
Search by semantic similarity (embeddings) rather than exact keyword match — the retrieval half of a RAG pipeline. Bedrock Knowledge Bases with OpenSearch Serverless covers most of this without you standing up a dedicated vector database. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

*A technique, not a single product — see the OpenSearch and Bedrock Knowledge Bases entries above for the actual services and their alternatives/scores.*

### VPC (Virtual Private Cloud) and VPC topology
Your own private, isolated network inside AWS. The standard topology splits it into public subnets (things the internet can reach directly, like a load balancer) and private subnets (your app and database, reachable only from inside the VPC), connected outbound through a NAT gateway. Worth setting up once, properly, before there's data to migrate around it later. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

- **Alternatives:** GCP VPC, Azure VNet — every cloud has its own version of this same primitive
- **Open-source:** n/a
- **Recommendation:** 5/5 — not optional on any cloud; the only real decision is topology, not whether to use it

### WAF (Web Application Firewall)
Sits in front of your load balancer and blocks traffic matching known-bad patterns (bot traffic, common exploit signatures, credential-stuffing patterns) using managed rule sets, without you writing custom filtering logic. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

- **Alternatives:** GCP Cloud Armor, Azure Front Door WAF, Cloudflare WAF
- **Open-source:** ModSecurity, Coraza
- **Recommendation:** 3/5 — decent managed rules; Cloudflare's WAF is frequently cited as more effective out of the box and easier to tune

### X-Ray / APM (Datadog, Honeycomb)
Distributed tracing: follows one request across multiple services so "which of these five services is actually slow" has an answer instead of a guess. X-Ray is AWS-native; Datadog and Honeycomb are the common third-party alternatives. → [Tier 3]({% post_url cloud/2026-05-16-building-a-startup-0-to-1-million %}#tier-3--100-to-1000-users-public-launch)

- **Alternatives:** Datadog, Honeycomb, GCP Cloud Trace, Azure Application Insights
- **Open-source:** Jaeger, Grafana Tempo (both via OpenTelemetry)
- **Recommendation:** 2/5 for X-Ray specifically — a clunky UI and limited query capability compared to Honeycomb or Datadog; most teams that try both end up switching off X-Ray

## Reading it back

None of these are complicated once named — the actual difficulty in the original post was density, not depth. The alternatives and scores above aren't a case for leaving AWS; they're the same "do I need this yet, and is this specific version of it actually good" question the main post asks about *tiers*, just aimed at each *service* instead. If a term still feels shaky after the definition, that's what the tier link is for: go read it in context, where it's solving an actual problem instead of sitting in a list.
