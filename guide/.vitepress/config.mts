import { defineConfig } from "vitepress";
import { MermaidMarkdown } from "vitepress-plugin-mermaid";

const mermaidConfigId = "virtual:mermaid-config";
const resolvedMermaidConfigId = `\0${mermaidConfigId}`;
// The package's withMermaid helper injects its renderer into every page.
// Supplying the renderer's virtual config directly lets the theme load it
// only when a page contains a Mermaid fence.
const mermaidConfigPlugin = {
  name: "cliewen-mermaid-config",
  resolveId(id: string) {
    return id === mermaidConfigId ? resolvedMermaidConfigId : undefined;
  },
  load(id: string) {
    if (id === resolvedMermaidConfigId) {
      return "export default { securityLevel: 'loose', startOnLoad: false };";
    }
  },
};

export default defineConfig({
  lang: "en-US",
  title: "Cliewen",
  description: "A verifiable thread from goal to test for agent-driven development.",
  base: "/",
  cleanUrls: true,
  ignoreDeadLinks: false,
  markdown: {
    config: (md) => MermaidMarkdown(md),
  },
  vite: {
    plugins: [mermaidConfigPlugin],
    optimizeDeps: {
      include: [
        "mermaid",
        "@braintree/sanitize-url",
        "dayjs",
        "cytoscape-cose-bilkent",
        "cytoscape",
      ],
    },
    resolve: {
      alias: {
        "dayjs/plugin/advancedFormat.js": "dayjs/esm/plugin/advancedFormat",
        "dayjs/plugin/customParseFormat.js": "dayjs/esm/plugin/customParseFormat",
        "dayjs/plugin/isoWeek.js": "dayjs/esm/plugin/isoWeek",
        "cytoscape/dist/cytoscape.umd.js": "cytoscape/dist/cytoscape.esm.js",
      },
    },
  },
  head: [
    ["meta", { name: "theme-color", content: "#3b5bdb" }],
  ],
  themeConfig: {
    nav: [
      { text: "Guide", link: "/what-is-cliewen" },
      { text: "Get started", link: "/getting-started" },
      { text: "GitHub", link: "https://github.com/cliewen/cliewen" },
    ],
    sidebar: [
      {
        text: "Start here",
        items: [
          { text: "What is Cliewen?", link: "/what-is-cliewen" },
          { text: "Get started", link: "/getting-started" },
          { text: "Greenfield and brownfield", link: "/adoption" },
        ],
      },
      {
        text: "How the method works",
        items: [
          { text: "The verifiable thread", link: "/methodology" },
          { text: "The corpus", link: "/corpus" },
          { text: "The change loop", link: "/change-loop" },
          { text: "The skills", link: "/skills" },
        ],
      },
    ],
    search: {
      provider: "local",
    },
    socialLinks: [
      { icon: "github", link: "https://github.com/cliewen/cliewen" },
    ],
    editLink: {
      pattern: "https://github.com/cliewen/cliewen/edit/main/guide/:path",
      text: "Edit this page on GitHub",
    },
    footer: {
      message: "Released under the Apache 2.0 License.",
      copyright: "Cliewen",
    },
  },
});
