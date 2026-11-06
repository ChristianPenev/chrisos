import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTs from "eslint-config-next/typescript";
import prettier from "eslint-config-prettier/flat";

const aliases = "@{components,shadcn,ui,domains,lib,hooks,generated}/**";

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  prettier,
  {
    rules: {
      "max-len": [
        "error",
        { code: 100, ignoreUrls: true, ignoreRegExpLiterals: true, ignoreTemplateLiterals: true },
      ],
      "import/order": [
        "error",
        {
          groups: [["builtin", "external"], "internal", ["parent", "sibling", "index"]],
          pathGroups: [{ pattern: aliases, group: "internal" }],
          pathGroupsExcludedImportTypes: ["builtin"],
          "newlines-between": "always",
          alphabetize: { order: "asc", caseInsensitive: true },
        },
      ],
      "no-restricted-imports": [
        "error",
        {
          patterns: [
            { group: ["@/*"], message: "Use a specific alias (@shadcn, @ui, @domains, @lib)." },
            {
              group: ["../*"],
              message: "Import across folders through an alias; ./ is for the same folder.",
            },
          ],
        },
      ],
    },
  },
  {
    files: ["src/components/shadcn/**"],
    rules: { "max-len": "off" },
  },
  globalIgnores([
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
    "src/generated/**",
    "**/*.generated.ts",
  ]),
]);

export default eslintConfig;
