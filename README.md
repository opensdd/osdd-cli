# OpenSDD CLI (`osdd`)

The OpenSDD CLI lets you discover and run OpenSDD recipes from your terminal using your AI IDE of choice.

Think of it as Terraform for AI flows.

---

## Installation

### Option 1: Homebrew (macOS & Linux)

```bash
brew install opensdd/tap/osdd

# verify
osdd version
```

### Option 2: Download a Release Binary

1. Visit the [releases page](https://github.com/opensdd/osdd-cli/releases).
2. Download the artifact for your platform:
   - `osdd-macos-x64` (macOS)
   - `osdd-linux-x64` (Linux)
   - `osdd-windows-x64.exe` (Windows)
3. Place the binary on your `PATH` and make it executable (macOS/Linux):

```bash
chmod +x osdd
sudo mv osdd /usr/local/bin/
osdd version
```

Windows users can run `osdd.exe` directly or move it to a directory listed in `%PATH%`.

### Option 3: Build from Source (advanced)

Requires Go 1.25.1 or newer.

```bash
git clone https://github.com/opensdd/osdd-cli.git
cd osdd-cli
make build        # or: make build-dev for a dev build
./osdd version
```

---
## Example

Let's start with an example — the [documentation website](https://opensdd.ai/osdd-cli/) was prepared using an OpenSDD recipe!

We have created an `astro_site` recipe, which is instructed to generate a documentation website based on
[Astro Starlight](https://astro.build/themes/details/starlight/) and published to
[GitHub Pages](https://docs.github.com/en/pages). All you need to do is install the OSDD CLI

```bash
brew install opensdd/tap/osdd
```

and run the following command:

```bash
osdd recipe execute astro_site -i codex
```

![osdd astro docs](/resources/osdd_astro.png)

The CLI asks for user input, specifically in this case for the name of the web site, repo where the website should
be generated, any context information (other websites, local files, other git repos, etc) and for any other instructions
on how to generate the website. Then it launches the provided IDE (in this case - `codex`) and gets to work!

**How exactly does this work?**

In a nutshell, this command:
1. Downloaded this recipe - [astro_site](https://github.com/opensdd/recipes/tree/main/global/astro_site).
2. Asked for user input, which is also [pre-configured in the recipe](https://github.com/opensdd/recipes/blob/012cb69b08d261d1ae65661b0855b2c8fbb5075a/global/astro_site/recipe.yaml#L14-L27),
3. Created a new workspace directory in `~/osdd/workspace/astro_site` for [recipe's instruction](https://github.com/opensdd/recipes/blob/012cb69b08d261d1ae65661b0855b2c8fbb5075a/global/astro_site/recipe.yaml#L6C11-L6C36).
4. Downloaded all the context files and commands into the workspace.
5. Launched `codex` with the instruction to run `/astro_run` command, which maps to [this prompt](https://github.com/opensdd/recipes/blob/main/global/astro_site/resources/run.md).

After a while, the `codex` was done and the website was generated and after a couple tweaks, it was published and here
it is!
---

## Quick Start

1. **Pick a recipe**
   Browse [`opensdd/recipes`](https://github.com/opensdd/recipes/tree/main/global) to find an automation. Recipes are identified by the name of the folder in this repository, for example:
   `docs_update`

2. **Run the recipe**
   ```bash
   osdd recipe execute astro_site --ide claude
   ```
   - `--ide` names the IDE integration (`claude`, `codex`).
   - The CLI fetches the recipe, and prompts for any declared user inputs.

3. **Follow the prompts**
   Provide requested information (multi-line text, options, and so on). When you finish, the CLI materializes files, configures workspaces, and executes the recipe steps.

4. **Follow the execution in the IDE**
   The CLI will start the IDE requested and will optionally prompt it to start the work (depends on the recipe). You may need to pay attention to what the IDE is doing, since it may request permissions or confirmations.

---

## Command Reference

```bash
osdd recipe execute <ID>                # run a recipe by ID
osdd recipe execute <ID> --ide <name>   # required IDE identifier
```

---

## Running a Recipe

```bash
osdd recipe execute docs_update --ide claude
```

Flags to know:

- `--ide` selects which IDE integration to launch (Codex, Claude, etc.).
- `--recipe-file` lets you point to a local manifest when authoring new automations.

## Using Your Own Recipes

The CLI also supports running recipes from an arbitrary publicly available repository. All you need to do is to
create a folder `opensdd_recipes` in the repository root and add a folder with then name of the recipe in it and a
`recipe.yaml` as a recipe declaration.

For example, a recipe from
```
https://github.com/<OWNER>/<REPO_NAME>/blob/main/opensdd_recipes/<RECIPE_NAME>/recipe.yaml
```
can be executed using
```bash
osdd recipe execute <OWNER>/<REPO_NAME>/<RECIPE_NAME> --ide claude
```

## Why It Matters

- **Helper role** – The CLI is a facilitator, not a gatekeeper. It ensures every automation starts with the right context and guardrails.
- **Repeatability** – Recipes encode best practices, like mandatory planning or review loops, so teams don’t reinvent the process every time.
- **Observability** – Logs and generated artifacts live alongside your specs, enriching the OpenSDD knowledge base.
- **Shareable** – You can now just share id of a recipe you created and let others try it out. No need to copy/paste
  prompts, download files manually, etc.
- **Interoperability** – Recipes work the same way (well, almost) across different IDEs/CLIs. E.g. the example above
  can be executed with either Claude or Codex, despite Codex formally not having support for slash-commands
  or even project-level prompts. OSDD takes care of instructing each coding agent the right way.

## Support & Related Projects

- Open issues or feature requests in [opensdd/osdd-cli](https://github.com/opensdd/osdd-cli/issues).
- Explore the broader OpenSDD ecosystem:
  - [opensdd/osdd-api](https://github.com/opensdd/osdd-api) – Protobuf definitions & clients referenced by the CLI.
  - [opensdd/osdd-core](https://github.com/opensdd/osdd-core) – Core runtime used to materialize recipes.
  - [opensdd/recipes](https://github.com/opensdd/recipes) – Official recipe catalog (see its README for authoring guidance).

Licensed under the terms of the [LICENSE](LICENSE) file.
