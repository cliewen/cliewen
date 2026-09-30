// Vitest discovery example. Test callbacks are not run. Adapt to the pinned API.
import path from 'node:path';
import { pathToFileURL } from 'node:url';

const ac = /^[A-Z][A-Z0-9]*(?:-[A-Z][A-Z0-9]*)*-\d+[a-z]*$/;
const types = new Set(['Unit', 'Integration', 'E2E', 'Performance']);

export function collectModules(modules, root, files) {
  const references = [], diagnostics = [], seen = new Map();
  const allowed = new Set(files);
  const report = (source, subject, message) => diagnostics.push({ path: source, subject, message });
  function visit(collection, source, parents = []) {
    for (const task of collection) {
      const subject = [...parents, task.name].join(' > ');
      const tags = task.options?.tags ?? [];
      const ids = tags.filter(tag => ac.test(tag));
      const bracket = /^\[([A-Z][A-Z0-9]*(?:-[A-Z][A-Z0-9]*)*-\d+[a-z]*) (Unit|Integration|E2E|Performance) (Positive|Negative)\]/.exec(task.name);
      if (!bracket && /^\[[A-Z][A-Z0-9-]*-\d/.test(task.name)) {
        report(source, subject, 'malformed executable evidence title');
        continue;
      }
      if (task.type === 'suite') {
        if (ids.length || bracket) report(source, subject, 'container AC metadata cannot prove an executable');
        visit(task.children, source, [...parents, task.name]);
        continue;
      }
      const metadata = [];
      if (ids.length) {
        const classes = tags.filter(tag => types.has(tag));
        const directions = tags.filter(tag => tag === 'positive' || tag === 'negative');
        if (ids.length !== 1 || classes.length !== 1 || directions.length !== 1) {
          report(source, subject, 'ambiguous or incomplete executable tags');
          continue;
        }
        metadata.push({ id: ids[0], type: classes[0], direction: directions[0] });
      }
      if (bracket) metadata.push({ id: bracket[1], type: bracket[2], direction: bracket[3].toLowerCase() });
      if (!metadata.length) continue;
      if (!allowed.has(source)) { report(source, subject, 'discovered source is outside producer input scope'); continue; }
      if (metadata.some(m => JSON.stringify(m) !== JSON.stringify(metadata[0]))) {
        report(source, subject, 'native tags and fallback title disagree'); continue;
      }
      // One each declaration counts once; location groups its expanded instances.
      const identity = task.options?.each
        ? `${parents.join(' > ')} > parameterized ${task.location?.line}:${task.location?.column}`
        : subject;
      if (task.options?.each && !task.location) {
        report(source, subject, 'parameterized executable needs declaration location'); continue;
      }
      const key = source + '\0' + identity;
      const ref = { ...metadata[0], path: source, subject: identity };
      if (seen.has(key)) {
        if (!task.options?.each || JSON.stringify(seen.get(key)) !== JSON.stringify(ref)) report(source, subject, 'duplicate or conflicting executable');
      } else { seen.set(key, ref); references.push(ref); }
    }
  }
  for (const module of modules) {
    const source = path.relative(root, module.moduleId).split(path.sep).join('/');
    if (module.errors?.().length) report(source, 'module', 'framework collection failed');
    visit(module.children, source);
  }
  return { references, diagnostics };
}

if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
  let input = ''; for await (const chunk of process.stdin) input += chunk;
  const request = JSON.parse(input);
  const { createVitest } = await import('vitest/node');
  const vitest = await createVitest('test', { root: request.root, watch: false, reporters: [], includeTaskLocation: true });
  try {
    const result = await vitest.collect(undefined, { staticParse: false });
    if (result.unhandledErrors?.length) throw new Error('Vitest collection returned unhandled errors');
    process.stdout.write(JSON.stringify(collectModules(result.testModules, request.root, request.files)));
  } finally { await vitest.close(); }
}
