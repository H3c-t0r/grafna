import fs from 'fs';
import { glob } from 'glob';
import path from 'path';
import process from 'process';

export function getPackageJson() {
  return require(path.resolve(process.cwd(), 'package.json'));
}

export function getPluginJson() {
  return require(path.resolve(process.cwd(), 'plugin.json'));
}

// Support bundling nested plugins by finding all plugin.json files in src directory
// then checking for a sibling module.[jt]sx? file.
export async function getEntries(): Promise<Record<string, string>> {
  const pluginsJson = await glob(path.resolve('**/plugin.json'), { absolute: true });

  const plugins = await Promise.all(
    pluginsJson.map((pluginJson) => {
      const folder = path.dirname(pluginJson);
      return glob(`${folder}/module.{ts,tsx,js,jsx}`, { absolute: true });
    })
  );

  return plugins.reduce((result, modules) => {
    return modules.reduce((result: { [s: string]: string }, module) => {
      result['module'] = module;
      return result;
    }, result);
  }, {});
}

export function hasLicense() {
  return fs.existsSync(path.resolve(process.cwd(), 'LICENSE'));
}
