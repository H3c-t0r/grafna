export function generateCopiedName(originalName, exisitingNames) {
    const nonDuplicateName = originalName.replace(/\(copy( [0-9]+)?\)$/, '').trim();
    let newName = `${nonDuplicateName} (copy)`;
    for (let i = 2; exisitingNames.includes(newName); i++) {
        newName = `${nonDuplicateName} (copy ${i})`;
    }
    return newName;
}
//# sourceMappingURL=duplicate.js.map