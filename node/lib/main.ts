import * as fs from 'fs';
import { promisify } from 'util';
import * as path from 'path';
import * as del from 'del';

const readFileAsync = promisify(fs.readFile);
const writeFileAsync = promisify(fs.writeFile);
const existsAsync = promisify(fs.exists);
const mkdirAsync = promisify(fs.mkdir);

const ROOT_DIR = 'huge_dir';

(async () => {
  const args = process.argv.slice(2);

  // iteration argument
  if (args.length < 1) {
    console.log('Missing iteration argument');
    process.exit(1);
  }

  const iteration = parseInt(args[0], 10);
  if (iteration <= 0) {
    console.log(`Invalid iteration value: ${iteration}`);
    process.exit(1);
  }

  // parseJSON argument
  let parseJSON = false;
  if (args.length >= 2 && args[1] === '--parse-json') {
    parseJSON = true;
  }

  if (await existsAsync(ROOT_DIR)) {
    await del([ROOT_DIR]);
  }

  await mkdirAsync(ROOT_DIR);
  const paths: string[] = new Array(iteration);
  const content = await readFileAsync('../common/bench_data.json');
  for (let i = 0; i < iteration; i++) {
    paths[i] = path.join(ROOT_DIR, i.toString() + '.json');
  }

  console.time('write');
  console.log(`Creating ${iteration} files...`);
  const writeJobs = paths.map((p) => writeFileAsync(p, content));
  await Promise.all(writeJobs);
  console.timeEnd('write');

  if (parseJSON) {
    console.log(`Reading and parsing ${iteration} files...`);
  } else {
    console.log(`Reading ${iteration} files...`);
  }
  console.time('read');
  const readJobs = paths.map(async (p) => {
    const bytes = await readFileAsync(p);
    if (parseJSON) {
      return JSON.parse(bytes.toString());
    } else {
      return bytes;
    }
  });
  await Promise.all(readJobs);
  console.timeEnd('read');

  console.log('Completed');
})();
