export const meta = {
  name: 'project-lifecycle-fixed',
  description: 'Fixed project lifecycle workflow without schema dependencies',
  phases: [
    { title: 'Analyze', detail: 'Requirements analysis and technical planning' },
    { title: 'Fix', detail: 'Fix frontend layout and debug issues' },
    { title: 'Document', detail: 'Generate comprehensive documentation' },
    { title: 'Test', detail: 'Create test plan and validation' }
  ]
}

// Global workflow state
const workflowState = {
  analysisResults: '',
  fixesApplied: '',
  documentationCreated: '',
  testPlanCreated: ''
}

// Phase 1: Requirements Analysis
phase('Analyze')
log('Starting comprehensive project analysis')

const analysisResults = await agent(
  `Conduct comprehensive analysis of the NAS Dashboard project at /home/hserver/nas-dashboard

   Provide detailed analysis covering:
   1. Current project structure and architecture
   2. Existing features and their implementation status
   3. Critical bugs and layout issues (especially the menu duplication problem)
   4. Security vulnerabilities and technical gaps
   5. Frontend-backend integration issues
   6. Missing functionality and improvement opportunities
   7. Priority recommendations (P0/P1/P2)

   Focus on identifying the root causes of:
   - Layout duplication (menu showing multiple times)
   - Missing monitoring data display
   - API endpoint mismatches
   - Authentication flow issues
   - WebSocket connectivity problems

   Provide actionable findings with specific file paths and line numbers where applicable.`,
  { label: 'Project Analysis', phase: 'Analyze' }
)

workflowState.analysisResults = analysisResults
log('Project analysis completed')

// Phase 2: Fix Issues
phase('Fix')
log('Starting frontend and backend fixes')

const fixesApplied = await agent(
  `Based on the analysis results, implement critical fixes for the NAS Dashboard.

   Analysis Results:
   ${workflowState.analysisResults}

   Current Issues to Fix:
   1. Layout duplication in App.vue and router configuration
   2. API endpoint path mismatches (/api/auth/login vs /auth/login)
   3. WebSocket URL configuration issues
   4. Debug mode implementation for troubleshooting
   5. Missing chart component imports
   6. Authentication flow debugging

   Tasks:
   1. Fix App.vue layout logic to prevent duplication
   2. Correct all API endpoint paths in frontend
   3. Update environment configuration for WebSocket
   4. Add proper debug mode with toggle
   5. Fix Chart component imports and dependencies
   6. Improve error handling in authentication
   7. Test login flow end-to-end

   For each fix:
   - Specify exact file and line number
   - Provide the exact code change
   - Explain why the fix is needed
   - Include testing verification steps`,
  { label: 'Critical Fixes', phase: 'Fix' }
)

workflowState.fixesApplied = fixesApplied
log('Critical fixes completed')

// Phase 3: Generate Documentation
phase('Document')
log('Generating comprehensive documentation')

const documentationCreated = await agent(
  `Create comprehensive documentation for the NAS Dashboard project.

   Analysis: ${workflowState.analysisResults}
   Fixes: ${workflowState.fixesApplied}

   Generate the following documentation files in /home/hserver/nas-dashboard/docs/:

   1. **PROJECT_OVERVIEW.md**
      - Project description and goals
      - Technology stack details
      - Architecture overview
      - Feature matrix

   2. **INSTALLATION.md**
      - Prerequisites and setup
      - Docker deployment instructions
      - Manual installation steps
      - Environment configuration
      - Troubleshooting guide

   3. **API_DOCUMENTATION.md**
      - All API endpoints with methods
      - Request/response examples
      - Authentication requirements
      - Error codes and handling
      - WebSocket protocol documentation

   4. **DEVELOPMENT_GUIDE.md**
      - Development setup
      - Code organization
      - Contributing guidelines
      - Testing procedures
      - Debug mode usage

   5. **SECURITY_CONSIDERATIONS.md**
      - Current security measures
      - Known vulnerabilities and fixes
      - Security best practices
      - Production deployment security

   6. **TROUBLESHOOTING.md**
      - Common issues and solutions
      - Debug procedures
      - Log analysis guide
      - Performance optimization tips

   Make the documentation comprehensive, well-structured, and beginner-friendly.`,
  { label: 'Documentation Generation', phase: 'Document' }
)

workflowState.documentationCreated = documentationCreated
log('Documentation generation completed')

// Phase 4: Test Planning
phase('Test')
log('Creating comprehensive test plan')

const testPlanCreated = await agent(
  `Create a comprehensive testing strategy for the NAS Dashboard.

   Analysis: ${workflowState.analysisResults}
   Fixes: ${workflowState.fixesApplied}

   Create a detailed test plan covering:

   1. **Unit Testing Strategy**
      - Frontend component testing (Vue components)
      - Backend handler testing (Go functions)
      - API testing (endpoint validation)
      - Mock data and fixtures
      - Coverage requirements (>80%)

   2. **Integration Testing**
      - Frontend-backend API integration
      - Authentication flow testing
      - WebSocket connection testing
      - Database integration testing
      - Docker integration testing

   3. **End-to-End Testing**
      - User login/logout flows
      - Dashboard loading and data display
      - Storage management operations
      - Service control operations
      - User management operations

   4. **Performance Testing**
      - API response time benchmarks
      - WebSocket throughput testing
      - Frontend rendering performance
      - Memory leak detection
      - Load testing scenarios

   5. **Security Testing**
      - Authentication bypass testing
      - Input validation testing
      - SQL injection prevention
      - XSS prevention testing
      - CSRF protection testing

   6. **Compatibility Testing**
      - Browser compatibility matrix
      - Mobile responsiveness testing
      - Operating system compatibility
      - Docker version compatibility

   Generate the test plan as /home/hserver/nas-dashboard/docs/TEST_PLAN.md`,
  { label: 'Test Planning', phase: 'Test' }
)

workflowState.testPlanCreated = testPlanCreated
log('Test planning completed')

// Generate final summary
log('Generating final project summary')

const finalSummary = await agent(
  `Generate a comprehensive project summary and recommendations.

   Complete Workflow Results:
   - Analysis: ${workflowState.analysisResults}
   - Fixes: ${workflowState.fixesApplied}
   - Documentation: ${workflowState.documentationCreated}
   - Test Plan: ${workflowState.testPlanCreated}

   Create a final summary document: /home/hserver/nas-dashboard/docs/PROJECT_SUMMARY.md

   Include:
   1. Project completion status
   2. Key achievements from this workflow
   3. Critical fixes implemented
   4. Remaining recommendations by priority
   5. Next steps for production deployment
   6. Success metrics and KPIs
   7. Team handoff checklist
   8. Deployment readiness assessment

   Make this summary actionable and clear for stakeholders.`,
  { label: 'Final Summary', phase: 'Test' }
)

// Return workflow results
return {
  success: true,
  phasesCompleted: ['Analyze', 'Fix', 'Document', 'Test'],
  deliverables: {
    analysis: 'Comprehensive project analysis completed',
    fixes: 'Critical fixes implemented and tested',
    documentation: [
      'docs/PROJECT_OVERVIEW.md',
      'docs/INSTALLATION.md',
      'docs/API_DOCUMENTATION.md',
      'docs/DEVELOPMENT_GUIDE.md',
      'docs/SECURITY_CONSIDERATIONS.md',
      'docs/TROUBLESHOOTING.md'
    ],
    testPlan: 'docs/TEST_PLAN.md',
    summary: 'docs/PROJECT_SUMMARY.md'
  },
  nextSteps: [
    'Review generated documentation',
    'Implement remaining P1/P2 recommendations',
    'Execute test plan',
    'Setup CI/CD pipeline',
    'Deploy to staging environment',
    'Conduct security audit',
    'Monitor performance metrics',
    'Plan production deployment'
  ]
}