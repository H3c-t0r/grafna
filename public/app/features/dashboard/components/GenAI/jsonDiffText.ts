import { get, isEqual } from 'lodash';

import { Dashboard } from '@grafana/schema';
import { Diff, jsonDiff } from 'app/features/dashboard/components/VersionHistory/utils';

import { DashboardModel } from '../../state';

export type JSONValue = null | boolean | number | string | JSONArray | JSONObject;

export type JSONArray = JSONValue[];

export type JSONObject = {
  [key: string]: JSONValue;
};

export function orderProperties(obj1: JSONValue, obj2: JSONValue) {
  // If obj1 and obj2 are the same object, return obj2
  if (obj1 === obj2) {
    return obj2; // No need to order properties, they are already the same
  }

  if (Array.isArray(obj1) && Array.isArray(obj2)) {
    // They are both arrays
    return orderArrayProperties(obj1, obj2);
  }

  // Use a type guard to check if they are both non-array objects
  else if (isObject(obj1) && isObject(obj2)) {
    // Both non-array objects
    return orderObjectProperties(obj1, obj2);
  }

  return obj2;
}

export function isObject(obj: JSONValue): obj is JSONObject {
  return typeof obj === 'object' && !Array.isArray(obj) && obj !== null;
}

export function orderObjectProperties(obj1: JSONObject, obj2: JSONObject) {
  const orderedProperties = Object.keys(obj1);
  const orderedObj2: Record<string, JSONValue> = {};

  for (const prop of orderedProperties) {
    if (obj2.hasOwnProperty(prop)) {
      if (Array.isArray(obj1[prop]) && Array.isArray(obj2[prop])) {
        // Recursive call orderProperties for arrays
        orderedObj2[prop] = orderProperties(obj1[prop], obj2[prop]);
      } else if (typeof obj1[prop] === 'object' && typeof obj2[prop] === 'object') {
        // Recursively call orderProperties for nested objects
        orderedObj2[prop] = orderProperties(obj1[prop], obj2[prop]);
      } else {
        orderedObj2[prop] = obj2[prop];
      }
    }
  }
  return orderedObj2;
}

export function orderArrayProperties(obj1: JSONArray, obj2: JSONArray) {
  const orderedObj2: JSONValue[] = new Array(obj1.length).fill(undefined);

  const unseen1 = new Set<number>([...Array(obj1.length).keys()]);
  const unseen2 = new Set<number>([...Array(obj2.length).keys()]);

  // Loop to match up elements that match exactly
  for (let i = 0; i < obj1.length; i++) {
    if (unseen2.size === 0) {
      break;
    }
    let item1 = obj1[i];
    for (let j = 0; j < obj2.length; j++) {
      if (!unseen2.has(j)) {
        continue;
      }
      let item2 = obj2[j];
      item2 = orderProperties(item1, item2);
      if (isEqual(item1, item2)) {
        unseen1.delete(i);
        unseen2.delete(j);
        orderedObj2[i] = item2;
      }
    }
  }

  fillBySimilarity(obj1, obj2, orderedObj2, unseen1, unseen2);

  return orderedObj2.filter((value) => value !== undefined);
}

// Compare all pairings by similarity and match greedily from highest to lowest
// Similarity is simply measured by number of k:v pairs that are identical
// O(n^2), which is more or less unavoidable
// Can be made a better match by using levenshtein distance and Hungarian matching
export function fillBySimilarity(
  // TODO: Investigate not using any
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  obj1: any[],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  obj2: any[],
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  orderedObj2: any[],
  unseen1: Set<number>,
  unseen2: Set<number>
): void {
  let rankings: Record<number, number[][]> = {}; // Maps scores to arrays of value pairs
  // Unpacking it because I'm not sure removing items while iterating is safe

  unseen2.forEach((j: number) => {
    // Index name matches calling function
    let item2 = obj2[j];

    // If not object, or if array, just push item2 to orderedObj2 and remove j from unseen2
    if (typeof item2 !== 'object' || Array.isArray(item2)) {
      orderedObj2.push(item2);
      unseen2.delete(j);
      return;
    }

    unseen1.forEach((i: number) => {
      let item1 = obj1[i];
      if (typeof item1 !== 'object' || Array.isArray(item1)) {
        unseen1.delete(i);
        return;
      }

      let score = 0;

      for (const key in item1) {
        let val1 = item1[key];
        if (!item2.hasOwnProperty(key)) {
          continue;
        }
        let val2 = item2[key];
        if ((typeof val1 !== 'string' && typeof val1 !== 'number') || typeof val1 !== typeof val2) {
          continue;
        }
        if (val1 === val2) {
          if (key === 'id') {
            score += 1000; // Can probably be caught earlier in the call tree.
          }
          score++;
        }
      }

      if (score !== 0) {
        if (rankings[score] === undefined) {
          rankings[score] = [];
        }
        rankings[score].push([i, j]);
      }
    });
  });

  const keys: number[] = Object.keys(rankings).map(Number); // Get keys as an array of numbers
  keys.sort((a, b) => b - a); // Sort in descending order

  for (const key of keys) {
    let pairs: number[][] = rankings[key];
    for (const pair of pairs) {
      const [i, j] = pair;
      if (unseen1.has(i) && unseen2.has(j)) {
        orderedObj2[i] = obj2[j];
        unseen1.delete(i);
        unseen2.delete(j);
      }
    }
  }

  // Get anything that had no matches whatsoever
  for (const j of unseen2) {
    orderedObj2.push(obj2[j]);
  }
}

// Function for removing empty fields from JSON objects
export function removeEmptyFields(input: JSONValue): JSONValue {
  if (input === null || input === '') {
    return null;
  }

  if (Array.isArray(input)) {
    // Filter out empty values and recursively process the non-empty ones
    const filteredArray = input.map((item) => removeEmptyFields(item)).filter((item) => item !== null);

    return filteredArray.length > 0 ? filteredArray : null;
  }

  if (typeof input !== 'object') {
    // If it's not an object, return as is
    return input;
  }

  // For objects, recursively process each key-value pair
  const result: JSONObject = {};
  for (const key in input) {
    const processedValue = removeEmptyFields(input[key]);

    if (processedValue !== null) {
      if (Array.isArray(processedValue) && processedValue.length === 0) {
        continue;
      }

      if (typeof processedValue === 'object') {
        const keys = Object.keys(processedValue);
        if (keys.length === 0) {
          continue;
        }
      }

      result[key] = processedValue;
    }
  }

  return Object.keys(result).length > 0 ? result : null;
}

// Function for reorganizing diffs by path
export function reorganizeDiffs(diffRecord: Record<string, Diff[]>): Record<string, Diff[]> {
  const reorganized: Record<string, Diff[]> = {};

  for (const key in diffRecord) {
    const changes = diffRecord[key];
    for (const change of changes) {
      const newKey = change.path.length === 1 ? key : change.path.slice(0, -1).join('/');
      if (!Object.hasOwn(reorganized, newKey)) {
        reorganized[newKey] = [];
      }
      reorganized[newKey].push(change);
    }
  }

  return reorganized;
}

// Function for separating root and non-root diffs
export function separateRootAndNonRootDiffs(diffRecord: Record<string, Diff[]>): {
  rootDiffs: Record<string, Diff[]>;
  nonRootDiffs: Record<string, Diff[]>;
} {
  const reorganizedDiffs = reorganizeDiffs(diffRecord);

  const rootDiffs: Record<string, Diff[]> = { dashboard: [] };
  const nonRootDiffs: Record<string, Diff[]> = {};

  for (const key in reorganizedDiffs) {
    if (reorganizedDiffs[key][0].path.length === 1) {
      rootDiffs['dashboard'].push(...reorganizedDiffs[key]);
    } else {
      nonRootDiffs[key] = reorganizedDiffs[key];
    }
  }

  return { rootDiffs, nonRootDiffs };
}

// Function for taking a diff and returning human-readable string
export function getDiffString(diff: Diff): string {
  if (diff.path.length >= 2 && diff.path[0] === 'templating' && diff.path[1] === 'list') {
    return '';
  }
  if (diff.path.length >= 2 && 'thresholds' in diff.path) {
    return '';
  }
  let diffString = '';

  const key: string = diff.path[diff.path.length - 1];
  if (diff.op === 'add') {
    diffString = `+\t"${key}": ${JSON.stringify(diff.value)}`;
  } else if (diff.op === 'remove') {
    diffString = `-\t"${key}": ${JSON.stringify(diff.originalValue)}`;
  } else if (diff.op === 'replace') {
    let minusString = `-\t"${key}": ${JSON.stringify(diff.originalValue)}`;
    let plusString = `+\t"${key}": ${JSON.stringify(diff.value)}`;
    diffString = minusString + '\n' + plusString;
  }
  return diffString;
}

// Function for taking a diff record and returning human-readable string
// Extremely specific to panels, to ensure they have displayed titles
function formatDiffsAsString(lhs: unknown, diffRecord: Record<string, Diff[]>): string[] {
  return Object.entries(diffRecord).map(([key, diffs]) => {
    if (diffs.length === 0) {
      return '';
    }
    const diffStrings = diffs.map(getDiffString).filter((diffString) => diffString !== '');
    let titleString = '';
    if (key.includes('panels')) {
      const path = [...diffs[0].path.slice(0, diffs[0].path.indexOf('panels') + 2), 'title'];
      const title = get(lhs, path);
      if (title !== undefined) {
        titleString = ` with title: ${title}`;
      }
    }
    return `Changes for path ${key}${titleString}:\n {\n${diffStrings.join('\n')}\n }`;
  });
}

function jsonSanitize(obj: Dashboard | DashboardModel | null) {
  return JSON.parse(JSON.stringify(obj, null, 2));
}

export function getDashboardStringDiff(dashboard: DashboardModel): { migrationDiff: string; userDiff: string } {
  let originalDashboard = jsonSanitize(dashboard.getOriginalDashboard());
  let dashboardAfterMigration = jsonSanitize(new DashboardModel(originalDashboard).getSaveModelClone());
  let currentDashboard = jsonSanitize(dashboard.getSaveModelClone());

  dashboardAfterMigration = removeEmptyFields(orderProperties(originalDashboard, dashboardAfterMigration));
  currentDashboard = removeEmptyFields(orderProperties(dashboardAfterMigration, currentDashboard));
  originalDashboard = removeEmptyFields(originalDashboard);

  let jsonMigrationDiff = jsonDiff(originalDashboard, dashboardAfterMigration);
  let jsonUserDiff = jsonDiff(dashboardAfterMigration, currentDashboard);

  const { rootDiffs: rootDiffsMigration, nonRootDiffs: nonRootDiffsMigration } =
    separateRootAndNonRootDiffs(jsonMigrationDiff);
  const { rootDiffs: rootDiffsUser, nonRootDiffs: nonRootDiffsUser } = separateRootAndNonRootDiffs(jsonUserDiff);

  const migrationDiff = `${formatDiffsAsString(originalDashboard, rootDiffsMigration)}\n${formatDiffsAsString(
    originalDashboard,
    nonRootDiffsMigration
  )}`.trim();
  const userDiff = `${formatDiffsAsString(dashboardAfterMigration, rootDiffsUser)}\n${formatDiffsAsString(
    dashboardAfterMigration,
    nonRootDiffsUser
  )}`.trim();

  return { migrationDiff, userDiff };
}
