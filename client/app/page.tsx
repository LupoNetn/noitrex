"use client";

import { useEffect, useRef, useState } from "react";
import {
  Zap,
  Calculator,
  Mail,
  ShieldCheck,
  Infinity as InfinityIcon,
  Database,
  Square,
  X,
  Layers,
  XCircle,
  ArrowRight,
  ArrowDown,
  ArrowLeftRight,
  Server,
  Activity,
  Terminal,
  Code2
} from "lucide-react";

export default function Home() {
  const [scrollY, setScrollY] = useState(0);
  const heroRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleScroll = () => setScrollY(window.scrollY);
    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <main className="bg-[#080810] text-white overflow-x-hidden font-sans selection:bg-[#7c6dfa]/30 w-full relative">
      {/* NAV */}
      <nav
        className="fixed top-0 left-0 right-0 z-50 transition-all duration-300 w-full"
        style={{
          background:
            scrollY > 40
              ? "rgba(8,8,16,0.75)"
              : "transparent",
          backdropFilter: scrollY > 40 ? "blur(16px)" : "none",
          borderBottom:
            scrollY > 40 ? "1px solid rgba(255,255,255,0.08)" : "1px solid transparent",
        }}
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 py-4 flex items-center justify-between">
          <span className="font-bold text-lg tracking-tight flex items-center gap-2">
            <Activity className="w-5 h-5 text-[#7c6dfa]" />
            <span className="hidden sm:inline">Noitrex</span>
          </span>
          <div className="hidden md:flex items-center gap-8 text-sm text-white/60 font-medium">
            <a href="#features" className="hover:text-white transition-colors">Features</a>
            <a href="#how-it-works" className="hover:text-white transition-colors">Architecture</a>
            <a href="#pricing-models" className="hover:text-white transition-colors">Pricing</a>
            <a href="#docs" className="hover:text-white transition-colors">Docs</a>
          </div>
          <a
            href="#get-started"
            className="text-xs sm:text-sm px-4 sm:px-5 py-2 sm:py-2.5 rounded-full border border-white/10 bg-white/5 hover:bg-white/10 transition-all font-medium flex items-center gap-2"
          >
            Get Started <ArrowRight className="w-3 h-3 sm:w-4 sm:h-4" />
          </a>
        </div>
      </nav>

      {/* HERO */}
      <section
        ref={heroRef}
        className="relative min-h-[100svh] flex flex-col items-center justify-center text-center px-4 sm:px-6 pt-32 pb-20 w-full overflow-hidden"
      >
        {/* background grid */}
        <div
          className="absolute inset-0 opacity-[0.04]"
          style={{
            backgroundImage:
              "linear-gradient(rgba(255,255,255,0.8) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,0.8) 1px, transparent 1px)",
            backgroundSize: "40px 40px",
          }}
        />

        {/* purple glow orbs */}
        <div
          className="absolute top-1/3 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[300px] h-[300px] sm:w-[600px] sm:h-[600px] rounded-full pointer-events-none"
          style={{
            background:
              "radial-gradient(circle, rgba(124,109,250,0.15) 0%, transparent 60%)",
            transform: `translate(-50%, calc(-50% + ${scrollY * 0.15}px))`,
          }}
        />
        <div
          className="absolute top-2/3 left-1/4 w-[200px] h-[200px] sm:w-[500px] sm:h-[500px] rounded-full pointer-events-none"
          style={{
            background:
              "radial-gradient(circle, rgba(56,189,248,0.08) 0%, transparent 60%)",
          }}
        />

        {/* pill badge */}
        <div className="inline-flex items-center gap-2 px-3 py-1.5 sm:px-4 sm:py-2 rounded-full border border-[#7c6dfa]/30 bg-[#7c6dfa]/10 text-[10px] sm:text-xs font-medium text-[#a89dfc] mb-6 sm:mb-8 animate-fade-in backdrop-blur-sm z-10">
          <Server className="w-3 h-3 sm:w-3.5 sm:h-3.5" />
          Self-hostable · Usage-based · Built in Go
        </div>

        <h1 className="text-4xl sm:text-6xl md:text-8xl font-extrabold tracking-tight leading-[1.1] max-w-5xl mb-6 sm:mb-8 animate-fade-in-up z-10">
          Stop building<br className="sm:hidden" />{" "}
          <span className="relative inline-block mt-2 sm:mt-0">
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#7c6dfa] via-[#a89dfc] to-[#38bdf8]">
              billing infra.
            </span>
            <span
              className="absolute inset-0 blur-2xl sm:blur-3xl opacity-40 text-transparent bg-clip-text bg-gradient-to-r from-[#7c6dfa] to-[#38bdf8]"
              aria-hidden
            >
              billing infra.
            </span>
          </span>
          <br />
          Start building.
        </h1>

        <p className="text-sm sm:text-lg md:text-xl text-white/50 max-w-2xl mb-10 sm:mb-12 leading-relaxed animate-fade-in-up animation-delay-100 px-2 sm:px-0 z-10">
          Noitrex is a high-performance billing engine that counts usage, computes
          invoices with tiered pricing, and fires signed webhooks — without
          touching your payment provider.
        </p>

        <div className="flex flex-col w-full sm:w-auto sm:flex-row items-center gap-3 sm:gap-4 animate-fade-in-up animation-delay-200 z-10">
          <a
            id="get-started"
            href="https://github.com"
            className="w-full sm:w-auto flex items-center justify-center gap-2 px-6 sm:px-8 py-3.5 sm:py-4 rounded-full bg-gradient-to-r from-[#7c6dfa] to-[#5b4fcf] hover:from-[#8f81fc] hover:to-[#6d5fe0] text-white font-semibold text-sm transition-all shadow-[0_0_30px_-10px_rgba(124,109,250,0.5)] hover:shadow-[0_0_50px_-10px_rgba(124,109,250,0.6)] hover:-translate-y-0.5"
          >
            Deploy for free
            <ArrowRight className="w-4 h-4 transition-transform group-hover:translate-x-1" />
          </a>
          <a
            href="#how-it-works"
            className="w-full sm:w-auto flex items-center justify-center px-6 sm:px-8 py-3.5 sm:py-4 rounded-full border border-white/10 text-white/70 hover:text-white hover:border-white/20 hover:bg-white/5 text-sm font-medium transition-all"
          >
            See architecture
          </a>
        </div>

        {/* terminal snippet */}
        <div className="mt-16 sm:mt-20 w-full max-w-3xl mx-auto rounded-xl border border-white/10 bg-[#0c0c14]/80 backdrop-blur-xl overflow-hidden text-left shadow-2xl animate-fade-in-up animation-delay-300 z-10">
          <div className="flex items-center justify-between px-3 sm:px-4 py-2.5 sm:py-3 border-b border-white/10 bg-white/[0.02]">
            <div className="flex items-center gap-2">
              <div className="flex items-center gap-1.5 mr-2 sm:mr-4">
                <span className="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full bg-red-500/80" />
                <span className="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full bg-yellow-500/80" />
                <span className="w-2.5 h-2.5 sm:w-3 sm:h-3 rounded-full bg-green-500/80" />
              </div>
              <Terminal className="w-3 h-3 sm:w-4 sm:h-4 text-white/40 hidden sm:block" />
              <span className="text-[10px] sm:text-xs text-white/40 font-mono font-medium truncate max-w-[120px] sm:max-w-none">Ingest Event</span>
            </div>
            <span className="text-[10px] sm:text-xs text-[#4ade80] font-mono bg-[#4ade80]/10 px-1.5 sm:px-2 py-0.5 sm:py-1 rounded whitespace-nowrap">200 OK &lt; 5ms</span>
          </div>
          <div className="p-4 sm:p-5 overflow-x-auto">
            <pre className="text-xs sm:text-sm font-mono text-white/70 leading-relaxed whitespace-pre sm:whitespace-pre-wrap word-break">
              <span className="text-[#7c6dfa]">curl</span>{" "}
              <span className="text-[#38bdf8]">-X POST</span>{" "}
              https://api.noitrex.com/v1/events/ingest \{"\n"}
              {"  "}-H{" "}
              <span className="text-[#a89dfc]">&quot;Authorization: Bearer sk_live_...&quot;</span>{" "}\{"\n"}
              {"  "}-d{" "}
              <span className="text-white/90">&#39;&#123;</span>{"\n"}
              {"    "}
              <span className="text-[#4ade80]">&quot;customer_id&quot;</span>:{" "}
              <span className="text-amber-300">&quot;cus_abc123&quot;</span>,{"\n"}
              {"    "}
              <span className="text-[#4ade80]">&quot;event_name&quot;</span>:{" "}
              <span className="text-amber-300">&quot;api.request&quot;</span>,{"\n"}
              {"    "}
              <span className="text-[#4ade80]">&quot;quantity&quot;</span>:{" "}
              <span className="text-sky-300">1</span>,{"\n"}
              {"    "}
              <span className="text-[#4ade80]">&quot;idempotency_key&quot;</span>:{" "}
              <span className="text-amber-300">&quot;req_xyz789&quot;</span>{"\n"}
              {"  "}
              <span className="text-white/90">&#125;&#39;</span>
            </pre>
          </div>
        </div>
      </section>

      {/* STATS BAR */}
      <section className="border-y border-white/10 bg-white/[0.01] py-8 sm:py-12 relative z-10 w-full overflow-hidden">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 grid grid-cols-2 md:grid-cols-4 gap-y-8 sm:gap-8 text-center divide-x-0 md:divide-x divide-white/10">
          {[
            { value: "< 5ms", label: "Event ingestion p99" },
            { value: "Exactly-once", label: "Idempotency guarantee" },
            { value: "Integer kobo", label: "No float arithmetic" },
            { value: "HMAC-SHA256", label: "Signed webhooks" },
          ].map((stat, i) => (
            <div key={stat.label} className={i % 2 === 0 ? "border-l-0" : "md:border-l border-white/10 border-l-0"}>
              <div className="text-xl sm:text-2xl md:text-3xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-[#7c6dfa] to-[#38bdf8] mb-1 sm:mb-2 tracking-tight">
                {stat.value}
              </div>
              <div className="text-[10px] sm:text-xs md:text-sm text-white/40 font-medium">{stat.label}</div>
            </div>
          ))}
        </div>
      </section>

      {/* FEATURES */}
      <section id="features" className="py-20 sm:py-32 px-4 sm:px-6 relative z-10">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-12 sm:mb-20">
            <div className="inline-flex items-center justify-center px-3 py-1 rounded-full border border-white/10 bg-white/5 text-[10px] sm:text-xs uppercase tracking-widest text-[#7c6dfa] mb-4 sm:mb-6 font-semibold">
              <Code2 className="w-3 h-3 sm:w-3.5 sm:h-3.5 mr-2" />
              Core Capabilities
            </div>
            <h2 className="text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight px-2">
              The billing layer your API{" "}
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#7c6dfa] to-[#38bdf8]">
                deserves
              </span>
            </h2>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 lg:gap-8">
            {[
              {
                icon: Zap,
                title: "Real-time event ingestion",
                desc: "POST usage events and get a 200 back in microseconds. Noitrex persists, publishes internally via NexusMQ, and increments your Redis counters — all atomically.",
                accent: "text-[#7c6dfa]",
                bg: "bg-[#7c6dfa]/10",
                border: "border-[#7c6dfa]/20"
              },
              {
                icon: Calculator,
                title: "Tiered pricing",
                desc: "Flat rate, per-unit, or volume tiers — computed correctly every time. All arithmetic uses BIGINT kobo. Floats are strictly prohibited.",
                accent: "text-[#38bdf8]",
                bg: "bg-[#38bdf8]/10",
                border: "border-[#38bdf8]/20"
              },
              {
                icon: Mail,
                title: "Signed webhooks",
                desc: "Invoice creation triggers a NexusMQ event. The dispatcher delivers a signed HMAC-SHA256 payload to your registered endpoint via NexusRelay.",
                accent: "text-[#a78bfa]",
                bg: "bg-[#a78bfa]/10",
                border: "border-[#a78bfa]/20"
              },
              {
                icon: ShieldCheck,
                title: "Operator isolation",
                desc: "Every row in every table is scoped by operator_id. No cross-tenant data leaks possible — enforced at the PostgreSQL query level.",
                accent: "text-[#34d399]",
                bg: "bg-[#34d399]/10",
                border: "border-[#34d399]/20"
              },
              {
                icon: InfinityIcon,
                title: "Exactly-once counting",
                desc: "Idempotency keys have a UNIQUE constraint at the schema level. Duplicate events are rejected natively at the INSERT step.",
                accent: "text-[#fb923c]",
                bg: "bg-[#fb923c]/10",
                border: "border-[#fb923c]/20"
              },
              {
                icon: Database,
                title: "Aggregate flushing",
                desc: "Redis counters flush to PostgreSQL on a 30-second interval using atomic INSERT ... ON CONFLICT DO UPDATE. Always consistent.",
                accent: "text-[#f472b6]",
                bg: "bg-[#f472b6]/10",
                border: "border-[#f472b6]/20"
              },
            ].map((f) => (
              <div
                key={f.title}
                className="group relative p-6 sm:p-8 rounded-[2rem] border border-white/10 bg-white/[0.02] hover:bg-white/[0.04] transition-all duration-500 overflow-hidden"
              >
                <div className={`absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-500 bg-gradient-to-br from-transparent to-white/[0.02]`} />
                
                <div
                  className={`w-12 h-12 sm:w-14 sm:h-14 rounded-2xl ${f.bg} ${f.border} border flex items-center justify-center mb-5 sm:mb-6 shadow-inner`}
                >
                  <f.icon className={`w-6 h-6 sm:w-7 sm:h-7 ${f.accent}`} />
                </div>
                <h3 className="font-bold text-lg sm:text-xl text-white mb-2 sm:mb-3 tracking-tight">{f.title}</h3>
                <p className="text-white/50 leading-relaxed text-xs sm:text-sm">{f.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* THREE-PARTY MODEL */}
      <section className="py-20 sm:py-24 px-4 sm:px-6 border-y border-white/10 bg-white/[0.01]">
        <div className="max-w-6xl mx-auto text-center">
          <div className="inline-flex items-center justify-center px-3 py-1 rounded-full border border-white/10 bg-white/5 text-[10px] sm:text-xs uppercase tracking-widest text-[#7c6dfa] mb-4 sm:mb-6 font-semibold">
            Architecture
          </div>
          <h2 className="text-3xl sm:text-4xl md:text-5xl font-bold mb-4 sm:mb-6 tracking-tight">The Three-Party Model</h2>
          <p className="text-white/50 max-w-2xl mx-auto mb-10 sm:mb-16 text-sm sm:text-lg leading-relaxed px-4 sm:px-0">
            Noitrex sits cleanly in the middle of the billing relationship between you and your customers.
          </p>

          <div className="flex flex-col md:flex-row items-stretch justify-center gap-4 sm:gap-6 lg:gap-8 px-2 sm:px-0">
            {[
              { label: "Operator", desc: "You — the SaaS API company. Send usage events, receive webhooks.", color: "#7c6dfa", border: "border-[#7c6dfa]/30", icon: Terminal },
              { label: "Noitrex", desc: "The engine. Counts usage, computes invoices, fires webhooks.", color: "#38bdf8", border: "border-[#38bdf8]/40", bg: "bg-[#38bdf8]/5", isCenter: true, icon: Activity },
              { label: "Customer", desc: "Your end-users. Billed transparently based on actual consumption.", color: "#4ade80", border: "border-[#4ade80]/30", icon: Zap },
            ].map((party, i) => (
              <div key={party.label} className="flex flex-col md:flex-row items-center flex-1 w-full">
                {i > 0 && (
                  <div className="py-2 sm:py-4 md:py-0 md:px-2 flex-shrink-0">
                    <ArrowLeftRight className="w-5 h-5 sm:w-6 sm:h-6 text-white/20 hidden md:block" />
                    <ArrowDown className="w-5 h-5 sm:w-6 sm:h-6 text-white/20 block md:hidden" />
                  </div>
                )}
                <div
                  className={`flex-1 w-full rounded-[2rem] border p-6 sm:p-8 text-center transition-all ${
                    party.isCenter
                      ? `${party.border} ${party.bg} shadow-[0_0_30px_-10px_rgba(56,189,248,0.2)] md:-translate-y-4`
                      : "border-white/10 bg-white/[0.02]"
                  }`}
                >
                  <div
                    className="w-12 h-12 sm:w-16 sm:h-16 rounded-2xl mx-auto mb-4 sm:mb-6 flex items-center justify-center shadow-inner"
                    style={{ background: party.color + "15", border: `1px solid ${party.color}40` }}
                  >
                    <party.icon className="w-6 h-6 sm:w-8 sm:h-8" style={{ color: party.color }} />
                  </div>
                  <h3 className="font-bold text-xl sm:text-2xl mb-2 sm:mb-3" style={{ color: party.color }}>{party.label}</h3>
                  <p className="text-white/60 leading-relaxed text-xs sm:text-sm">{party.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* HOW IT WORKS (DATA FLOW) */}
      <section id="how-it-works" className="py-20 sm:py-32 px-4 sm:px-6 relative overflow-hidden">
        <div
          className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[400px] h-[400px] sm:w-[800px] sm:h-[600px] pointer-events-none"
          style={{ background: "radial-gradient(ellipse, rgba(124,109,250,0.08) 0%, transparent 60%)" }}
        />
        <div className="max-w-4xl mx-auto relative z-10">
          <div className="text-center mb-16 sm:mb-20">
            <div className="inline-flex items-center justify-center px-3 py-1 rounded-full border border-white/10 bg-white/5 text-[10px] sm:text-xs uppercase tracking-widest text-[#7c6dfa] mb-4 sm:mb-6 font-semibold">
              Data Flow
            </div>
            <h2 className="text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight">How Noitrex Works</h2>
          </div>

          <div className="relative pl-0 sm:pl-10 md:pl-20 mx-auto">
            {/* connecting line - hidden on very small mobile, visible on sm and up */}
            <div className="absolute left-[20px] sm:left-[31px] md:left-[71px] top-[40px] bottom-[40px] w-px bg-gradient-to-b from-[#7c6dfa] via-[#38bdf8] to-transparent hidden sm:block opacity-30" />

            <div className="space-y-10 sm:space-y-16">
              {[
                {
                  step: "01",
                  title: "Event Ingestion",
                  desc: "The Operator calls POST /v1/events/ingest. Noitrex persists the raw event to PostgreSQL, publishes internally via NexusMQ, and returns HTTP 200 immediately.",
                  tag: "< 5ms",
                },
                {
                  step: "02",
                  title: "Fast Counting",
                  desc: "The metering engine subscribes to NexusMQ events and atomically increments a Redis counter scoped to the customer and billing period.",
                  tag: "Redis atomic INCR",
                },
                {
                  step: "03",
                  title: "Aggregate Flushing",
                  desc: "Every 30 seconds, Redis counters flush to PostgreSQL aggregates using INSERT ... ON CONFLICT DO UPDATE — preventing race conditions.",
                  tag: "30s interval",
                },
                {
                  step: "04",
                  title: "Invoice Generation",
                  desc: "At period end, the billing worker reads aggregates, applies your pricing plan, computes the invoice amount in integer kobo, and writes an immutable record.",
                  tag: "BIGINT kobo only",
                },
                {
                  step: "05",
                  title: "Webhook Dispatch",
                  desc: "Invoice creation fires a NexusMQ event. The dispatcher signs the payload with HMAC-SHA256 and delivers it to your endpoint via NexusRelay.",
                  tag: "HMAC-SHA256",
                },
              ].map((item, i) => (
                <div key={item.step} className="relative flex gap-4 sm:gap-6 md:gap-8 items-start group">
                  <div className="relative flex-shrink-0 z-10 pt-1 sm:pt-0">
                    <div className="w-10 h-10 sm:w-16 sm:h-16 rounded-full border border-[#7c6dfa]/30 sm:border-2 bg-[#080810] flex items-center justify-center text-[#7c6dfa] font-bold text-xs sm:text-base font-mono group-hover:border-[#7c6dfa] group-hover:shadow-[0_0_20px_-5px_rgba(124,109,250,0.5)] transition-all">
                      {item.step}
                    </div>
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4 mb-2 sm:mb-3">
                      <h3 className="font-bold text-white text-lg sm:text-2xl tracking-tight truncate">{item.title}</h3>
                      <span className="text-[10px] sm:text-xs px-2 sm:px-3 py-0.5 sm:py-1 rounded-full bg-[#38bdf8]/10 text-[#38bdf8] border border-[#38bdf8]/20 font-mono w-fit">
                        {item.tag}
                      </span>
                    </div>
                    <p className="text-white/50 text-xs sm:text-base leading-relaxed">{item.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* PRICING MODELS */}
      <section id="pricing-models" className="py-20 sm:py-32 px-4 sm:px-6 bg-white/[0.01]">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16 sm:mb-20">
            <div className="inline-flex items-center justify-center px-3 py-1 rounded-full border border-white/10 bg-white/5 text-[10px] sm:text-xs uppercase tracking-widest text-[#7c6dfa] mb-4 sm:mb-6 font-semibold">
              Billing Logic
            </div>
            <h2 className="text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight">Every pricing model, built in</h2>
          </div>

          <div className="grid md:grid-cols-3 gap-4 sm:gap-6 lg:gap-8">
            {[
              {
                name: "Flat Rate",
                icon: Square,
                color: "#7c6dfa",
                bg: "bg-[#7c6dfa]/10",
                border: "border-[#7c6dfa]/30",
                desc: "Fixed price per billing period regardless of volume. Simple, predictable, zero computation overhead.",
                example: "$99/month — always.",
              },
              {
                name: "Per Unit",
                icon: X,
                color: "#38bdf8",
                bg: "bg-[#38bdf8]/10",
                border: "border-[#38bdf8]/30",
                desc: "Linear pricing that scales directly with consumption. Perfect for API calls, tokens, or seats.",
                example: "1k calls × $0.002 = $2.00",
              },
              {
                name: "Tiered",
                icon: Layers,
                color: "#a78bfa",
                bg: "bg-[#a78bfa]/10",
                border: "border-[#a78bfa]/30",
                desc: "Usage walks through price buckets sequentially. Each tier prices only the units within its specific range.",
                example: "0-1k: free · 10k+: $0.0005",
              },
            ].map((model) => (
              <div
                key={model.name}
                className="relative p-6 sm:p-8 rounded-[2rem] border border-white/10 bg-white/[0.02] overflow-hidden group hover:bg-white/[0.04] transition-all"
              >
                <div
                  className="absolute -top-12 -right-12 w-32 h-32 sm:w-40 sm:h-40 rounded-full opacity-10 blur-2xl sm:blur-3xl group-hover:opacity-20 transition-opacity"
                  style={{ background: model.color }}
                />
                <div className={`w-12 h-12 sm:w-14 sm:h-14 rounded-2xl ${model.bg} ${model.border} border flex items-center justify-center mb-5 sm:mb-6`}>
                  <model.icon className="w-6 h-6 sm:w-7 sm:h-7" style={{ color: model.color }} />
                </div>
                <h3 className="font-bold text-xl sm:text-2xl mb-2 sm:mb-3 tracking-tight">{model.name}</h3>
                <p className="text-white/50 leading-relaxed mb-6 sm:mb-8 text-xs sm:text-sm">{model.desc}</p>
                <div
                  className="text-[10px] sm:text-xs font-mono px-3 sm:px-4 py-2 sm:py-3 rounded-xl bg-[#080810] border border-white/10 shadow-inner truncate"
                  style={{ color: model.color }}
                >
                  {model.example}
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* WHAT IT DOESN'T DO */}
      <section className="py-20 sm:py-24 px-4 sm:px-6 border-t border-white/10">
        <div className="max-w-5xl mx-auto text-center">
          <div className="inline-flex items-center justify-center px-3 py-1 rounded-full border border-white/10 bg-white/5 text-[10px] sm:text-xs uppercase tracking-widest text-white/40 mb-4 sm:mb-6 font-semibold">
            Strict Boundaries
          </div>
          <h2 className="text-2xl sm:text-4xl md:text-5xl font-bold mb-8 sm:mb-12 tracking-tight">
            What Noitrex{" "}
            <span className="text-white/30">does not do</span>
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 sm:gap-6">
            {[
              { title: "No card charging", desc: "Noitrex computes what is owed. You collect it via Stripe or your provider." },
              { title: "No external queues", desc: "NexusMQ is strictly in-process for v1. No Kafka, no Redis Streams." },
              { title: "No float math", desc: "All money is stored as BIGINT kobo. Floating point is architecturally prohibited." },
            ].map((item) => (
              <div key={item.title} className="p-5 sm:p-8 rounded-[1.5rem] sm:rounded-[2rem] border border-white/10 bg-white/[0.02] text-left">
                <XCircle className="w-6 h-6 sm:w-8 sm:h-8 text-red-400/80 mb-3 sm:mb-4" />
                <h3 className="font-bold text-base sm:text-lg mb-2 text-white/90">{item.title}</h3>
                <p className="text-white/50 leading-relaxed text-xs sm:text-sm">{item.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-24 sm:py-32 px-4 sm:px-6 relative overflow-hidden">
        <div
          className="absolute inset-0 pointer-events-none"
          style={{ background: "radial-gradient(ellipse at 50% 50%, rgba(124,109,250,0.15) 0%, transparent 60%)" }}
        />
        <div className="max-w-4xl mx-auto text-center relative z-10">
          <h2 className="text-4xl sm:text-5xl md:text-7xl font-extrabold tracking-tight mb-6 sm:mb-8">
            Ready to ship<br className="sm:hidden" />{" "}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#7c6dfa] to-[#38bdf8]">
              faster?
            </span>
          </h2>
          <p className="text-white/50 text-base sm:text-xl mb-10 sm:mb-12 max-w-2xl mx-auto leading-relaxed px-4">
            Self-host Noitrex in minutes. Focus on your product — let Noitrex handle the billing infrastructure.
          </p>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-6 px-4 sm:px-0 w-full">
            <a
              href="https://github.com"
              className="w-full sm:w-auto group flex items-center justify-center gap-2 px-8 py-3.5 sm:py-4 rounded-full bg-gradient-to-r from-[#7c6dfa] to-[#5b4fcf] hover:from-[#8f81fc] hover:to-[#6d5fe0] text-white font-semibold text-sm sm:text-base transition-all shadow-[0_0_40px_-10px_rgba(124,109,250,0.5)] hover:shadow-[0_0_60px_-15px_rgba(124,109,250,0.6)] hover:-translate-y-0.5"
            >
              View on GitHub
              <ArrowRight className="w-4 h-4 sm:w-5 sm:h-5 group-hover:translate-x-1 transition-transform" />
            </a>
            <a
              href="#docs"
              className="w-full sm:w-auto flex items-center justify-center px-8 py-3.5 sm:py-4 rounded-full border border-white/10 text-white/70 hover:text-white hover:border-white/20 hover:bg-white/5 text-sm sm:text-base font-medium transition-all"
            >
              Read the docs
            </a>
          </div>

          {/* bottom tagline */}
          <div className="mt-16 sm:mt-20 flex flex-col sm:flex-row items-center justify-center gap-2 text-[10px] sm:text-xs text-white/30 font-mono">
            <Activity className="w-3 h-3 sm:w-4 sm:h-4 hidden sm:block" />
            <span>Noitrex: Stop building billing infra.</span>
          </div>
        </div>
      </section>

      {/* FOOTER */}
      <footer className="border-t border-white/10 bg-[#080810] py-8 sm:py-10 px-4 sm:px-6">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-6 text-center md:text-left">
          <div className="flex items-center gap-2 sm:gap-3">
            <div className="w-6 h-6 sm:w-8 sm:h-8 rounded-lg bg-[#7c6dfa]/10 border border-[#7c6dfa]/30 flex items-center justify-center">
              <Activity className="w-3 h-3 sm:w-4 sm:h-4 text-[#7c6dfa]" />
            </div>
            <span className="font-semibold tracking-tight text-sm sm:text-base">Noitrex</span>
          </div>
          <div className="flex flex-wrap items-center justify-center gap-4 sm:gap-8 text-xs sm:text-sm text-white/50 font-medium">
            <a href="#" className="hover:text-white transition-colors">GitHub</a>
            <a href="#" className="hover:text-white transition-colors">Documentation</a>
            <a href="#" className="hover:text-white transition-colors">Twitter</a>
            <a href="#" className="hover:text-white transition-colors">License</a>
          </div>
          <div className="text-[10px] sm:text-xs text-white/30 font-mono">
            © {new Date().getFullYear()} Noitrex
          </div>
        </div>
      </footer>

      <style>{`
        @keyframes fade-in {
          from { opacity: 0; }
          to { opacity: 1; }
        }
        @keyframes fade-in-up {
          from { opacity: 0; transform: translateY(20px); }
          to { opacity: 1; transform: translateY(0); }
        }
        .animate-fade-in {
          animation: fade-in 0.6s cubic-bezier(0.16, 1, 0.3, 1) both;
        }
        .animate-fade-in-up {
          animation: fade-in-up 0.6s cubic-bezier(0.16, 1, 0.3, 1) both;
        }
        .animation-delay-100 { animation-delay: 0.1s; }
        .animation-delay-200 { animation-delay: 0.2s; }
        .animation-delay-300 { animation-delay: 0.3s; }
      `}</style>
    </main>
  );
}
