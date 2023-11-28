import { urlUtil, locationUtil } from '@grafana/data';
export function nameExits(dataSources, name) {
    return (dataSources.filter((dataSource) => {
        return dataSource.name.toLowerCase() === name.toLowerCase();
    }).length > 0);
}
export function findNewName(dataSources, name) {
    // Need to loop through current data sources to make sure
    // the name doesn't exist
    while (nameExits(dataSources, name)) {
        // If there's a duplicate name that doesn't end with '-x'
        // we can add -1 to the name and be done.
        if (!nameHasSuffix(name)) {
            name = `${name}-1`;
        }
        else {
            // if there's a duplicate name that ends with '-x'
            // we can try to increment the last digit until the name is unique
            // remove the 'x' part and replace it with the new number
            name = `${getNewName(name)}${incrementLastDigit(getLastDigit(name))}`;
        }
    }
    return name;
}
function nameHasSuffix(name) {
    return name.endsWith('-', name.length - 1);
}
function getLastDigit(name) {
    return parseInt(name.slice(-1), 10);
}
function incrementLastDigit(digit) {
    return isNaN(digit) ? 1 : digit + 1;
}
function getNewName(name) {
    return name.slice(0, name.length - 1);
}
export const constructDataSourceExploreUrl = (dataSource) => {
    const exploreState = JSON.stringify({ datasource: dataSource.name, context: 'explore' });
    const exploreUrl = urlUtil.renderUrl(locationUtil.assureBaseUrl('/explore'), { left: exploreState });
    return exploreUrl;
};
//# sourceMappingURL=utils.js.map