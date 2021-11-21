module.exports = {
  extends: [
    'eslint:recommended',
    'plugin:react/recommended', // Uses the recommended rules from @eslint-plugin-react
    'plugin:prettier/recommended',
  ],
  parserOptions: {
    ecmaVersion: 6,
    browser: true,
    sourceType: 'module',
  },
  env: {
    node: true,
    es6: true,
  },
  rules: {
    'no-unused-vars': [
      'warn',
      { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
    ],
  },
  overrides: [
    {
      files: ['*.ts', '*.tsx'],
      parser: '@typescript-eslint/parser',
      plugins: ['@typescript-eslint'],
      parserOptions: {
        ecmaVersion: 2020, // Allows for the parsing of modern ECMAScript features
        sourceType: 'module', // Allows for the use of imports
        ecmaFeatures: {
          jsx: true, // Allows for the parsing of JSX
        },
      },
      settings: {
        react: {
          version: 'detect', // Tells eslint-plugin-react to automatically detect the version of React to use
        },
      },
      extends: [
        'plugin:@typescript-eslint/eslint-recommended',
        'plugin:@typescript-eslint/recommended',
      ],
      rules: {
        // Follow eslint-config-prettier by disabling no-extra-semi.
        // This conflicts with prettier and is used as a prefix for things like inline type casts.
        '@typescript-eslint/no-extra-semi': 'off',
        'react/jsx-uses-react': 'off',
        'react/react-in-jsx-scope': 'off',
        '@typescript-eslint/no-unused-vars': [
          'warn',
          { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
        ],
      },
    },
  ],
}
