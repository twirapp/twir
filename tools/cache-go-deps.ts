import { spawnSync } from 'node:child_process';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

const regexp = /use\s*\((.*?)\)/gs;
const goWork = await readFile(resolve(process.cwd(), 'go.work'), 'utf-8');

const packagesMatch = regexp.exec(goWork)![1];
const packages = packagesMatch
  .split('\n')
  .map(p => p.replaceAll('\n', '').replaceAll('\t', ''))
  .filter(p => Boolean(p));

process.stdout.write('Caching golang deps');

for (const packagePath of packages) {
  process.stdout.write('.');
  spawnSync('go mod download', {
    cwd: resolve(process.cwd(), packagePath),
  });
}