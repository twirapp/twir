import { execSync } from 'node:child_process';
import { platform } from 'node:os';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));

const workDir = resolve(__dirname);
const cwd = resolve(workDir, '..', '..');

const pathName = platform() === 'win32' ? 'Path' : 'PATH';

const binPath = resolve(cwd, '.bin');
const execOptions = {
	env: {
		...process.env,
		[pathName]: `${process.env[pathName]}:${binPath}`,
	},
};

execSync(`buf ${process.argv.slice(2).join(' ')}`, execOptions);
