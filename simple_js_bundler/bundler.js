import fs from 'fs';
import babylon from 'babylon';

function createAsset(filename) {
  const src = fs.readFileSync(filename, 'utf-8');
  const parsedAst = babylon.parse(src, { sourceType: 'module' });
  console.log(parsedAst);
}

const rootFolder = 'src';
createAsset(rootFolder + '/author.js');