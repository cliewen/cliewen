# Install from Claude Code

If you already work inside Claude Code, you can install `clue` without leaving the session:

```
/plugin marketplace add cliewen/cliewen
/plugin install cliewen@cliewen
```

The plugin adds one skill, `/cliewen:setup`. It works out your platform, installs `clue` through the same checksum-verifying script the [Install](./install) page documents, confirms the binary reports a release version, and then **asks** before running `clue init`. Nothing is written into your repository until you say so.

That is the entire plugin. The rest of this page is about what it does not do, which is the part worth understanding before you go looking for it.

## What the plugin does not install

It does not ship `clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, `clue-upgrade`, or `clue-verify` — the six managed skills that run the Cliewen loop. Those arrive when you run `clue init` in a repository, and only then.

This is deliberate, and it is worth a paragraph because the opposite is what most plugins do.

Cliewen's skills are committed repository files. `clue init` writes them into `.agents/skills/`, each stamped with the version of the binary that wrote it, and `clue validate` fails if the binary and those skills ever disagree. That check is what makes a green validate mean something: the judge, the guidance, and the corpus conventions in a repository are the same generation, and you can see the version in the diff.

A plugin does not work that way. Its components are copied into a per-user cache and enabled through a settings file, so they sit outside `.agents/skills/` where the drift check cannot reach them — one copy, spanning every repository you open, while the whole point of a stamped skill is that it is pinned per repository. Bundle the six skills into the plugin and you get a second, unversioned set of Cliewen instructions beside the committed set, contradicting whichever repository happens to pin a different release, with nothing reporting it.

So the plugin installs the binary and gets out of the way. `clue init` is the only supported way to put Cliewen skills in a repository, and it is the right one: the skills land in version control, where your collaborators get them, your reviewers can read them, and `clue validate` can check them.

The same reasoning applies to copying a Cliewen skill into your personal `~/.claude/skills/` directory by hand. It will appear to work, and it will be invisible to the check that exists to catch it going stale.

## Which route to use

The plugin and the install script install the same binary from the same release, and neither is more official than the other.

| You are | Use |
|---|---|
| Working in Claude Code | The plugin — it saves you the context switch |
| At a terminal, or scripting a machine | [`install.sh` or `install.ps1`](./install#install-clue) |
| Setting up CI | The [pinned, checksum-verified binary](./ci-wall) the generated workflow installs — never a plugin |

Upgrading `clue` moves the binary only. For a repository already using Cliewen that is half an upgrade; preview and apply the coordinated repository migration with `clue migrate`, and use the drift report as the check that the pair still needs work. [Operate safely](./operations) explains the recovery boundary.

## Next

[Watch the judge accept and reject a thread.](./getting-started)
