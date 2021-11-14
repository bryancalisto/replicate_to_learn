import fs from 'fs';
import babylon from 'babylon';
import babelTraverse from 'babel-traverse';
import path from 'path';

let fileId = 0;

function createAsset(filename) {
  // Read the source code
  const src = fs.readFileSync(filename, 'utf-8');
  // Parse the source to a AST (abstract syntax tree) so later we can find 'import' tokens easily
  const parsedAst = babylon.parse(src, { sourceType: 'module' });
  // Extract 'ImportDeclaration' nodes (the child dependencies)
  const dependencies = [];
  babelTraverse.default(parsedAst, {
    ImportDeclaration: ({ node }) => {
      dependencies.push(node.source.value);
    }
  })

  return {
    id: fileId++,
    filename,
    dependencies
  }
}

function createGraph(entryFile) {
  const directory = path.dirname(entryFile);

  // Read the dependencies of the entry file
  const entryAsset = createAsset(entryFile);

  // In this list we will have our dependencies graph
  const assets = [entryAsset];

  // Going to use iterative approach, not recursive
  for (const asset of assets) {
    //The associations between this asset and others will be here
    asset.mapping = {};

    // Get the child assets from each dependency
    asset.dependencies.forEach(relativePath => {
      const childAsset = createAsset(path.join(directory, relativePath));
      asset.mapping[relativePath] = childAsset.id;
      assets.push(childAsset);
    });
  }

  return assets;
}

const graph = createGraph('src\\entry.js');
console.log(graph);