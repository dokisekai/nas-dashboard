export const meta = {
  name: 'project-lifecycle',
  description: 'Complete project lifecycle workflow from requirements to deployment',
  phases: [
    { title: 'Analyze', detail: 'Requirements analysis and technical planning' },
    { title: 'Frontend', detail: 'Frontend implementation and testing' },
    { title: 'Backend', detail: 'Backend implementation and testing' },
    { title: 'Review', detail: 'Code review and quality assurance' },
    { title: 'Execute', detail: 'Execution and deployment' }
  ]
}

// Global workflow state
const workflowState = {
  requirements: '',
  techSpec: null,
  frontendCode: '',
  backendCode: '',
  testResults: null,
  reviewFindings: null,
  deploymentPlan: ''
}

// Phase 1: Requirements Analysis
phase('Analyze')
log('Starting requirements analysis and technical planning')

const requirementsSpec = await agent(
  `Analyze the current project state and provide comprehensive requirements analysis.
   Current project context: NAS dashboard application at /home/hserver/nas-dashboard

   Tasks:
   1. Analyze existing codebase structure
   2. Identify current features and gaps
   3. Generate detailed requirements document
   4. Create technical specification
   5. Identify potential risks and mitigation strategies

   Output format: Provide structured requirements with priorities (P0/P1/P2)`,
  { label: 'Requirements Analysis', phase: 'Analyze', schema: {
    type: 'object',
    properties: {
      projectOverview: { type: 'string', description: 'Project overview and goals' },
      existingFeatures: { type: 'array', items: { type: 'string' }, description: 'Currently implemented features' },
      requiredFeatures: { type: 'array', items: { type: 'object', properties: {
        name: { type: 'string' },
        priority: { type: 'string', enum: ['P0', 'P1', 'P2'] },
        description: { type: 'string' }
      }}},
      technicalGaps: { type: 'array', items: { type: 'string' }, description: 'Technical gaps identified' },
      risks: { type: 'array', items: { type: 'object', properties: {
        risk: { type: 'string' },
        impact: { type: 'string', enum: ['High', 'Medium', 'Low'] },
        mitigation: { type: 'string' }
      }}}
    }
  }}
)

workflowState.requirements = JSON.stringify(requirementsSpec, null, 2)
log('Requirements analysis completed')

// Generate technical specification
const techSpec = await agent(
  `Based on the requirements analysis, create a detailed technical specification.

   Requirements Analysis:
   ${workflowState.requirements}

   Tasks:
   1. Design system architecture
   2. Define API specifications
   3. Database schema design
   4. Component structure for frontend
   5. Integration patterns
   6. Security considerations

   Output format: Complete technical specification document`,
  { label: 'Technical Specification', phase: 'Analyze', schema: {
    type: 'object',
    properties: {
      architecture: { type: 'string', description: 'System architecture overview' },
      apiSpec: { type: 'array', items: { type: 'object', properties: {
        endpoint: { type: 'string' },
        method: { type: 'string' },
        description: { type: 'string' },
        auth: { type: 'boolean' }
      }}},
      databaseSchema: { type: 'string', description: 'Database schema design' },
      frontendComponents: { type: 'array', items: { type: 'string' }, description: 'Required frontend components' },
      securityConsiderations: { type: 'array', items: { type: 'string' } }
    }
  }}
)

workflowState.techSpec = techSpec
log('Technical specification completed')

// Generate markdown documentation
await agent(
  `Create comprehensive markdown documentation based on the analysis.

   Requirements:
   ${workflowState.requirements}

   Technical Spec:
   ${JSON.stringify(workflowState.techSpec, null, 2)}

   Generate:
   1. README.md with project overview
   2. REQUIREMENTS.md with detailed requirements
   3. ARCHITECTURE.md with system architecture
   4. API.md with API documentation
   5. DEVELOPMENT.md with development guidelines

   Write these files to /home/hserver/nas-dashboard/documentation/`,
  { label: 'Generate Documentation', phase: 'Analyze' }
)

// Phase 2: Frontend Implementation
phase('Frontend')
log('Starting frontend implementation and testing')

const frontendImplementation = await agent(
  `Implement frontend components based on technical specification.

   Technical Spec:
   ${JSON.stringify(workflowState.techSpec, null, 2)}

   Tasks:
   1. Analyze current frontend code at /home/hserver/nas-dashboard/frontend/src
   2. Identify missing components based on requirements
   3. Generate Vue.js components for required features
   4. Implement responsive layouts
   5. Add proper TypeScript types
   6. Include error handling and loading states
   7. Follow Vue 3 Composition API best practices

   Current issues to address:
   - Layout duplication problems
   - Missing monitoring data display
   - Debug mode implementation
   - Component structure optimization

   Output: Complete frontend code with proper structure`,
  { label: 'Frontend Implementation', phase: 'Frontend', schema: {
    type: 'object',
    properties: {
      components: { type: 'array', items: { type: 'object', properties: {
        name: { type: 'string' },
        path: { type: 'string' },
        code: { type: 'string' },
        description: { type: 'string' }
      }}},
      layoutFixes: { type: 'array', items: { type: 'string' }, description: 'Layout fixes applied' },
      newFeatures: { type: 'array', items: { type: 'string' }, description: 'New features implemented' },
      testingRecommendations: { type: 'array', items: { type: 'string' } }
    }
  }}
)

workflowState.frontendCode = JSON.stringify(frontendImplementation, null, 2)

// Frontend testing and validation
const frontendTests = await agent(
  `Create comprehensive frontend tests and validate implementation.

   Frontend Implementation:
   ${workflowState.frontendCode}

   Tasks:
   1. Generate unit tests for Vue components using Vitest
   2. Create integration tests for user flows
   3. Add E2E tests for critical paths
   4. Test responsive design
   5. Validate accessibility (ARIA)
   6. Performance testing
   7. Cross-browser compatibility checks

   Output: Test suite with coverage reports`,
  { label: 'Frontend Testing', phase: 'Frontend', schema: {
    type: 'object',
    properties: {
      unitTests: { type: 'array', items: { type: 'string' }, description: 'Unit test files created' },
      integrationTests: { type: 'array', items: { type: 'string' }, description: 'Integration test scenarios' },
      e2eTests: { type: 'array', items: { type: 'string' }, description: 'E2E test scenarios' },
      coverage: { type: 'number', description: 'Code coverage percentage' },
      accessibilityScore: { type: 'number', description: 'Accessibility score' },
      performanceMetrics: { type: 'object', description: 'Performance metrics' }
    }
  }}
)

log('Frontend implementation and testing completed')

// Phase 3: Backend Implementation
phase('Backend')
log('Starting backend implementation and testing')

const backendImplementation = await agent(
  `Implement backend services based on technical specification.

   Technical Spec:
   ${JSON.stringify(workflowState.techSpec, null, 2)}

   Current backend at /home/hserver/nas-dashboard/backend

   Tasks:
   1. Analyze existing Go backend code
   2. Implement missing API endpoints
   3. Add proper error handling and middleware
   4. Implement authentication and authorization
   5. Add database integration
   6. Implement WebSocket for real-time updates
   7. Add input validation and sanitization
   8. Implement proper logging and monitoring

   Output: Complete backend implementation with Go best practices`,
  { label: 'Backend Implementation', phase: 'Backend', schema: {
    type: 'object',
    properties: {
      apiEndpoints: { type: 'array', items: { type: 'object', properties: {
        path: { type: 'string' },
        method: { type: 'string' },
        handler: { type: 'string' },
        auth: { type: 'boolean' }
      }}},
      middleware: { type: 'array', items: { type: 'string' }, description: 'Middleware implemented' },
      databaseIntegration: { type: 'string', description: 'Database integration details' },
      websockets: { type: 'array', items: { type: 'string' }, description: 'WebSocket endpoints' },
      securityFeatures: { type: 'array', items: { type: 'string' }, description: 'Security features implemented' }
    }
  }}
)

workflowState.backendCode = JSON.stringify(backendImplementation, null, 2)

// Backend testing and validation
const backendTests = await agent(
  `Create comprehensive backend tests and validate implementation.

   Backend Implementation:
   ${workflowState.backendCode}

   Tasks:
   1. Generate unit tests for Go handlers and services
   2. Create integration tests for API endpoints
   3. Add load and stress testing
   4. Security testing (SQL injection, XSS, etc.)
   5. Database testing
   6. WebSocket testing
   7. Performance benchmarking

   Output: Complete backend test suite with coverage`,
  { label: 'Backend Testing', phase: 'Backend', schema: {
    type: 'object',
    properties: {
      unitTests: { type: 'array', items: { type: 'string' }, description: 'Unit test files' },
      apiTests: { type: 'array', items: { type: 'string' }, description: 'API integration tests' },
      loadTests: { type: 'array', items: { type: 'string' }, description: 'Load test scenarios' },
      securityTests: { type: 'array', items: { type: 'string' }, description: 'Security test results' },
      coverage: { type: 'number', description: 'Code coverage percentage' },
      performanceBenchmarks: { type: 'object', description: 'Performance benchmarks' }
    }
  }}
)

log('Backend implementation and testing completed')

// Phase 4: Code Review and Quality Assurance
phase('Review')
log('Starting comprehensive code review and quality assurance')

// Parallel review processes
const [securityReview, performanceReview, codeQualityReview] = await parallel([
  // Security review
  () => agent(
    `Conduct comprehensive security review of the entire codebase.

     Frontend: ${workflowState.frontendCode}
     Backend: ${workflowState.backendCode}

     Tasks:
     1. Check for OWASP Top 10 vulnerabilities
     2. Review authentication and authorization
     3. Check for sensitive data exposure
     4. Validate input sanitization
     5. Review API security
     6. Check for dependency vulnerabilities
     7. Review CORS and security headers

     Output: Detailed security findings with severity levels`,
    { label: 'Security Review', phase: 'Review', schema: {
      type: 'object',
      properties: {
        criticalIssues: { type: 'array', items: { type: 'object', properties: {
          issue: { type: 'string' },
          severity: { type: 'string', enum: ['Critical', 'High', 'Medium', 'Low'] },
          description: { type: 'string' },
          remediation: { type: 'string' }
        }}},
        securityScore: { type: 'number', description: 'Overall security score (0-100)' },
        recommendations: { type: 'array', items: { type: 'string' } }
      }
    }}
  ),

  // Performance review
  () => agent(
    `Conduct comprehensive performance review.

     Frontend: ${workflowState.frontendCode}
     Backend: ${workflowState.backendCode}

     Tasks:
     1. Database query optimization
     2. API response time analysis
     3. Frontend bundle size optimization
     4. Caching strategy review
     5. Memory leak detection
     6. Network optimization
     7. Rendering performance

     Output: Performance findings with optimization recommendations`,
    { label: 'Performance Review', phase: 'Review', schema: {
      type: 'object',
      properties: {
        bottlenecks: { type: 'array', items: { type: 'object', properties: {
          area: { type: 'string' },
          issue: { type: 'string' },
          impact: { type: 'string' },
          solution: { type: 'string' }
        }}},
        optimizationSuggestions: { type: 'array', items: { type: 'string' } },
        performanceScore: { type: 'number', description: 'Performance score (0-100)' }
      }
    }}
  ),

  // Code quality review
  () => agent(
    `Conduct comprehensive code quality review.

     Frontend: ${workflowState.frontendCode}
     Backend: ${workflowState.backendCode}

     Tasks:
     1. Code style and consistency
     2. Design pattern adherence
     3. Error handling completeness
     4. Code duplication detection
     5. Naming conventions
     6. Documentation completeness
     7. Test coverage analysis

     Output: Code quality findings with improvement suggestions`,
    { label: 'Code Quality Review', phase: 'Review', schema: {
      type: 'object',
      properties: {
        codeSmells: { type: 'array', items: { type: 'object', properties: {
          type: { type: 'string' },
          location: { type: 'string' },
          description: { type: 'string' },
          suggestion: { type: 'string' }
        }}},
        duplications: { type: 'array', items: { type: 'string' }, description: 'Code duplications found' },
        bestPracticeViolations: { type: 'array', items: { type: 'string' }, description: 'Best practice violations' },
        qualityScore: { type: 'number', description: 'Code quality score (0-100)' },
        improvementSuggestions: { type: 'array', items: { type: 'string' } }
      }
    }}
  )
])

workflowState.reviewFindings = {
  security: securityReview,
  performance: performanceReview,
  codeQuality: codeQualityReview
}

log('Code review and quality assurance completed')

// Generate review report
await agent(
  `Generate comprehensive review report based on all findings.

   Security Review:
   ${JSON.stringify(securityReview, null, 2)}

   Performance Review:
   ${JSON.stringify(performanceReview, null, 2)}

   Code Quality Review:
   ${JSON.stringify(codeQualityReview, null, 2)}

   Generate:
   1. SECURITY_REVIEW.md with security findings
   2. PERFORMANCE_REVIEW.md with performance findings
   3. CODE_QUALITY_REVIEW.md with code quality findings
   4. IMPROVEMENT_PLAN.md with prioritized improvement plan

   Write to /home/hserver/nas-dashboard/reviews/`,
  { label: 'Generate Review Reports', phase: 'Review' }
)

// Phase 5: Execution and Deployment
phase('Execute')
log('Starting execution and deployment planning')

const deploymentPlan = await agent(
  `Create comprehensive deployment and execution plan.

   Review Findings:
   ${JSON.stringify(workflowState.reviewFindings, null, 2)}

   Tasks:
   1. Create deployment checklist
   2. Define rollback procedures
   3. Setup monitoring and alerting
   4. Create CI/CD pipeline configuration
   5. Generate deployment scripts
   6. Create runbook for operations
   7. Setup maintenance procedures

   Output: Complete deployment package`,
  { label: 'Deployment Planning', phase: 'Execute', schema: {
    type: 'object',
    properties: {
      deploymentChecklist: { type: 'array', items: { type: 'string' } },
      rollbackProcedures: { type: 'array', items: { type: 'string' } },
      monitoringSetup: { type: 'array', items: { type: 'string' } },
      cicdPipeline: { type: 'string', description: 'CI/CD pipeline configuration' },
      deploymentScripts: { type: 'array', items: { type: 'object', properties: {
        name: { type: 'string' },
        script: { type: 'string' },
        description: { type: 'string' }
      }}},
      operationsRunbook: { type: 'string', description: 'Operations runbook' }
    }
  }}
)

workflowState.deploymentPlan = JSON.stringify(deploymentPlan, null, 2)

// Generate deployment scripts and documentation
await agent(
  `Generate deployment scripts and documentation.

   Deployment Plan:
   ${workflowState.deploymentPlan}

   Generate:
   1. deploy.sh - Main deployment script
   2. rollback.sh - Rollback script
   3. setup-monitoring.sh - Monitoring setup script
   4. DEPLOYMENT.md - Deployment guide
   5. OPERATIONS.md - Operations manual

   Write to /home/hserver/nas-dashboard/deployment/`,
  { label: 'Generate Deployment Scripts', phase: 'Execute' }
)

log('Deployment package generated')

// Generate final project summary
const finalSummary = await agent(
  `Generate final project summary and recommendations.

   Complete Workflow State:
   ${JSON.stringify({
     requirements: workflowState.requirements,
     techSpec: workflowState.techSpec,
     frontendCode: 'Implementation completed',
     backendCode: 'Implementation completed',
     reviewFindings: workflowState.reviewFindings,
     deploymentPlan: 'Deployment package ready'
   }, null, 2)}

   Generate comprehensive project summary including:
   1. Project completion status
   2. Key achievements
   3. Remaining recommendations
   4. Next steps
   5. Success metrics

   Output as PROJECT_SUMMARY.md in /home/hserver/nas-dashboard/`,
  { label: 'Generate Final Summary', phase: 'Execute' }
)

// Return final workflow results
return {
  success: true,
  phasesCompleted: ['Analyze', 'Frontend', 'Backend', 'Review', 'Execute'],
  deliverables: {
    documentation: [
      'documentation/README.md',
      'documentation/REQUIREMENTS.md',
      'documentation/ARCHITECTURE.md',
      'documentation/API.md',
      'documentation/DEVELOPMENT.md'
    ],
    reviews: [
      'reviews/SECURITY_REVIEW.md',
      'reviews/PERFORMANCE_REVIEW.md',
      'reviews/CODE_QUALITY_REVIEW.md',
      'reviews/IMPROVEMENT_PLAN.md'
    ],
    deployment: [
      'deployment/deploy.sh',
      'deployment/rollback.sh',
      'deployment/setup-monitoring.sh',
      'deployment/DEPLOYMENT.md',
      'deployment/OPERATIONS.md'
    ],
    summary: 'PROJECT_SUMMARY.md'
  },
  metrics: {
    securityScore: securityReview.securityScore,
    performanceScore: performanceReview.performanceScore,
    codeQualityScore: codeQualityReview.qualityScore,
    overallScore: (securityReview.securityScore + performanceReview.performanceScore + codeQualityReview.qualityScore) / 3
  },
  recommendations: {
    immediate: securityReview.criticalIssues.map((issue) => issue.remediation),
    shortTerm: performanceReview.optimizationSuggestions.slice(0, 5),
    longTerm: codeQualityReview.improvementSuggestions.slice(0, 5)
  }
}