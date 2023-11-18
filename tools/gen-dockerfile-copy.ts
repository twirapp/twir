import fs from 'fs';
import path from 'path';

function findPackageJsonFiles(startPath: string) {
	const result: string[] = [];
	if (!fs.existsSync(startPath)) {
		console.log('Directory path does not exist: ', startPath);
		return result;
	}

	const files = fs.readdirSync(startPath);

	for (const file of files) {
		const fullPath = path.join(startPath, file);
		if (!fs.statSync(fullPath).isDirectory()) {
			continue;
		}

		const hasPackage = fs.existsSync(path.join(fullPath, 'package.json'));

		if (hasPackage) {
			result.push(fullPath);
		}
	}

	return result;
}

function findGoModFiles(startPath: string) {
	const result: string[] = [];
	if (!fs.existsSync(startPath)) {
		console.log('Directory path does not exist: ', startPath);
		return result;
	}

	const files = fs.readdirSync(startPath);

	for (const file of files) {
		const fullPath = path.join(startPath, file);
		if (!fs.statSync(fullPath).isDirectory()) {
			continue;
		}

		const hasModOrSum = fs.existsSync(path.join(fullPath, 'go.mod')) || fs.existsSync(path.join(fullPath, 'go.sum'));
		if (hasModOrSum) {
			result.push(fullPath);
		}
	}

	return result;
}


const nodePaths = [
	...findPackageJsonFiles('./libs'),
	...findPackageJsonFiles('./apps'),
	...findPackageJsonFiles('./frontend'),
];

const goPaths = [
	...findGoModFiles('./libs'),
	...findGoModFiles('./apps'),
];

const generatedStrings = [
	...nodePaths.map((p) => `COPY ${p}/package.json ${p}/package.json`),
	...goPaths.map((p) => `COPY ${p}/go.mod ${p}/go.mod`),
].join('\n');

const originalFile = fs.readFileSync('./base.Dockerfile', 'utf-8');

function replaceBetweenComments(input: string, newContent: string) {
	// Define start and end comments
	const startComment = '# START COPYGEN';
	const endComment = '# END COPYGEN';

	// Define regular expression
	const regex = new RegExp(`${startComment}[\\s\\S]*${endComment}`);

	// Replace content between comments
	const replacedContent = input.replace(regex, `${startComment}\n${newContent}\n${endComment}`);

	return replacedContent;
}

const newFileContent = replaceBetweenComments(originalFile, generatedStrings);

fs.writeFileSync('./base.Dockerfile', newFileContent);
