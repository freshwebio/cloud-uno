const CracoAlias = require('craco-alias')

module.exports = {
  plugins: [
    {
      plugin: CracoAlias,
      options: {
        source: 'tsconfig',
        baseUrl: '.',
        tsConfigPath: './tsconfig.extend.json',
      },
    },
  ],
  jest: {
    configure: {
      testResultsProcessor: 'jest-sonar-reporter',
      coverageThreshold: {
        global: {
          branches: 80,
          functions: 80,
          lines: 80,
          statements: 80,
        },
      },
      collectCoverageFrom: [
        'src/**/*.{ts,tsx}',
        '!<rootDir>/node_modules/',
        '!src/test-utils/**/*.{ts,tsx}',
        // The following files come from CRA boilerplate
        // and shouldn't be changed, so not worth worrying about code
        // coverage or trying to write tests for them.
        '!src/reportWebVitals.ts',
        '!src/serviceWorkerRegistration.ts',
      ],
    },
  },
}
