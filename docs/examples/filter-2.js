//Now we start slimming things down
//First we switch to a for...of, which lets us loop over every dragon of 
//x.dragons. We also switch to the modern equivalence test and Array.push().
//When we match a dragon we add it to our results array, to return
//at the end.
{ 
    var r = []; 
    for(dragon of x.dragons) { 
        if(dragon.family === 'Metallic') {
            r.push(dragon); 
        } 
    } 
    return r; 
}

