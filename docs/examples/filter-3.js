//Continuing to trim stuff we switch to the Array filter method
//It takes an arrow function that filters out any elements that don't return
//true for our test, and returns the filtered array, which we put in r and then 
//return
{ 
    var r = x.dragons.filter((dragon)=>{return dragon.family === "Metallic"});
    return r; 
}

