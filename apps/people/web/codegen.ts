import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  overwrite: true,
  schema: "../api/schema/*.graphqls",
  documents: ["src/graphql/**/*.graphql"],
  ignoreNoDocuments: true,
  config: {
    inputValueDeprecation: true,
    scalars: {
      Time: "string",
    },
  },
  generates: {
    "src/generated/types.ts": {
      plugins: ["typescript"],
    },
    "src/graphql/": {
      preset: "near-operation-file",
      presetConfig: {
        extension: ".generated.ts",
        baseTypesPath: "../generated/types",
      },
      plugins: ["typescript-operations", "typed-document-node"],
    },
  },
};

export default config;
