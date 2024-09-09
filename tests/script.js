{
    let a = x.filter(i=>i.filePath.match(/.*dragon.*/)); 
    return {issues:a, numIssues:a.length};
}

