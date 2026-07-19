import { withMermaid } from "vitepress-plugin-mermaid";

export default withMermaid({
  lang: "en-US",
  title: "Cliewen",
  description: "A verifiable thread from goal to test for agent-driven development.",
  base: "/cliewen/",
  cleanUrls: true,
  ignoreDeadLinks: false,
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
