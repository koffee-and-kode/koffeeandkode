---
title: "0 to 1 Million with AWS"
date: 2026-05-16
tags: [aws, cloud, scaling, ai, bedrock, architecture]
excerpt: "An AI startup's AWS bill tells a story. Walk through what services you actually need at each user tier — from the alpha with 10 design partners to a million-user product — and, more importantly, what pressure forces each jump."
read_time: 14
---

Most AWS 101 posts are a glossary: here is S3, here is EC2, here is Lambda. Useful as a reference, useless as a guide — because the question you're actually asking when you're building isn't *what is this service*, it's *do I need it yet*. A startup at 10 users and a startup at 100,000 users are running the same business but a very different cloud, and the interesting question is what pressure forces each jump.

This post walks an AI startup from zero to a million users across six tiers. AI startup because Bedrock, vector search, and inference cost behave differently from a CRUD app, and that shapes the trajectory in ways worth being explicit about. At each tier I'll cover three things: the pressure that forces the change, the services that come in, and — equally important — what you can still safely *not* have. The biggest mistake at every tier is adopting the next tier's complexity before you have its problems.

A note on framing: "users" is doing a lot of work. For a B2C consumer AI product, 10,000 users is small. For a B2B SaaS where each customer is a $50k ACV, 10,000 users might be a $50M ARR company. The numbers below are calibrated to a B2C-shaped AI product (chat, generation, agent workflows) with mid-intensity per-user usage. Adjust mentally if your shape is different.

This post moves fast through the acronyms — if something below is unfamiliar, there's a companion [AWS glossary]({% post_url cloud/2026-07-13-aws-glossary-0-to-1-million %}) covering every term used here, each one linking back to exactly this post.

## Tier 1 — 0 to 10 users (alpha, design partners)

You're trying to find out if anyone wants this. The architecture should reflect that.

**The pressure:** none, really. Your bottleneck is not infrastructure, it's whether the product works and people will use it. Every hour spent on infra is an hour not spent talking to users.

**What you actually need:**

- **A single compute target.** App Runner or Lambda + API Gateway. Whichever you can ship in an afternoon. Do not stand up ECS, do not write a Terraform module, do not build a VPC.
- **Managed Postgres** — RDS single-AZ on the smallest instance, or Aurora Serverless v2 with a 0.5 ACU floor. Database choices are sticky; pick boring and move on.
- **S3** for any artifact storage (uploaded files, generated outputs, model artifacts).
- **Bedrock** for inference. This is the unfair advantage of building an AI startup *now* — you don't run a model, you call an API. No GPUs, no triton, no autoscaling group of inference workers. A `bedrock-runtime:InvokeModel` call and you have Claude / Llama / Mistral in your product.
- **Route 53** for DNS, **ACM** for TLS. Both effectively free.
- **CloudWatch** for logs, default-tier. You will read them by hand. That's fine for 10 users.

**What you don't need yet:** VPC design, multi-AZ anything, auth provider (you can ship with magic links via SES and a `users` table), CDN, observability tooling beyond `console.log` to CloudWatch, IaC, CI/CD beyond `git push` triggering a deploy.

**The shape:** one app server, one database, S3, Bedrock. You could draw the whole architecture on a napkin and have room for the coffee stain.

**Rough cost:** $50–$200/month if you're careful. Most of it is Bedrock tokens, not infra. The infra under an idle alpha product is *cheap* — the trap people fall into is provisioning for a load that doesn't exist yet.

## Tier 2 — 10 to 100 users (private beta)

A handful of design partners is now a small community. People show up at irregular hours, send feedback, and a couple of them start to use the product seriously enough that downtime is embarrassing.

**The pressure:**
1. You can't push at 3pm anymore without breaking someone's flow.
2. Auth-by-magic-link starts to feel hacky when one customer wants SSO.
3. Bedrock costs are no longer a rounding error and you have no idea which user or which feature is responsible.
4. You're losing track of what's in the database.

**What comes in:**

- **Cognito** (or a third-party — Clerk, WorkOS, Auth0) for proper user accounts, sessions, and the door to SSO/SAML later. Worth choosing carefully because migrating auth providers is painful.
- **ElastiCache (Redis)** the moment you have sessions, rate limits, or a hot cache key. Cheap, transformative.
- **CloudFront** in front of the app for static assets and global TLS termination. Pages snap to life instead of feeling like they're served from one US region.
- **Secrets Manager** for the API keys you've been keeping in `.env` files. Cheap insurance.
- **A CI/CD pipeline** — GitHub Actions is enough, but at least one of: tests on PR, deploy on merge, ability to roll back.
- **Bedrock Guardrails** if you're putting model output in front of users you don't fully trust. Cheap to enable, catches the obvious stuff.
- **Per-user cost tracking** at the application layer. Tag every Bedrock call with `user_id`, `feature`, and log token counts. You will thank yourself later.

**What you still don't need:** multi-AZ databases (probably), a proper VPC topology, EKS, Kinesis, a data warehouse, an on-call rotation, fine-grained IAM. RDS with automated backups and a few CloudWatch alarms is still enough.

**The shape:** still one application tier, but now in front of a CDN and behind an auth provider. Redis next to the database. Bedrock calls instrumented with usage metadata going into a `usage_events` table.

**Rough cost:** $300–$1,500/month, increasingly dominated by Bedrock token spend. This is when you start noticing that one verbose user costs more than ten quiet ones.

## Tier 3 — 100 to 1,000 users (public launch)

Now there's a launch post, a waitlist, and a usage curve that doesn't follow your timezone anymore. Things break in ways you didn't predict.

**The pressure:**
1. A single instance can no longer absorb traffic spikes — you need to be able to scale horizontally.
2. A bad deploy now affects hundreds of paying users.
3. Bedrock latency starts to matter; users are noticing the 2-second pause.
4. You have actual usage data and want to do something with it (analytics, billing, abuse detection).
5. The CFO (or your co-founder wearing that hat) wants to know unit economics per user.

**What comes in:**

- **ALB + ECS Fargate** with multiple tasks. Now you have horizontal scaling, blue/green deploys, and a health check that takes a sick instance out automatically. (Or stay on Lambda if your workload fits — but most AI workloads with streaming responses outgrow Lambda's 15-minute, no-streaming-from-API-Gateway constraints quickly.)
- **A real VPC** — public subnets for the ALB, private subnets for app and database, NAT for egress. Set this up *once*, properly, before you have data to migrate.
- **RDS Multi-AZ** or Aurora with a reader. The cost roughly doubles; the recovery story stops being "we hope the AZ comes back".
- **SQS** for any work that doesn't need to happen in the request: sending notifications, kicking off long-running Bedrock jobs, billing reconciliation. The moment you have a "retry this later" need, SQS is the answer.
- **EventBridge** for fan-out and cross-service eventing. A `user.created` event lights up onboarding, billing, and a Slack alert without coupling those three to your signup endpoint.
- **Kinesis Data Streams** for usage telemetry. Every Bedrock call, every feature use, fire-and-forget into Kinesis. You'll process it elsewhere and you stop coupling analytics to the hot request path.
- **WAF** in front of the ALB. Bot traffic on a public AI product is *aggressive* — content scrapers, prompt-injection probes, credential stuffers. WAF managed rules will block 80% of the obvious garbage cheap.
- **X-Ray** or a third-party APM (Datadog, Honeycomb) for distributed traces. The first time you debug "why is this user's request slow" you'll wish you'd had this two tiers ago.
- **A proper IaC story** — Terraform or CDK. Manual console changes from here on become technical debt.

**The shape:** ALB → Fargate cluster → RDS Multi-AZ + Redis, with SQS and EventBridge wiring async work, Kinesis collecting telemetry, CloudFront in front of everything. Bedrock is still just an API call but now wrapped in retry, timeout, and per-tenant rate-limiting logic.

**What you still don't need:** Kubernetes, multi-region, SageMaker, a Redshift warehouse, your own vector database. Bedrock Knowledge Bases with OpenSearch Serverless covers most RAG use cases now and is one fewer thing to operate.

**Rough cost:** $3k–$10k/month. Bedrock is now likely the biggest line item, with Fargate compute second and data transfer creeping up as a surprise third.

## Tier 4 — 1,000 to 10,000 users (growing)

The product works, the team has grown, and the AWS bill is no longer something one person can hold in their head. You're optimizing in two directions at once: reliability (because customers complain) and cost (because the gross margin matters to the next round).

**The pressure:**
1. Bedrock spend is now five figures monthly and on-demand pricing hurts. You need predictable cost without throttling users.
2. Analytics-on-Postgres is slowing down OLTP queries. You need a separate place for usage data.
3. Customer-success wants to ask "show me every prompt this user sent" and you don't have a good answer.
4. Auto-scaling is reactive — there's still a "warm up" period that's noticeable.
5. Compliance starts to come up in sales conversations (SOC 2, sometimes HIPAA).

**What comes in:**

- **Bedrock Provisioned Throughput** for your hot models. You commit to throughput, you get predictable cost and predictable latency. Run an analysis on your token mix first — provisioned only wins above a certain utilization threshold.
- **Aurora with multiple read replicas**, or split into Aurora (transactions) + a managed analytics target. **Kinesis Data Firehose** → S3 (Parquet) → **Athena** is the cheapest analytics warehouse you'll ever run: pay-per-query, no cluster to manage.
- **DynamoDB** for the high-fanout, low-relational data: session state, rate-limit counters, feature flag evaluations, agent step state. The right tool the moment your access pattern is "look up by one key, very fast, at very high rate".
- **Step Functions** for multi-step agent workflows. Once an "agent" call is actually a chain of `retrieve → think → tool-call → think → respond`, expressing it as code in a single request handler becomes fragile. Step Functions makes it visible, restartable, and observable.
- **API Gateway with usage plans** if you start to expose APIs to customers — per-key throttling and quota for free.
- **Reserved capacity / Savings Plans** on the stable parts of the workload. Easy 30–50% cost cut on compute you're confident you'll keep running.
- **AWS Backup**, **AWS Config**, basic **GuardDuty** — the cheap-to-enable compliance hygiene that auditors expect to see turned on.
- **A real on-call rotation** with **PagerDuty** (or **AWS Incident Manager**) wired to CloudWatch alarms and Sentry. The team is now too big for "Slack on whoever's online".

**The shape:** the front of the system looks similar to Tier 3 but the back has split — transactional data on Aurora, analytics flowing through Kinesis to S3, session/agent-step data on DynamoDB, agent orchestration on Step Functions, provisioned Bedrock capacity behind a thin tenant-aware routing layer.

**Rough cost:** $20k–$80k/month. Provisioning decisions (reserved vs on-demand, provisioned-throughput on Bedrock vs not) are now individually worth a senior engineer's salary to get right.

## Tier 5 — 10,000 to 100,000 users (scaled)

You're a real company. The platform-engineering function exists. Failure modes are weirder — not "the server is down" but "a slow query in one tenant is starving another tenant's threadpool".

**The pressure:**
1. Multi-tenant noisy-neighbor problems are now happening. One customer's batch job tanks latency for everyone.
2. Region-level outages now exceed your tolerance budget. A 4-hour `us-east-1` blip is a 4-hour outage of your business.
3. The data team wants a real warehouse. Athena was great for ad-hoc; reporting now needs joins, dbt, and dashboards.
4. Custom embeddings, custom fine-tunes, and proprietary models start to differentiate the product.
5. The security team is real, has opinions, and is starting to push back on the IAM sprawl from Tier 3.

**What comes in:**

- **Multi-region** — at minimum, an active-passive failover region with **Aurora Global Database**, S3 cross-region replication, Route 53 health-check failover. Full active-active is much harder and rarely needed before Tier 6.
- **Bedrock cross-region inference profiles** to route inference to the lowest-loaded region transparently. Real win on tail latency.
- **Redshift** (or Snowflake / Databricks) as the analytics warehouse. Glue or Airbyte to land transactional data alongside the Kinesis-fed usage stream.
- **SageMaker** if you've decided to fine-tune or to train custom embedding models. The build-vs-buy moment for ML infrastructure: most startups should still buy from Bedrock as long as possible.
- **OpenSearch** as a proper cluster (not Serverless) once your vector + lexical search workload is too large or too custom for Bedrock Knowledge Bases.
- **AWS Organizations** with separate accounts per environment (prod, staging, dev) and per-business-unit if relevant. **Service Control Policies** to enforce guardrails (no `us-east-2` resources, mandatory tags, no internet-exposed S3).
- **Transit Gateway** as the VPC count grows.
- **EKS** *maybe* — only if you have workloads that genuinely benefit (large stateful inference, complex pod-level scheduling, Helm-shaped tooling your team already knows). Most companies should stay on ECS as long as they can. EKS is a tax you pay forever.
- **CloudFront Functions / Lambda@Edge** for request-shaping at the edge (per-region routing, auth-token validation, A/B bucketing).
- **A FinOps practice** — someone whose explicit responsibility is reading the Cost & Usage Report, tagging discipline, and reserved-capacity strategy. Saves a percentage of a multi-hundred-thousand-dollar bill.

**The shape:** two-region (active-passive) deployment, separate AWS accounts for prod and non-prod, transactional path on Aurora Global, analytics path landing in Redshift via Kinesis+Glue, vector workloads either on Knowledge Bases or a dedicated OpenSearch cluster, agent orchestration on Step Functions, and inference routed across regions on Bedrock provisioned profiles.

**Rough cost:** $150k–$600k/month, with single-digit-percent improvements in token efficiency or compute utilization now worth six-figure annualized savings.

## Tier 6 — 100,000 to 1,000,000 users (at scale)

You are operating at a scale where small architectural choices have large dollar and reliability consequences. The interesting work is no longer "what service do we adopt next" — it's "how do we run this one really well".

**The pressure:**
1. Cost per user must continue to come down to preserve gross margin as you grow.
2. Tail latency at the 99.9th percentile is the user experience for tens of thousands of people *every day*. The averages don't matter.
3. Regulatory and data-residency requirements (GDPR in EU, data residency in specific verticals) shape architecture.
4. The blast radius of any single change is huge. A bad config rollout is a public incident.

**What comes in (and what gets harder):**

- **Active-active multi-region**, with conflict resolution strategies for the parts of the data model that allow it and pinned-region tenancy for the parts that don't.
- **MSK (managed Kafka)** if Kinesis-shaped streaming has outgrown Kinesis: higher fan-out, longer retention, more flexible consumer model. Most teams should stick with Kinesis longer than they think.
- **Dedicated tenancy / VPC endpoints / PrivateLink** for enterprise customers who refuse to send traffic over the public internet.
- **Custom inference infrastructure** — if your token volume and latency profile justify it, this is the point at which running a fleet of self-hosted models on EC2 with **Inferentia/Trainium** or GPUs can beat Bedrock on cost. Worth a real build-vs-buy analysis. For most companies the answer is still: keep using Bedrock and put the engineering hours into product.
- **Chaos engineering**, **GameDay** rehearsals, formal disaster-recovery testing on a regular cadence. The "we'd survive a regional outage" claim becomes a thing you've actually rehearsed.
- **A platform engineering team** whose product is the developer experience for everyone else. Internal abstractions over Terraform, golden paths, opinionated service templates.
- **Compliance breadth** — SOC 2 Type II, ISO 27001, HIPAA, FedRAMP, depending on where you sell. Most of this is process, but the cloud configuration (KMS, CloudTrail, Config, GuardDuty, Security Hub, Macie) is the substrate it sits on.

**The shape:** a multi-region, multi-account topology where the "AWS architecture diagram" is no longer one diagram but a set of them — control plane, data plane, ML plane, analytics plane — each evolving on its own cadence. The interesting questions are organizational and economic as much as technical.

**Rough cost:** $1M+/month. At this scale the unit economics of inference are typically the largest single optimization lever; everything else is rounding.

## The meta-lesson — adopt one tier at a time

The single most expensive mistake I've seen startups make is to architect for Tier 5 while operating at Tier 1. Multi-region active-active in front of 12 users is not insurance; it's complexity that slows you down for years and a bill that doesn't end. The teams that scale well adopt one tier at a time, when the pressure that forces the jump is *actually felt*, and not before.

A useful test before adopting any service: *which specific failure or cost or velocity problem is this fixing right now?* If the answer is "we might need it later" or "it would be more correct", you're not at that tier yet. Wait until a real outage, a real bill, a real lost deal, or a real velocity ceiling tells you it's time.

The corollary is that the *order* of adoption isn't arbitrary. SQS before EKS. Multi-AZ before multi-region. Provisioned-throughput Bedrock before SageMaker fine-tuning. The progression above isn't a checklist; it's a sequence of pressure responses, each one earning the next.

Bedrock changes one thing fundamentally for AI startups: you can spend two tiers worth of headcount on *product* before you start spending it on inference infrastructure. Use that gift. The companies that win in AI right now are the ones that get to Tier 3 with a sharper product than their competitors, not a more elaborate cloud.
