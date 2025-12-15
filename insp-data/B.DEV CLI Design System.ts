// ═══════════════════════════════════════════════════════════════════════════════
// B.DEV CLI DESIGN SYSTEM v1.0
// Military-Grade Terminal Interface Components
// ═══════════════════════════════════════════════════════════════════════════════

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 1. ASCII ART HEADERS (Compact, One-Line Format)                             │
// └─────────────────────────────────────────────────────────────────────────────┘

const HEADERS = {
  BUILD: `╔═══════════════════════════════════════╗  ██████╗ ██╗   ██╗██╗██╗     ██████╗  ╔═══════════════════════════════════════╗`,
  
  SECURITY: `╔═[SECURITY]═══════════════════════════╗  ███████╗███████╗ ██████╗██╗   ██╗██████╗ ██╗████████╗██╗   ██╗  ╔═══════╗`,
  
  NETWORK: `╔═══════════════════════════════════════╗  ███╗   ██╗███████╗████████╗██╗    ██╗ ██████╗ ██████╗ ██╗  ██╗  ╔═════╗`,
  
  TEST: `╔═══════════════════════════════════════╗  ████████╗███████╗███████╗████████╗  ╔═══════════════════════════════════╗`,
  
  DEPLOY: `╔═══════════════════════════════════════╗  ██████╗ ███████╗██████╗ ██╗      ██████╗ ██╗   ██╗  ╔═════════════════╗`,
  
  SYSTEM: `╔═══════════════════════════════════════╗  ███████╗██╗   ██╗███████╗████████╗███████╗███╗   ███╗  ╔═════════════╗`
};

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 2. PROGRESS BARS (6+ Sophistiqués avec Gradients & Animations)             │
// └─────────────────────────────────────────────────────────────────────────────┘

// Progress Bar Type 1: BUILD COMPILATION (Gradient Block Style)
const PROGRESS_BUILD = {
  frames: [
    `[▰▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱] 03% ┃ Initializing compiler...`,
    `[▰▰▰▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱] 12% ┃ Parsing TypeScript modules...`,
    `[▰▰▰▰▰▰▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱] 24% ┃ Type checking definitions...`,
    `[▰▰▰▰▰▰▰▰▰▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱] 38% ┃ Transpiling to JavaScript...`,
    `[▰▰▰▰▰▰▰▰▰▰▰▰▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱] 51% ┃ Bundling dependencies...`,
    `[▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱] 63% ┃ Optimizing chunks...`,
    `[▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▱▱▱▱▱▱▱▱▱▱▱▱] 77% ┃ Minifying output...`,
    `[▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▱▱▱▱▱▱▱▱] 89% ┃ Generating source maps...`,
    `[▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰] 100% ┃ Build complete. Ready.`
  ]
};

// Progress Bar Type 2: TEST EXECUTION (Detailed Counter Style)
const PROGRESS_TEST = {
  frames: [
    `[░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 000/847 ┃ Starting test suite...`,
    `[▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 042/847 ┃ auth.service.spec.ts`,
    `[▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 127/847 ┃ user.controller.spec.ts`,
    `[▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░] 234/847 ┃ database.integration.spec.ts`,
    `[▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░] 391/847 ┃ api.endpoints.spec.ts`,
    `[▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░] 518/847 ┃ middleware.chain.spec.ts`,
    `[▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░] 623/847 ┃ validation.pipes.spec.ts`,
    `[▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░] 741/847 ┃ guards.security.spec.ts`,
    `[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓] 847/847 ┃ All tests passed. Coverage: 94.2%`
  ]
};

// Progress Bar Type 3: NETWORK TRANSFER (Speed Indicator Style)
const PROGRESS_NETWORK = {
  frames: [
    `[━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━] 000 MB/512 MB ┃ 0.0 MB/s ┃ Initializing...`,
    `[█▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒] 032 MB/512 MB ┃ 8.4 MB/s ┃ Establishing connection...`,
    `[██▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒] 089 MB/512 MB ┃ 12.7 MB/s ┃ Downloading chunks...`,
    `[████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒] 156 MB/512 MB ┃ 15.3 MB/s ┃ Buffering stream...`,
    `[██████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒] 234 MB/512 MB ┃ 18.9 MB/s ┃ Receiving payload...`,
    `[████████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒] 318 MB/512 MB ┃ 21.2 MB/s ┃ Decompressing data...`,
    `[██████████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒] 401 MB/512 MB ┃ 19.8 MB/s ┃ Writing to disk...`,
    `[████████████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒] 467 MB/512 MB ┃ 17.1 MB/s ┃ Finalizing transfer...`,
    `[██████████████████████████████] 512 MB/512 MB ┃ Avg: 16.4 MB/s ┃ Complete.`
  ]
};

// Progress Bar Type 4: SECURITY SCAN (Analysis Progress)
const PROGRESS_SECURITY = {
  frames: [
    `[◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐] 00% ┃ Initializing ISO 27001 audit...`,
    `[◓◓◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐] 08% ┃ Scanning dependencies for CVEs...`,
    `[◓◓◓◓◓◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐] 18% ┃ Checking encryption protocols...`,
    `[◓◓◓◓◓◓◓◓◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐] 29% ┃ Auditing authentication flows...`,
    `[◓◓◓◓◓◓◓◓◓◓◓◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐] 41% ┃ Analyzing access control lists...`,
    `[◓◓◓◓◓◓◓◓◓◓◓◓◓◓◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐] 53% ┃ Validating data sanitization...`,
    `[◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◐◐◐◐◐◐◐◐◐◐◐◐◐] 64% ┃ Testing SQL injection vectors...`,
    `[◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◐◐◐◐◐◐◐◐◐] 77% ┃ Reviewing security headers...`,
    `[◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◐◐◐◐◐] 88% ┃ Generating compliance report...`,
    `[◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓] 100% ┃ Security audit complete. 0 critical issues.`
  ]
};

// Progress Bar Type 5: BUILD OPTIMIZATION (Multi-Phase Style)
const PROGRESS_OPTIMIZE = {
  frames: [
    `[Phase 1/4] [▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁] 00% ┃ Tree shaking unused code...`,
    `[Phase 1/4] [▃▃▃▃▃▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁] 27% ┃ Eliminating dead code paths...`,
    `[Phase 1/4] [▃▃▃▃▃▃▃▃▃▃▁▁▁▁▁▁▁▁▁▁] 51% ┃ Removing debug statements...`,
    `[Phase 2/4] [▅▅▅▅▅▅▅▅▅▅▅▅▅▃▃▁▁▁▁▁] 68% ┃ Compressing assets...`,
    `[Phase 2/4] [▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅▅] 100% ┃ Asset compression complete.`,
    `[Phase 3/4] [▆▆▆▆▆▆▁▁▁▁▁▁▁▁▁▁▁▁▁▁] 34% ┃ Code splitting chunks...`,
    `[Phase 3/4] [▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆▆] 100% ┃ Chunk optimization done.`,
    `[Phase 4/4] [▇▇▇▇▇▇▇▇▇▇▁▁▁▁▁▁▁▁▁▁] 53% ┃ Applying gzip compression...`,
    `[Phase 4/4] [▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇▇] 100% ┃ Build optimized. Size: 2.4 MB → 847 KB`
  ]
};

// Progress Bar Type 6: DEPLOYMENT (Multi-Stage Pipeline)
const PROGRESS_DEPLOY = {
  frames: [
    `[Stage 1/5] [◯◯◯◯◯◯◯◯◯◯◯◯◯◯◯◯◯◯◯◯] 00% ┃ Provisioning infrastructure...`,
    `[Stage 1/5] [●●●●●◯◯◯◯◯◯◯◯◯◯◯◯◯◯◯] 23% ┃ Spinning up containers...`,
    `[Stage 2/5] [●●●●●●●●●◯◯◯◯◯◯◯◯◯◯◯] 41% ┃ Configuring load balancer...`,
    `[Stage 3/5] [●●●●●●●●●●●●◯◯◯◯◯◯◯◯] 58% ┃ Deploying application code...`,
    `[Stage 3/5] [●●●●●●●●●●●●●●●◯◯◯◯◯] 72% ┃ Running database migrations...`,
    `[Stage 4/5] [●●●●●●●●●●●●●●●●●◯◯◯] 84% ┃ Warming up cache layers...`,
    `[Stage 4/5] [●●●●●●●●●●●●●●●●●●●◯] 93% ┃ Health checks in progress...`,
    `[Stage 5/5] [●●●●●●●●●●●●●●●●●●●●] 100% ┃ Deployment successful. Live on production.`
  ]
};

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 3. STATUS INDICATORS (10+ Types avec Animations Fluides)                   │
// └─────────────────────────────────────────────────────────────────────────────┘

const STATUS_LOADING = {
  frames: [
    `⠋ LOADING ┃ Fetching resources...`,
    `⠙ LOADING ┃ Fetching resources...`,
    `⠹ LOADING ┃ Fetching resources...`,
    `⠸ LOADING ┃ Fetching resources...`,
    `⠼ LOADING ┃ Fetching resources...`,
    `⠴ LOADING ┃ Fetching resources...`,
    `⠦ LOADING ┃ Fetching resources...`,
    `⠧ LOADING ┃ Fetching resources...`,
    `⠇ LOADING ┃ Fetching resources...`,
    `⠏ LOADING ┃ Fetching resources...`
  ]
};

const STATUS_PROCESSING = {
  frames: [
    `◐ PROCESSING ┃ Executing pipeline stage 1/4...`,
    `◓ PROCESSING ┃ Executing pipeline stage 1/4...`,
    `◑ PROCESSING ┃ Executing pipeline stage 2/4...`,
    `◒ PROCESSING ┃ Executing pipeline stage 2/4...`,
    `◐ PROCESSING ┃ Executing pipeline stage 3/4...`,
    `◓ PROCESSING ┃ Executing pipeline stage 3/4...`,
    `◑ PROCESSING ┃ Executing pipeline stage 4/4...`,
    `◒ PROCESSING ┃ Executing pipeline stage 4/4...`
  ]
};

const STATUS_SCANNING = {
  frames: [
    `▱▱▱▱▱▱▱▱▱▱ SCANNING ┃ Analyzing codebase structure...`,
    `▰▱▱▱▱▱▱▱▱▱ SCANNING ┃ Analyzing codebase structure...`,
    `▰▰▱▱▱▱▱▱▱▱ SCANNING ┃ Indexing file dependencies...`,
    `▰▰▰▱▱▱▱▱▱▱ SCANNING ┃ Indexing file dependencies...`,
    `▰▰▰▰▱▱▱▱▱▱ SCANNING ┃ Mapping import graphs...`,
    `▰▰▰▰▰▱▱▱▱▱ SCANNING ┃ Mapping import graphs...`,
    `▰▰▰▰▰▰▱▱▱▱ SCANNING ┃ Detecting circular references...`,
    `▰▰▰▰▰▰▰▱▱▱ SCANNING ┃ Detecting circular references...`,
    `▰▰▰▰▰▰▰▰▱▱ SCANNING ┃ Building syntax tree...`,
    `▰▰▰▰▰▰▰▰▰▱ SCANNING ┃ Building syntax tree...`,
    `▰▰▰▰▰▰▰▰▰▰ SCANNING ┃ Finalizing analysis results...`
  ]
};

const STATUS_BUILDING = {
  frames: [
    `┤ BUILDING ┃ Compiling source files...`,
    `┴ BUILDING ┃ Compiling source files...`,
    `┬ BUILDING ┃ Resolving dependencies...`,
    `├ BUILDING ┃ Resolving dependencies...`,
    `┼ BUILDING ┃ Bundling modules...`,
    `│ BUILDING ┃ Bundling modules...`,
    `┤ BUILDING ┃ Optimizing chunks...`,
    `┴ BUILDING ┃ Optimizing chunks...`
  ]
};

const STATUS_TESTING = {
  frames: [
    `▹ TESTING ┃ Running unit tests...`,
    `▸ TESTING ┃ Running unit tests...`,
    `▹ TESTING ┃ Running integration tests...`,
    `▸ TESTING ┃ Running integration tests...`,
    `▹ TESTING ┃ Running e2e tests...`,
    `▸ TESTING ┃ Running e2e tests...`,
    `▹ TESTING ┃ Calculating coverage...`,
    `▸ TESTING ┃ Calculating coverage...`
  ]
};

const STATUS_CONNECTING = {
  frames: [
    `⣾ CONNECTING ┃ Establishing TCP handshake...`,
    `⣽ CONNECTING ┃ Establishing TCP handshake...`,
    `⣻ CONNECTING ┃ Negotiating TLS...`,
    `⢿ CONNECTING ┃ Negotiating TLS...`,
    `⡿ CONNECTING ┃ Authenticating session...`,
    `⣟ CONNECTING ┃ Authenticating session...`,
    `⣯ CONNECTING ┃ Synchronizing state...`,
    `⣷ CONNECTING ┃ Synchronizing state...`
  ]
};

const STATUS_DEPLOYING = {
  frames: [
    `▁ DEPLOYING ┃ Pushing container to registry...`,
    `▂ DEPLOYING ┃ Pushing container to registry...`,
    `▃ DEPLOYING ┃ Updating kubernetes manifests...`,
    `▄ DEPLOYING ┃ Updating kubernetes manifests...`,
    `▅ DEPLOYING ┃ Rolling out new version...`,
    `▆ DEPLOYING ┃ Rolling out new version...`,
    `▇ DEPLOYING ┃ Validating health checks...`,
    `█ DEPLOYING ┃ Validating health checks...`
  ]
};

const STATUS_ANALYZING = {
  frames: [
    `◰ ANALYZING ┃ Parsing abstract syntax tree...`,
    `◳ ANALYZING ┃ Parsing abstract syntax tree...`,
    `◲ ANALYZING ┃ Computing complexity metrics...`,
    `◱ ANALYZING ┃ Computing complexity metrics...`,
    `◰ ANALYZING ┃ Detecting code smells...`,
    `◳ ANALYZING ┃ Detecting code smells...`,
    `◲ ANALYZING ┃ Generating insights report...`,
    `◱ ANALYZING ┃ Generating insights report...`
  ]
};

const STATUS_VERIFYING = {
  frames: [
    `● VERIFYING ┃ Checking digital signatures...`,
    `◐ VERIFYING ┃ Checking digital signatures...`,
    `◓ VERIFYING ┃ Validating checksums...`,
    `◑ VERIFYING ┃ Validating checksums...`,
    `◒ VERIFYING ┃ Confirming integrity...`,
    `● VERIFYING ┃ Confirming integrity...`
  ]
};

const STATUS_FINALIZING = {
  frames: [
    `∙ FINALIZING ┃ Cleaning up temporary files...`,
    `⋅ FINALIZING ┃ Cleaning up temporary files...`,
    `∘ FINALIZING ┃ Updating metadata...`,
    `∙ FINALIZING ┃ Updating metadata...`,
    `⋅ FINALIZING ┃ Writing logs...`,
    `∘ FINALIZING ┃ Writing logs...`
  ]
};

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 4. SEPARATORS & DIVIDERS (8+ Styles Variés)                                │
// └─────────────────────────────────────────────────────────────────────────────┘

const SEPARATORS = {
  // Simple styles
  SIMPLE: `────────────────────────────────────────────────────────────────────────────────`,
  DOUBLE: `════════════════════════════════════════════════════════════════════════════════`,
  DOTTED: `┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄`,
  DASHED: `╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌`,
  
  // Accented styles
  SECTION: `├────────────────────────────────────────────────────────────────────────────────┤`,
  HEAVY: `┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫`,
  
  // Thematic styles
  SECURITY: `╟─[SECURITY BOUNDARY]────────────────────────────────────────────────────────────╢`,
  BUILD: `╟─[BUILD STAGE]──────────────────────────────────────────────────────────────────╢`,
  NETWORK: `╟─[NETWORK ZONE]─────────────────────────────────────────────────────────────────╢`,
  
  // Compact variations
  COMPACT_SIMPLE: `──────────────────────────`,
  COMPACT_DOUBLE: `══════════════════════════`,
  COMPACT_SECTION: `├────────────────────────┤`
};

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 5. INFORMATION DISPLAYS (Métriques & Stats)                                │
// └─────────────────────────────────────────────────────────────────────────────┘

const INFO_DISPLAYS = {
  VERSION: `┃ v2.4.1 ┃ Node 20.11.0 ┃ Uptime: 7d 14h 23m ┃`,
  
  COMPLIANCE: `┃ ISO 27001:2022 ✓ ┃ GDPR Compliant ✓ ┃ SOC 2 Type II ✓ ┃`,
  
  NETWORK: `┃ Latency: 12ms ┃ Bandwidth: 847 Mbps ┃ Packet Loss: 0.00% ┃`,
  
  TEST_COVERAGE: `┃ Lines: 94.2% ┃ Branches: 91.7% ┃ Functions: 96.4% ┃ Statements: 93.8% ┃`,
  
  BUILD_TIME: `┃ Compile: 4.2s ┃ Bundle: 1.8s ┃ Optimize: 2.1s ┃ Total: 8.1s ┃`,
  
  SECURITY_STATUS: `┃ Auth: ACTIVE ┃ TLS: 1.3 ┃ Encryption: AES-256-GCM ┃ Last Scan: 2h ago ┃`,
  
  SYSTEM_METRICS: `┃ CPU: 34% ┃ MEM: 2.4/16 GB ┃ DISK: 187/512 GB ┃ TEMP: 48°C ┃`,
  
  GIT_STATUS: `┃ Branch: main ┃ Commits: +3 -0 ┃ Modified: 7 files ┃ Staged: 3 files ┃`
};

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 6. BREADCRUMBS & NAVIGATION (Path & Context)                               │
// └─────────────────────────────────────────────────────────────────────────────┘

const BREADCRUMBS = {
  PATH_SIMPLE: `~/projects/b-dev/src/app/core/services/auth.service.ts`,
  
  PATH_STYLED: `◆ home ▸ projects ▸ b-dev ▸ src ▸ app ▸ core ▸ services ▸ auth.service.ts`,
  
  GIT_CONTEXT: `[main ↑3] ~/projects/b-dev/src [+7 ~2]`,
  
  PROJECT_CONTEXT: `B.DEV ▸ Workspace: Production ▸ Environment: Staging ▸ Region: EU-West-1`,
  
  MULTI_LEVEL: `Root ▸ Infrastructure ▸ Kubernetes ▸ Deployments ▸ api-gateway-prod-v2.4.1`,
  
  HIERARCHY: `┌─ Organization: B.DEV
             ├─ Team: Core Engineering
             └─ Project: Authentication Service`
};

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 7. PROMPT STYLES (5+ Variations de Shell)                                  │
// └─────────────────────────────────────────────────────────────────────────────┘

const PROMPTS = {
  // Simple development prompt
  DEV_SIMPLE: `┃ dev ▸`,
  
  // Git-aware prompt
  GIT_AWARE: `┃ ~/b-dev [main] ▸`,
  
  // Secured/Authenticated prompt
  SECURED: `┃ [AUTH ✓] b.dev@production ▸`,
  
  // Root/Admin prompt
  ROOT: `╠═ root@system-core ▸`,
  
  // Production environment prompt
  PRODUCTION: `┃ [PROD] b-dev-api-v2 ▸`,
  
  // Multi-context prompt
  CONTEXT: `┃ [SECURE] b.dev │ main ↑3 │ EU-West │ 14:23:41 ▸`,
  
  // Minimal prompt
  MINIMAL: `▸`
};

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 8. COMPLETE STATUS MESSAGES (Success/Error/Warning/Info)                   │
// └─────────────────────────────────────────────────────────────────────────────┘

const MESSAGES = {
  SUCCESS: `┃ ✓ SUCCESS ┃ Build completed in 8.1s. Output: dist/bundle.min.js (847 KB)`,
  
  ERROR: `┃ ✗ ERROR ┃ Type 'string' is not assignable to type 'number' at auth.service.ts:42`,
  
  WARNING: `┃ ⚠ WARNING ┃ Deprecated API usage detected in 3 files. Update required before v3.0`,
  
  INFO: `┃ ℹ INFO ┃ New version available: v2.5.0. Run 'npm update' to upgrade.`,
  
  CRITICAL: `╔═══════════════════════════════════════════════════════════════════════════════╗
           ║ ⚠ CRITICAL SECURITY ALERT                                                     ║
           ║ CVE-2024-12345 detected in dependency 'example-lib@1.2.3'                    ║
           ║ Severity: HIGH | CVSS Score: 8.7/10                                          ║
           ║ Action Required: Update to version 1.2.4+ immediately                        ║
           ╚═══════════════════════════════════════════════════════════════════════════════╝`
};

// ┌─────────────────────────────────────────────────────────────────────────────┐
// │ 9. EXAMPLE USAGE SCENARIOS (Full CLI Sequences)                            │
// └─────────────────────────────────────────────────────────────────────────────┘

// Example 1: Full Build Sequence
const EXAMPLE_BUILD_SEQUENCE = `
╔═══════════════════════════════════════╗  ██████╗ ██╗   ██╗██╗██╗     ██████╗  ╔═══════════════════════════════════════╗
════════════════════════════════════════════════════════════════════════════════

┃ v2.4.1 ┃ Node 20.11.0 ┃ Uptime: 7d 14h 23m ┃
┃ [AUTH ✓] b.dev@production ▸ npm run build

⠋ LOADING ┃ Fetching resources...
[▰▰▰▰▰▰▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱▱] 24% ┃ Type checking definitions...
◐ PROCESSING ┃ Executing pipeline stage 2/4...
[▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰] 100% ┃ Build complete. Ready.

┃ Compile: 4.2s ┃ Bundle: 1.8s ┃ Optimize: 2.1s ┃ Total: 8.1s ┃
┃ ✓ SUCCESS ┃ Build completed in 8.1s. Output: dist/bundle.min.js (847 KB)

════════════════════════════════════════════════════════════════════════════════
`;

// Example 2: Security Audit Sequence
const EXAMPLE_SECURITY_SEQUENCE = `
╔═[SECURITY]═══════════════════════════╗  ███████╗███████╗ ██████╗██╗   ██╗██████╗ ██╗████████╗██╗   ██╗  ╔═══════╗
════════════════════════════════════════════════════════════════════════════════

┃ ISO 27001:2022 ✓ ┃ GDPR Compliant ✓ ┃ SOC 2 Type II ✓ ┃
┃ [SECURE] b.dev │ main ↑3 │ EU-West │ 14:23:41 ▸ npm run audit

▱▱▱▱▱▱▱▱▱▱ SCANNING ┃ Analyzing codebase structure...
[◓◓◓◓◓◓◓◓◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐◐] 29% ┃ Auditing authentication flows...
● VERIFYING ┃ Checking digital signatures...
[◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓◓] 100% ┃ Security audit complete. 0 critical issues.

┃ Auth: ACTIVE ┃ TLS: 1.3 ┃ Encryption: AES-256-GCM ┃ Last Scan: just now ┃
┃ ✓ SUCCESS ┃ All security checks passed. System compliant.

════════════════════════════════════════════════════════════════════════════════
`;

// ═══════════════════════════════════════════════════════════════════════════════
// EXPORT ALL COMPONENTS
// ═══════════════════════════════════════════════════════════════════════════════

export {
  HEADERS,
  PROGRESS_BUILD,
  PROGRESS_TEST,
  PROGRESS_NETWORK,
  PROGRESS_SECURITY,
  PROGRESS_OPTIMIZE,
  PROGRESS_DEPLOY,
  STATUS_LOADING,
  STATUS_PROCESSING,
  STATUS_SCANNING,
  STATUS_BUILDING,
  STATUS_TESTING,
  STATUS_CONNECTING,
  STATUS_DEPLOYING,
  STATUS_ANALYZING,
  STATUS_VERIFYING,
  STATUS_FINALIZING,
  SEPARATORS,
  INFO_DISPLAYS,
  BREADCRUMBS,
  PROMPTS,
  MESSAGES,
  EXAMPLE_BUILD_SEQUENCE,
  EXAMPLE_SECURITY_SEQUENCE
};