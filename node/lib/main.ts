import * as fs from 'fs';
import { promisify } from 'util';
import * as path from 'path';
import * as del from 'del';

const readFileAsync = promisify(fs.readFile);
const writeFileAsync = promisify(fs.writeFile);
const existsAsync = promisify(fs.exists);
const mkdirAsync = promisify(fs.mkdir);

const ROOT_DIR = 'huge_dir';
const ITERATION = 10000;

(async () => {
  if (await existsAsync(ROOT_DIR)) {
    await del([ROOT_DIR]);
  }

  await mkdirAsync(ROOT_DIR);
  const paths: string[] = new Array(ITERATION);
  const content = await readFileAsync('../common/bench_data.json');
  for (let i = 0; i < ITERATION; i++) {
    paths[i] = path.join(ROOT_DIR, i.toString() + '.json');
  }

  console.time('write');
  console.log(`Creating ${ITERATION} files...`);
  const writeJobs = paths.map((p) => writeFileAsync(p, content));
  await Promise.all(writeJobs);
  console.timeEnd('write');

  console.log(`Reading ${ITERATION} files...`);
  console.time('read');
  const readJobs = paths.map((p) => readFileAsync(p));
  await Promise.all(readJobs);
  console.timeEnd('read');

  console.log('Completed');
})();
