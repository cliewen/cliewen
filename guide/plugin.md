# Install from Claude Code

If you already work inside Claude Code, you can install `clue` without leaving the session:

```
/plugin marketplace add cliewen/cliewen
/plugin install cliewen@cliewen
```

The plugin adds one skill, `/cliewen:setup`. It works out your platform, installs `clue` through the same checksum-verifying script the [Install](./install) page documents, confirms the binary reports a release version, and then **asks** before running `clue init`. Nothing is written into your repository until you say so.

That is the whole plugin. The rest of this page explains its boundaries before you look for capabilities it does not provide.

## What the plugin does not install

It does not ship `clue-analysis`, `clue-delta`, `clue-extract`, `clue-plan`, `clue-upgrade`, or `clue-verify` — the six managed skills that run the Cliewen loop. Those arrive when you run `clue init` in a repository, and only then.

This is deliberate. Most plugins do the opposite.

Cliewen's skills are committed repository files. `clue init` writes them into `.agents/skills/`, each stamped with the version of the binary that wrote it. `clue validate` fails if the binary and skills disagree. That keeps the judge, guidance, and corpus conventions in the same generation, and you can see the version in the diff.

A plugin works differently. Its components are copied into a per-user cache and enabled through settings, outside `.agents/skills/` where the drift check cannot reach them. One copy would span every repository you open, while a stamped skill must be pinned per repository. Bundling the six skills into the plugin would create a second, unversioned set of instructions beside the committed set. It could conflict with a repository pinned to another release, with nothing to report the conflict.

The plugin installs the binary and stops there. `clue init` is the only supported way to add Cliewen skills to a repository. It puts them in version control, where collaborators and reviewers can read them and `clue validate` can check them.

The same reasoning applies to copying a Cliewen skill into your personal `~/.claude/skills/` directory by hand. It will appear to work, and it will be invisible to the check that exists to catch it going stale.

## Which route to use

The plugin and the install script install the same binary from the same release, and neither is more official than the other.

| You are | Use |
|---|---|
| Working in Claude Code | The plugin — it saves you the context switch |
| At a terminal, or scripting a machine | [`install.sh` or `install.ps1`](./install#install-clue) |
| Setting up CI | The [pinned, checksum-verified binary](./ci-wall) the generated workflow installs — never a plugin |

Upgrading `clue` moves the binary only. For a repository already using Cliewen, preview and apply the coordinated migration with `clue migrate`. The drift report tells you whether the repository still needs work. [Operate safely](./operations) explains the recovery boundary.

## Next

[Watch the judge accept and reject a thread.](./getting-started)
