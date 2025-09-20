module.exports = {
  root: true,
  extends: [
    'next/core-web-vitals',
    'plugin:import/recommended',
    'plugin:import/warnings',
    'plugin:tailwindcss/recommended',
    'prettier',
  ],
  plugins: ['@typescript-eslint', 'prettier', 'boundaries', 'tailwindcss'],
  rules: {
    'import/order': [
      'off',
      {
        alphabetize: {
          order: 'asc',
        },
      },
    ],
    'prettier/prettier': 'warn',
    'import/no-unresolved': 'off',
    'import/no-duplicates': 'off',
    'import/no-named-as-default': 'off',
    'import/no-named-as-default-member': 'off',
    '@next/next/no-img-element': 'off',
    // FSD boundaries rules
    'boundaries/element-types': [
      'error',
      {
        default: 'disallow',
        rules: [
          {
            from: 'shared',
            allow: ['shared'],
          },
          {
            from: 'entities',
            allow: ['shared'],
          },
          {
            from: 'features',
            allow: ['shared', 'entities'],
          },
          {
            from: 'widgets',
            allow: ['shared', 'entities', 'features'],
          },
          {
            from: 'page-components',
            allow: ['shared', 'entities', 'features', 'widgets'],
          },
          {
            from: 'app',
            allow: [
              'shared',
              'entities',
              'features',
              'widgets',
              'page-components',
              'app',
            ],
          },
        ],
      },
    ],
  },
  settings: {
    'import/parsers': {
      '@typescript-eslint/parser': ['.ts', '.tsx'],
    },
    'import/resolver': {
      typescript: {
        alwaysTryTypes: true,
      },
    },
    'boundaries/elements': [
      {
        type: 'shared',
        pattern: 'src/shared/**/*',
      },
      {
        type: 'entities',
        pattern: 'src/entities/**/*',
      },
      {
        type: 'features',
        pattern: 'src/features/**/*',
      },
      {
        type: 'widgets',
        pattern: 'src/widgets/**/*',
      },
      {
        type: 'page-components',
        pattern: 'src/page-components/**/*',
      },
      {
        type: 'app',
        pattern: 'src/app/**/*',
      },
    ],
  },
};
