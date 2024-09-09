{
    // Given an array like:
    // {"z":1, "foo":2, "x51":3, "1":4, "a":32} 
    // add the values

    let n = 0;
    Object.keys(x).forEach(i => {
        n += x[i]
    }); 
    return n;
}

