import { defineConfig } from 'orval';

export default defineConfig({
  api: {
    input: {
      target: `${process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'}/openapi.json`,
    },

    output: {
      target: './src/orval/soloist.ts',
      override: {
        mutator: {
          path: './src/orval/apiClient.ts',
        },
      },
    },
    hooks: {
      afterAllFilesWrite: 'prettier --write',
    },
  },
});